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

type AccountEventRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.AccountEventRepository
}

func (suite *AccountEventRepoTestSuite) SetupSuite() {
	suite.db = NewTestDB(suite.T())
	suite.repo = repository.NewAccountEventRepository(suite.db)
}

func (suite *AccountEventRepoTestSuite) SetupTest() {
	suite.db.Where("1 = 1").Delete(&model.AccountEvent{})
	suite.db.Where("1 = 1").Delete(&model.Event{})
	suite.db.Where("1 = 1").Delete(&model.Account{})
}

func (suite *AccountEventRepoTestSuite) createAccount(username string) model.Account {
	account := model.Account{Id: uuid.New(), UserName: &username}
	suite.db.Create(&account)
	return account
}

func (suite *AccountEventRepoTestSuite) createEvent(ownerId uuid.UUID) model.Event {
	event := model.Event{
		Id:       uuid.New(),
		Name:     "Test Event",
		Duration: 60,
		StartsAt: time.Now().Add(time.Hour),
		EndsAt:   time.Now().Add(2 * time.Hour),
		OwnerId:  ownerId,
		Status:   "IN_DECISION",
	}
	suite.db.Create(&event)
	return event
}

func (suite *AccountEventRepoTestSuite) TestCreate_Success() {
	owner := suite.createAccount("owner")
	event := suite.createEvent(owner.Id)

	accountEvent := model.AccountEvent{AccountId: owner.Id, EventId: event.Id}
	err := suite.repo.Create(&accountEvent)
	assert.NoError(suite.T(), err)

	var found model.AccountEvent
	assert.NoError(suite.T(), suite.db.Where("account_id = ? AND event_id = ?", owner.Id, event.Id).First(&found).Error)
}

func (suite *AccountEventRepoTestSuite) TestUpdates_Success() {
	owner := suite.createAccount("owner")
	event := suite.createEvent(owner.Id)
	accountEvent := model.AccountEvent{AccountId: owner.Id, EventId: event.Id}
	suite.Require().NoError(suite.repo.Create(&accountEvent))

	color := "#123456"
	accountEvent.Color = &color
	err := suite.repo.Updates(&accountEvent)
	assert.NoError(suite.T(), err)

	var found model.AccountEvent
	suite.Require().NoError(suite.db.Where("account_id = ? AND event_id = ?", owner.Id, event.Id).First(&found).Error)
	assert.NotNil(suite.T(), found.Color)
	assert.Equal(suite.T(), color, *found.Color)
}

func (suite *AccountEventRepoTestSuite) TestFindByAccountAndEventId_Found() {
	owner := suite.createAccount("owner")
	event := suite.createEvent(owner.Id)
	accountEvent := model.AccountEvent{AccountId: owner.Id, EventId: event.Id}
	suite.Require().NoError(suite.repo.Create(&accountEvent))

	var found model.AccountEvent
	err := suite.repo.FindByAccountAndEventId(owner.Id, event.Id, &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), owner.Id, found.AccountId)
	assert.Equal(suite.T(), event.Id, found.EventId)
}

func (suite *AccountEventRepoTestSuite) TestFindByAccountAndEventId_NotFound() {
	var found model.AccountEvent
	err := suite.repo.FindByAccountAndEventId(uuid.New(), uuid.New(), &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *AccountEventRepoTestSuite) TestFindAccountsByEventId() {
	owner := suite.createAccount("owner")
	participant := suite.createAccount("participant")
	event := suite.createEvent(owner.Id)

	suite.Require().NoError(suite.repo.Create(&model.AccountEvent{AccountId: owner.Id, EventId: event.Id}))
	suite.Require().NoError(suite.repo.Create(&model.AccountEvent{AccountId: participant.Id, EventId: event.Id}))

	var accounts []model.Account
	err := suite.repo.FindAccountsByEventId(event.Id, &accounts)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), accounts, 2)

	ids := map[uuid.UUID]bool{}
	for _, a := range accounts {
		ids[a.Id] = true
	}
	assert.True(suite.T(), ids[owner.Id])
	assert.True(suite.T(), ids[participant.Id])
}

func (suite *AccountEventRepoTestSuite) TestFindAccountsByEventId_Empty() {
	var accounts []model.Account
	err := suite.repo.FindAccountsByEventId(uuid.New(), &accounts)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), accounts, 0)
}

func TestAccountEventRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AccountEventRepoTestSuite))
}
