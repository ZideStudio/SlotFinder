package provider

import (
	"app/commons/constants"
	"app/commons/encryption"
	"app/commons/guard"
	"app/config"
	model "app/db/models"
	"app/db/repository"
	"app/pkg/account"
	"app/pkg/mail"
	"app/pkg/signin"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
	"net/url"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderService struct {
	accountProvidersRepository *repository.AccountProvidersRepository
	accountRepository          *repository.AccountRepository
	signinService              *signin.SigninService
	accountService             *account.AccountService
	avatarService              *account.AvatarService
	mailService                *mail.MailService
	config                     *config.Config
}

func NewProviderService(service *ProviderService) *ProviderService {
	if service != nil {
		return service
	}

	return &ProviderService{
		accountProvidersRepository: &repository.AccountProvidersRepository{},
		accountRepository:          &repository.AccountRepository{},
		signinService:              signin.NewSigninService(nil),
		accountService:             account.NewAccountService(nil),
		avatarService:              account.NewAvatarService(nil),
		mailService:                mail.NewMailService(nil),
		config:                     config.GetConfig(),
	}
}

var (
	PROVIDER_DISCORD_URL = "https://discord.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=identify+email&state=%s"
	PROVIDER_GOOGLE_URL  = "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email%%20profile&state=%s"
	PROVIDER_GITHUB_URL  = "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s"
)

func (*ProviderService) parseProvider(provider string) (constants.Provider, error) {
	switch provider {
	case string(constants.PROVIDER_GOOGLE):
		return constants.PROVIDER_GOOGLE, nil
	case string(constants.PROVIDER_DISCORD):
		return constants.PROVIDER_DISCORD, nil
	case string(constants.PROVIDER_GITHUB):
		return constants.PROVIDER_GITHUB, nil
	default:
		return "", errors.New("invalid provider")
	}
}

func (s *ProviderService) GetProviderUrl(providerEntry, returnUrl string, user *guard.Claims) (string, error) {
	provider, err := s.parseProvider(providerEntry)
	if err != nil {
		return "", err
	}

	var userId string
	if user != nil {
		userId = user.Id.String()
	}
	jsonState, _ := json.Marshal(map[string]string{
		"userId":    userId,
		"returnUrl": returnUrl,
	})
	jsonStateEncrypted, err := encryption.Encrypt(string(jsonState))
	if err != nil {
		return "", err
	}

	providerConfig := config.GetProviderConfig()
	switch provider {
	case constants.PROVIDER_GOOGLE:
		return fmt.Sprintf(PROVIDER_GOOGLE_URL, providerConfig.GoogleClientId, url.QueryEscape(providerConfig.GoogleRedirectUrl), url.QueryEscape(jsonStateEncrypted)), nil
	case constants.PROVIDER_DISCORD:
		return fmt.Sprintf(PROVIDER_DISCORD_URL, providerConfig.DiscordClientId, url.QueryEscape(providerConfig.DiscordRedirectUrl), url.QueryEscape(jsonStateEncrypted)), nil
	case constants.PROVIDER_GITHUB:
		return fmt.Sprintf(PROVIDER_GITHUB_URL, providerConfig.GithubClientId, url.QueryEscape(providerConfig.GithubRedirectUrl), url.QueryEscape(jsonStateEncrypted)), nil
	default:
		return "", errors.New("unsupported provider")
	}
}

type ProviderAccount struct {
	Id        string
	Username  string
	Email     *string
	AvatarUrl *string
}

type CreateProviderAccountDto struct {
	ProviderAccount ProviderAccount
	Provider        constants.Provider
}

type ProviderAccountReponse struct {
	Account         *repository.AccountCreateDto
	AccountProvider *model.AccountProvider
	Jwt             *signin.TokenResponseDto
}

func (s *ProviderService) createProviderAccount(providerUser CreateProviderAccountDto, authUserId string) (providerAccountResponse ProviderAccountReponse, err error) {
	var existingAccountProvider model.AccountProvider
	if err := s.accountProvidersRepository.FindOneById(providerUser.ProviderAccount.Id, string(providerUser.Provider), &existingAccountProvider); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return providerAccountResponse, fmt.Errorf("error finding provider: %w", err)
	}

	if authUserId != "" && authUserId != existingAccountProvider.AccountId.String() {
		if existingAccountProvider.AccountId != uuid.Nil {
			if existingAccountProvider.AccountId.String() != authUserId {
				return providerAccountResponse, fmt.Errorf("ALREADY_EXISTS: provider already exists connected to another account")
			}

			if err := s.accountProvidersRepository.Delete(providerUser.ProviderAccount.Id); err != nil {
				return providerAccountResponse, fmt.Errorf("error finding provider: %w", err)
			}
		}

		providerAccountResponse.AccountProvider = &model.AccountProvider{
			Provider: providerUser.Provider,
			Id:       providerUser.ProviderAccount.Id,
		}

		return providerAccountResponse, nil
	} else if existingAccountProvider.AccountId != uuid.Nil {
		jwt, err := s.signinService.GenerateToken(&guard.Claims{
			Id:       existingAccountProvider.AccountId,
			Username: existingAccountProvider.Account.UserName,
			Email:    existingAccountProvider.Account.Email,
		})
		if err != nil {
			return providerAccountResponse, err
		}

		providerAccountResponse.Jwt = &jwt

		return providerAccountResponse, nil
	}

	providerAccountResponse.Account = &repository.AccountCreateDto{
		UserName: &providerUser.ProviderAccount.Username,
		Email:    providerUser.ProviderAccount.Email,
		Providers: []model.AccountProvider{
			{
				Provider: providerUser.Provider,
				Id:       providerUser.ProviderAccount.Id,
			},
		},
	}

	return providerAccountResponse, nil
}

func (s *ProviderService) ProviderCallback(providerEntry string, code string, userId string) (signin.TokenResponseDto, error) {
	var tokenResponse signin.TokenResponseDto
	provider, err := s.parseProvider(providerEntry)
	if err != nil {
		return tokenResponse, err
	}

	var providerAccount ProviderAccount
	switch provider {
	case constants.PROVIDER_GOOGLE:
		providerAccount, err = s.getGoogleUserInfo(code)
		if err != nil {
			return tokenResponse, fmt.Errorf("OAUTH: failed to get Google user info: %w", err)
		}
	case constants.PROVIDER_DISCORD:
		providerAccount, err = s.getDiscordUserInfo(code)
		if err != nil {
			return tokenResponse, fmt.Errorf("OAUTH: failed to get Discord user info: %w", err)
		}
	case constants.PROVIDER_GITHUB:
		providerAccount, err = s.getGithubUserInfo(code)
		if err != nil {
			return tokenResponse, fmt.Errorf("OAUTH: failed to get Github user info: %w", err)
		}
	default:
		return tokenResponse, errors.New("unsupported provider")
	}

	providerAccountResponse, err := s.createProviderAccount(CreateProviderAccountDto{
		ProviderAccount: providerAccount,
		Provider:        provider,
	}, userId)
	if err != nil {
		return tokenResponse, fmt.Errorf("failed to create provider account: %w", err)
	}
	if providerAccountResponse.Jwt != nil { // log user
		return *providerAccountResponse.Jwt, nil
	}

	if providerAccountResponse.Account != nil { // New account
		// Ensure username and email is provided
		if providerAccountResponse.Account.UserName == nil || *providerAccountResponse.Account.UserName == "" {
			return tokenResponse, errors.New("username should be provided by provider")
		}
		if providerAccountResponse.Account.Email == nil || *providerAccountResponse.Account.Email == "" {
			return tokenResponse, errors.New("email should be provided by provider")
		}

		// Check if email already exists
		var existingAccount model.Account
		if err := s.accountRepository.FindOneByEmail(*providerAccountResponse.Account.Email, &existingAccount); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return tokenResponse, err
		}
		if existingAccount.Id != uuid.Nil {
			return tokenResponse, constants.ERR_EMAIL_ALREADY_EXISTS.Err
		}

		// Check if username is available
		isUsernameAvailable, err := s.accountService.CheckUserNameAvailability(*providerAccountResponse.Account.UserName, nil)
		if err != nil {
			return tokenResponse, err
		}
		if !isUsernameAvailable {
			providerAccountResponse.Account.UserName = nil
		}

		// Setup avatar
		accountId := uuid.New()
		avatarUrl := ""
		if providerAccount.AvatarUrl != nil {
			uploadedAvatarUrl, err := s.avatarService.UploadAvatar(providerAccount.AvatarUrl, nil, accountId.String())
			if err == nil {
				avatarUrl = uploadedAvatarUrl
			}
		}
		if avatarUrl == "" {
			avatarUrl = s.avatarService.GetGravatarURL(accountId.String())
		}

		// Choose a random color
		colors := constants.COLORS
		color := colors[mathrand.Intn(len(colors))]

		// Create account
		var account model.Account
		if err := s.accountRepository.Create(repository.AccountCreateDto{
			Id:        accountId,
			UserName:  providerAccountResponse.Account.UserName,
			Email:     providerAccountResponse.Account.Email,
			Color:     string(color),
			AvatarUrl: avatarUrl,
			Providers: providerAccountResponse.Account.Providers,
		}, &account); err != nil {
			return tokenResponse, fmt.Errorf("error creating account: %w", err)
		}

		// Generate token
		token, err := s.signinService.GenerateToken(&guard.Claims{
			Id:       account.Id,
			Username: account.UserName,
			Email:    account.Email,
		})
		if err != nil {
			_ = s.accountRepository.Delete(account.Id)
			return tokenResponse, err
		}

		// Send welcome email
		if account.Email != nil {
			go s.mailService.SendMail(mail.EmailParams{
				Template: constants.MAIL_TEMPLATE_WELCOME,
				To:       *account.Email,
				Subject:  "Welcome to SlotFinder!",
				Params: map[string]string{
					"LoginUrl": fmt.Sprintf("%s/login", s.config.Origin),
				},
			})
		}

		return token, nil
	}

	// Create provider on existing account
	if providerAccountResponse.AccountProvider != nil {
		if err := s.accountProvidersRepository.Create(*providerAccountResponse.AccountProvider); err != nil {
			return tokenResponse, fmt.Errorf("error updating account: %w", err)
		}

		userUuid, err := uuid.Parse(userId)
		if err != nil {
			return tokenResponse, fmt.Errorf("invalid userId: %w", err)
		}
		var account model.Account
		if err := s.accountRepository.FindOneById(userUuid, &account); err != nil {
			return tokenResponse, err
		}

		token, err := s.signinService.GenerateToken(&guard.Claims{
			Id:       account.Id,
			Username: account.UserName,
			Email:    account.Email,
		})
		if err != nil {
			return tokenResponse, err
		}

		return token, nil
	}

	return tokenResponse, fmt.Errorf("unhandled provider service error")
}
