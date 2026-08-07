package test

import (
	"app/commons/constants"
	model "app/db/models"
	"app/db/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AccountRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.AccountRepository
}

func (suite *AccountRepoTestSuite) SetupSuite() {
	suite.db = NewTestDB(suite.T())
	suite.repo = repository.NewAccountRepository(suite.db)
}

func (suite *AccountRepoTestSuite) SetupTest() {
	suite.db.Where("1 = 1").Delete(&model.AccountProvider{})
	suite.db.Where("1 = 1").Delete(&model.Account{})
}

func (suite *AccountRepoTestSuite) createDto() repository.AccountCreateDto {
	username := "johndoe"
	email := "john@example.com"
	termsVersion := "v1"
	return repository.AccountCreateDto{
		Id:           uuid.New(),
		UserName:     &username,
		Email:        &email,
		Color:        "#FFFFFF",
		Language:     constants.ACCOUNT_LANGUAGE_EN,
		Password:     "SuperSecret123!",
		AvatarUrl:    "https://example.com/avatar.png",
		TermsVersion: &termsVersion,
		TimeZone:     *time.UTC,
	}
}

func (suite *AccountRepoTestSuite) TestCreate_Success() {
	dto := suite.createDto()

	var account model.Account
	err := suite.repo.Create(dto, &account)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.Id, account.Id)
	assert.NotNil(suite.T(), account.Password)
	assert.NotEqual(suite.T(), dto.Password, *account.Password, "password should be hashed")
	assert.NotNil(suite.T(), account.TermsAcceptedAt, "TermsAcceptedAt should be set when TermsVersion is provided")
}

func (suite *AccountRepoTestSuite) TestCreate_PasswordTooLong() {
	dto := suite.createDto()
	dto.Password = string(make([]byte, 73)) // bcrypt rejects passwords > 72 bytes

	var account model.Account
	err := suite.repo.Create(dto, &account)
	assert.Error(suite.T(), err)
}

func (suite *AccountRepoTestSuite) TestUpdates_PasswordTooLong() {
	dto := suite.createDto()
	dto.Password = ""
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	tooLong := string(make([]byte, 73))
	account.Password = &tooLong

	err := suite.repo.Updates(account)
	assert.Error(suite.T(), err)
}

func (suite *AccountRepoTestSuite) TestCreate_NoPassword() {
	dto := suite.createDto()
	dto.Password = ""
	dto.TermsVersion = nil

	var account model.Account
	err := suite.repo.Create(dto, &account)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), account.Password)
	assert.Nil(suite.T(), account.TermsAcceptedAt)
}

func (suite *AccountRepoTestSuite) TestUpdates_HashesPassword() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	newPassword := "AnotherSecret456!"
	account.Password = &newPassword

	err := suite.repo.Updates(account)
	assert.NoError(suite.T(), err)

	var updated model.Account
	suite.Require().NoError(suite.repo.FindOneById(account.Id, &updated))
	assert.NotNil(suite.T(), updated.Password)
	assert.NotEqual(suite.T(), newPassword, *updated.Password)
}

func (suite *AccountRepoTestSuite) TestFindOneById_Found() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	var found model.Account
	err := suite.repo.FindOneById(account.Id, &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.Id)
}

func (suite *AccountRepoTestSuite) TestFindOneById_NotFound() {
	var found model.Account
	err := suite.repo.FindOneById(uuid.New(), &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestFindOneByUsername_Found() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	var found model.Account
	err := suite.repo.FindOneByUsername("JOHNDOE", &found, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.Id)
}

func (suite *AccountRepoTestSuite) TestFindOneByUsername_ExcludeId() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	var found model.Account
	err := suite.repo.FindOneByUsername("johndoe", &found, &account.Id)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestFindOneByUsername_NotFound() {
	var found model.Account
	err := suite.repo.FindOneByUsername("nobody", &found, nil)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestFindOneByEmail_Found() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	var found model.Account
	err := suite.repo.FindOneByEmail("JOHN@EXAMPLE.COM", &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.Id)
}

func (suite *AccountRepoTestSuite) TestFindOneByEmail_NotFound() {
	var found model.Account
	err := suite.repo.FindOneByEmail("missing@example.com", &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestFindOneByEmailOrUsername() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	var byEmail model.Account
	assert.NoError(suite.T(), suite.repo.FindOneByEmailOrUsername("john@example.com", &byEmail))
	assert.Equal(suite.T(), account.Id, byEmail.Id)

	var byUsername model.Account
	assert.NoError(suite.T(), suite.repo.FindOneByEmailOrUsername("johndoe", &byUsername))
	assert.Equal(suite.T(), account.Id, byUsername.Id)

	var notFound model.Account
	assert.ErrorIs(suite.T(), suite.repo.FindOneByEmailOrUsername("missing", &notFound), gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestFindOneByResetToken() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	resetToken := "reset-token-value"
	expiresAt := time.Now().Add(time.Hour)
	suite.Require().NoError(suite.repo.UpdateResetToken(account.Id, &resetToken, &expiresAt))

	var found model.Account
	err := suite.repo.FindOneByResetToken(resetToken, &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.Id)

	var notFound model.Account
	assert.ErrorIs(suite.T(), suite.repo.FindOneByResetToken("unknown-token", &notFound), gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestUpdateResetToken_Clear() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	resetToken := "reset-token-value"
	expiresAt := time.Now().Add(time.Hour)
	suite.Require().NoError(suite.repo.UpdateResetToken(account.Id, &resetToken, &expiresAt))

	assert.NoError(suite.T(), suite.repo.UpdateResetToken(account.Id, nil, nil))

	var found model.Account
	suite.Require().NoError(suite.repo.FindOneById(account.Id, &found))
	assert.Nil(suite.T(), found.ResetToken)
	assert.Nil(suite.T(), found.PasswordResetTokenAt)
}

func (suite *AccountRepoTestSuite) TestDelete_Success() {
	dto := suite.createDto()
	var account model.Account
	suite.Require().NoError(suite.repo.Create(dto, &account))

	err := suite.repo.Delete(account.Id)
	assert.NoError(suite.T(), err)

	var found model.Account
	assert.ErrorIs(suite.T(), suite.repo.FindOneById(account.Id, &found), gorm.ErrRecordNotFound)
}

func (suite *AccountRepoTestSuite) TestDelete_NotFound() {
	err := suite.repo.Delete(uuid.New())
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func TestAccountRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepoTestSuite))
}
