package test

import (
	model "app/db/models"
	"app/db/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RefreshTokenRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.RefreshTokenRepository
}

func (suite *RefreshTokenRepoTestSuite) SetupSuite() {
	suite.db = NewTestDB(suite.T())
	suite.repo = repository.NewRefreshTokenRepository(suite.db)
}

func (suite *RefreshTokenRepoTestSuite) SetupTest() {
	suite.db.Where("1 = 1").Delete(&model.RefreshToken{})
	suite.db.Where("1 = 1").Delete(&model.Account{})
}

func (suite *RefreshTokenRepoTestSuite) createAccount() model.Account {
	username := "owner"
	account := model.Account{Id: uuid.New(), UserName: &username}
	suite.db.Create(&account)
	return account
}

func (suite *RefreshTokenRepoTestSuite) TestGenerateRefreshToken_Unique() {
	token1, err1 := suite.repo.GenerateRefreshToken()
	token2, err2 := suite.repo.GenerateRefreshToken()

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEmpty(suite.T(), token1)
	assert.NotEqual(suite.T(), token1, token2)
}

func (suite *RefreshTokenRepoTestSuite) TestHashToken_Deterministic() {
	hash1 := suite.repo.HashToken("some-token")
	hash2 := suite.repo.HashToken("some-token")
	hash3 := suite.repo.HashToken("other-token")

	assert.Equal(suite.T(), hash1, hash2)
	assert.NotEqual(suite.T(), hash1, hash3)
}

func (suite *RefreshTokenRepoTestSuite) TestCreate_Success() {
	account := suite.createAccount()

	token, err := suite.repo.Create(account.Id, time.Now().Add(time.Hour))
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	var stored model.RefreshToken
	tokenHash := suite.repo.HashToken(token)
	assert.NoError(suite.T(), suite.db.Where("token_hash = ?", tokenHash).First(&stored).Error)
	assert.Equal(suite.T(), account.Id, stored.AccountId)
	assert.False(suite.T(), stored.IsRevoked)
}

func (suite *RefreshTokenRepoTestSuite) TestFindByTokenHash_Found() {
	account := suite.createAccount()
	token, err := suite.repo.Create(account.Id, time.Now().Add(time.Hour))
	suite.Require().NoError(err)

	var found model.RefreshToken
	err = suite.repo.FindByTokenHash(suite.repo.HashToken(token), &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), account.Id, found.AccountId)
}

func (suite *RefreshTokenRepoTestSuite) TestFindByTokenHash_NotFound() {
	var found model.RefreshToken
	err := suite.repo.FindByTokenHash("unknown-hash", &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RefreshTokenRepoTestSuite) TestFindByTokenHash_Revoked() {
	account := suite.createAccount()
	token, err := suite.repo.Create(account.Id, time.Now().Add(time.Hour))
	suite.Require().NoError(err)

	var stored model.RefreshToken
	suite.Require().NoError(suite.db.Where("token_hash = ?", suite.repo.HashToken(token)).First(&stored).Error)
	suite.Require().NoError(suite.repo.Revoke(stored.Id))

	var found model.RefreshToken
	err = suite.repo.FindByTokenHash(suite.repo.HashToken(token), &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound, "revoked tokens should not be findable")
}

func (suite *RefreshTokenRepoTestSuite) TestRevoke() {
	account := suite.createAccount()
	token, err := suite.repo.Create(account.Id, time.Now().Add(time.Hour))
	suite.Require().NoError(err)

	var stored model.RefreshToken
	suite.Require().NoError(suite.db.Where("token_hash = ?", suite.repo.HashToken(token)).First(&stored).Error)

	err = suite.repo.Revoke(stored.Id)
	assert.NoError(suite.T(), err)

	var updated model.RefreshToken
	suite.Require().NoError(suite.db.Where("id = ?", stored.Id).First(&updated).Error)
	assert.True(suite.T(), updated.IsRevoked)
	assert.NotNil(suite.T(), updated.RevokedAt)
}

func (suite *RefreshTokenRepoTestSuite) TestRevokeAllForAccount() {
	account := suite.createAccount()
	_, err := suite.repo.Create(account.Id, time.Now().Add(time.Hour))
	suite.Require().NoError(err)
	_, err = suite.repo.Create(account.Id, time.Now().Add(2*time.Hour))
	suite.Require().NoError(err)

	err = suite.repo.RevokeAllForAccount(account.Id)
	assert.NoError(suite.T(), err)

	var tokens []model.RefreshToken
	suite.db.Where("account_id = ?", account.Id).Find(&tokens)
	assert.Len(suite.T(), tokens, 2)
	for _, t := range tokens {
		assert.True(suite.T(), t.IsRevoked)
	}
}

func (suite *RefreshTokenRepoTestSuite) TestDeleteExpired() {
	account := suite.createAccount()

	// Expired more than a week ago and revoked: should be deleted.
	oldRevoked := model.RefreshToken{
		Id:        uuid.New(),
		AccountId: account.Id,
		TokenHash: "old-revoked",
		ExpiresAt: time.Now().AddDate(0, 0, -10),
		IsRevoked: true,
	}
	// Expired more than a week ago but NOT revoked: should be kept.
	oldNotRevoked := model.RefreshToken{
		Id:        uuid.New(),
		AccountId: account.Id,
		TokenHash: "old-not-revoked",
		ExpiresAt: time.Now().AddDate(0, 0, -10),
		IsRevoked: false,
	}
	// Recently expired and revoked: should be kept (not old enough).
	recentRevoked := model.RefreshToken{
		Id:        uuid.New(),
		AccountId: account.Id,
		TokenHash: "recent-revoked",
		ExpiresAt: time.Now().Add(-time.Hour),
		IsRevoked: true,
	}
	suite.Require().NoError(suite.db.Create(&oldRevoked).Error)
	suite.Require().NoError(suite.db.Create(&oldNotRevoked).Error)
	suite.Require().NoError(suite.db.Create(&recentRevoked).Error)

	err := suite.repo.DeleteExpired()
	assert.NoError(suite.T(), err)

	var remaining []model.RefreshToken
	suite.db.Where("account_id = ?", account.Id).Find(&remaining)
	assert.Len(suite.T(), remaining, 2)

	remainingIds := map[uuid.UUID]bool{}
	for _, t := range remaining {
		remainingIds[t.Id] = true
	}
	assert.False(suite.T(), remainingIds[oldRevoked.Id])
	assert.True(suite.T(), remainingIds[oldNotRevoked.Id])
	assert.True(suite.T(), remainingIds[recentRevoked.Id])
}

func TestRefreshTokenRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenRepoTestSuite))
}
