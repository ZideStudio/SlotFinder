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

// ErrorPathsRepoTestSuite exercises the "unexpected DB error" branches of every
// repository (as opposed to gorm.ErrRecordNotFound, which is expected/handled
// separately). Each test gets a fresh in-memory DB whose underlying connection
// is closed immediately, so every query fails with a generic driver error and
// hits the repositories' log.Error()+return err branches.
type ErrorPathsRepoTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *ErrorPathsRepoTestSuite) SetupTest() {
	suite.db = NewTestDB(suite.T())
	sqlDB, err := suite.db.DB()
	suite.Require().NoError(err)
	suite.Require().NoError(sqlDB.Close())
}

func (suite *ErrorPathsRepoTestSuite) TestAccountRepository_Errors() {
	repo := repository.NewAccountRepository(suite.db)

	var account model.Account
	assert.Error(suite.T(), repo.Create(repository.AccountCreateDto{Id: uuid.New(), Password: "pw"}, &account))
	assert.Error(suite.T(), repo.Updates(model.Account{Id: uuid.New(), Password: strPtr("pw")}))
	assert.Error(suite.T(), repo.FindOneById(uuid.New(), &account))
	assert.Error(suite.T(), repo.FindOneByUsername("user", &account, nil))
	assert.Error(suite.T(), repo.FindOneByEmail("mail@example.com", &account))
	assert.Error(suite.T(), repo.FindOneByEmailOrUsername("mail", &account))
	assert.Error(suite.T(), repo.FindOneByResetToken("token", &account))
	assert.Error(suite.T(), repo.UpdateResetToken(uuid.New(), strPtr("token"), timePtr(time.Now())))
	assert.Error(suite.T(), repo.Delete(uuid.New()))
}

func (suite *ErrorPathsRepoTestSuite) TestSlotRepository_Errors() {
	repo := repository.NewSlotRepository(suite.db)

	slot := &model.Slot{Id: uuid.New()}
	var slots []model.Slot
	assert.Error(suite.T(), repo.Create(slot))
	assert.Error(suite.T(), repo.Updates(slot))
	assert.Error(suite.T(), repo.FindOneById(uuid.New(), slot))
	assert.Error(suite.T(), repo.FindByEventId(uuid.New(), &slots))
	assert.Error(suite.T(), repo.FindValidatedSlotByEventId(uuid.New(), slot))
	assert.Error(suite.T(), repo.DeleteByEventId(uuid.New()))
	assert.Error(suite.T(), repo.DeleteValidatedSlotByEventId(uuid.New()))
}

func (suite *ErrorPathsRepoTestSuite) TestRefreshTokenRepository_Errors() {
	repo := repository.NewRefreshTokenRepository(suite.db)

	var token model.RefreshToken
	_, err := repo.Create(uuid.New(), time.Now().Add(time.Hour))
	assert.Error(suite.T(), err)
	assert.Error(suite.T(), repo.FindByTokenHash("hash", &token))
	assert.Error(suite.T(), repo.Revoke(uuid.New()))
	assert.Error(suite.T(), repo.RevokeAllForAccount(uuid.New()))
	assert.Error(suite.T(), repo.DeleteExpired())
}

func (suite *ErrorPathsRepoTestSuite) TestAccountEventRepository_Errors() {
	repo := repository.NewAccountEventRepository(suite.db)

	accountEvent := &model.AccountEvent{AccountId: uuid.New(), EventId: uuid.New()}
	var accounts []model.Account
	assert.Error(suite.T(), repo.Create(accountEvent))
	assert.Error(suite.T(), repo.Updates(accountEvent))
	assert.Error(suite.T(), repo.FindByAccountAndEventId(uuid.New(), uuid.New(), accountEvent))
	assert.Error(suite.T(), repo.FindAccountsByEventId(uuid.New(), &accounts))
}

func (suite *ErrorPathsRepoTestSuite) TestEventRepository_Errors() {
	repo := repository.NewEventRepository(suite.db)

	event := &model.Event{Id: uuid.New()}
	assert.Error(suite.T(), repo.Create(event))
	assert.Error(suite.T(), repo.Updates(event))
	assert.Error(suite.T(), repo.FindOneById(uuid.New(), event))
	_, _, err := repo.FindEventsByAccountId(uuid.New(), 10, 0)
	assert.Error(suite.T(), err)
	assert.Error(suite.T(), repo.Delete(uuid.New()))
}

func (suite *ErrorPathsRepoTestSuite) TestAccountProvidersRepository_Errors() {
	repo := repository.NewAccountProvidersRepository(suite.db)

	var provider model.AccountProvider
	assert.Error(suite.T(), repo.Create(model.AccountProvider{AccountId: uuid.New(), Provider: constants.PROVIDER_GOOGLE, Id: "id"}))
	assert.Error(suite.T(), repo.FindOneById("id", "google", &provider))
	assert.Error(suite.T(), repo.Delete("id"))
}

func (suite *ErrorPathsRepoTestSuite) TestAvailabilityRepository_Errors() {
	repo := repository.NewAvailabilityRepository(suite.db)

	availability := &model.Availability{Id: uuid.New(), AccountId: uuid.New(), EventId: uuid.New()}
	var availabilities []model.Availability
	ids := []uuid.UUID{uuid.New()}

	assert.Error(suite.T(), repo.FindOverlappingAvailabilities(availability, &availabilities))
	assert.Error(suite.T(), repo.DeleteByIds(&ids))
	assert.Error(suite.T(), repo.Create(availability))
	assert.Error(suite.T(), repo.FindOneById(availability.Id, availability))
	assert.Error(suite.T(), repo.DeleteById(&availability.Id))
	assert.Error(suite.T(), repo.FindByEventId(uuid.New(), &availabilities))
	assert.Error(suite.T(), repo.Update(availability))
	assert.Error(suite.T(), repo.DeleteOutOfEventRangeAndAdjustOverlaps(uuid.New(), time.Now(), time.Now().Add(time.Hour)))
}

func TestErrorPathsRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorPathsRepoTestSuite))
}

// TestConstructors_NilDatabaseFallback covers the `if database == nil` branch
// shared by every repository constructor. It only exercises construction (no
// query is issued), since the process-global DB singleton isn't initialized
// in unit tests.
func TestConstructors_NilDatabaseFallback(t *testing.T) {
	assert.NotNil(t, repository.NewAccountRepository(nil))
	assert.NotNil(t, repository.NewSlotRepository(nil))
	assert.NotNil(t, repository.NewRefreshTokenRepository(nil))
	assert.NotNil(t, repository.NewAccountEventRepository(nil))
	assert.NotNil(t, repository.NewEventRepository(nil))
	assert.NotNil(t, repository.NewAccountProvidersRepository(nil))
	assert.NotNil(t, repository.NewAvailabilityRepository(nil))
}

func strPtr(s string) *string        { return &s }
func timePtr(t time.Time) *time.Time { return &t }
