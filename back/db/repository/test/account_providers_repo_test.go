package test

import (
	"app/commons/constants"
	model "app/db/models"
	"app/db/repository"
	"app/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AccountProvidersRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.AccountProvidersRepository
}

func (suite *AccountProvidersRepoTestSuite) SetupTest() {
	suite.db = testutils.TestDB(suite.T())
	suite.repo = repository.NewAccountProvidersRepository(suite.db)
}

func (suite *AccountProvidersRepoTestSuite) createAccount() model.Account {
	username := "owner"
	account := model.Account{Id: uuid.New(), Username: &username}
	suite.db.Create(&account)
	return account
}

func (suite *AccountProvidersRepoTestSuite) TestCreate_Success() {
	account := suite.createAccount()

	provider := model.AccountProvider{
		AccountId: account.Id,
		Provider:  constants.PROVIDER_GOOGLE,
		Id:        "google-external-id",
	}

	err := suite.repo.Create(provider)
	assert.NoError(suite.T(), err)

	var found model.AccountProvider
	assert.NoError(suite.T(), suite.db.Where("account_id = ? AND provider = ?", account.Id, constants.PROVIDER_GOOGLE).First(&found).Error)
}

func (suite *AccountProvidersRepoTestSuite) TestFindOneById_Found() {
	account := suite.createAccount()
	provider := model.AccountProvider{AccountId: account.Id, Provider: constants.PROVIDER_DISCORD, Id: "Discord-Id-123"}
	suite.Require().NoError(suite.repo.Create(provider))

	var found model.AccountProvider
	err := suite.repo.FindOneById("discord-id-123", "DISCORD", &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.AccountId)
	assert.Equal(suite.T(), account.Id, found.Account.Id, "Account should be preloaded")
}

func (suite *AccountProvidersRepoTestSuite) TestFindOneById_NotFound() {
	var found model.AccountProvider
	err := suite.repo.FindOneById("unknown", "github", &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountProvidersRepoTestSuite) TestDelete_Success() {
	account := suite.createAccount()
	provider := model.AccountProvider{AccountId: account.Id, Provider: constants.PROVIDER_GITHUB, Id: "github-id-456"}
	suite.Require().NoError(suite.repo.Create(provider))

	err := suite.repo.Delete("GITHUB-ID-456")
	assert.NoError(suite.T(), err)

	var count int64
	suite.db.Model(&model.AccountProvider{}).Where("id = ?", provider.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *AccountProvidersRepoTestSuite) TestDelete_NotFound_NoError() {
	// Deleting a non-existent id is a no-op, not an error (GORM Delete with no matches).
	err := suite.repo.Delete("does-not-exist")
	assert.NoError(suite.T(), err)
}

func TestAccountProvidersRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AccountProvidersRepoTestSuite))
}
