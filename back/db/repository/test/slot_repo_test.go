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

type SlotRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.SlotRepository
}

func (suite *SlotRepoTestSuite) SetupSuite() {
	suite.db = NewTestDB(suite.T())
	suite.repo = repository.NewSlotRepository(suite.db)
}

func (suite *SlotRepoTestSuite) SetupTest() {
	suite.db.Where("1 = 1").Delete(&model.Slot{})
	suite.db.Where("1 = 1").Delete(&model.Event{})
	suite.db.Where("1 = 1").Delete(&model.Account{})
}

func (suite *SlotRepoTestSuite) createAccount() model.Account {
	username := "owner"
	account := model.Account{Id: uuid.New(), UserName: &username}
	suite.db.Create(&account)
	return account
}

func (suite *SlotRepoTestSuite) createEvent(ownerId uuid.UUID) model.Event {
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

func (suite *SlotRepoTestSuite) TestCreate_Success() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)

	slot := model.Slot{
		Id:       uuid.New(),
		EventId:  event.Id,
		StartsAt: event.StartsAt,
		EndsAt:   event.StartsAt.Add(30 * time.Minute),
	}

	err := suite.repo.Create(&slot)
	assert.NoError(suite.T(), err)

	var found model.Slot
	assert.NoError(suite.T(), suite.db.Where("id = ?", slot.Id).First(&found).Error)
}

func (suite *SlotRepoTestSuite) TestUpdates_Success() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	slot := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	suite.Require().NoError(suite.repo.Create(&slot))

	slot.IsValidated = true
	err := suite.repo.Updates(&slot)
	assert.NoError(suite.T(), err)

	var found model.Slot
	suite.Require().NoError(suite.db.Where("id = ?", slot.Id).First(&found).Error)
	assert.True(suite.T(), found.IsValidated)
}

func (suite *SlotRepoTestSuite) TestFindOneById_Found() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	slot := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	suite.Require().NoError(suite.repo.Create(&slot))

	var found model.Slot
	err := suite.repo.FindOneById(slot.Id, &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), slot.Id, found.Id)
	assert.Equal(suite.T(), event.Id, found.Event.Id)
}

func (suite *SlotRepoTestSuite) TestFindOneById_NotFound() {
	var found model.Slot
	err := suite.repo.FindOneById(uuid.New(), &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *SlotRepoTestSuite) TestFindByEventId() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	slot1 := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	slot2 := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(time.Hour)}
	suite.Require().NoError(suite.repo.Create(&slot1))
	suite.Require().NoError(suite.repo.Create(&slot2))

	var slots []model.Slot
	err := suite.repo.FindByEventId(event.Id, &slots)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), slots, 2)
}

func (suite *SlotRepoTestSuite) TestFindValidatedSlotByEventId() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	validated := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	notValidated := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(time.Hour)}
	suite.Require().NoError(suite.repo.Create(&validated))
	suite.Require().NoError(suite.repo.Create(&notValidated))

	var slot model.Slot
	err := suite.repo.FindValidatedSlotByEventId(event.Id, &slot)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), validated.Id, slot.Id)
}

func (suite *SlotRepoTestSuite) TestDeleteByEventId() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	slot := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	suite.Require().NoError(suite.repo.Create(&slot))

	err := suite.repo.DeleteByEventId(event.Id)
	assert.NoError(suite.T(), err)

	var slots []model.Slot
	suite.db.Where("event_id = ?", event.Id).Find(&slots)
	assert.Len(suite.T(), slots, 0)
}

func (suite *SlotRepoTestSuite) TestDeleteValidatedSlotByEventId() {
	account := suite.createAccount()
	event := suite.createEvent(account.Id)
	validated := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	notValidated := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt.Add(30 * time.Minute), EndsAt: event.StartsAt.Add(time.Hour)}
	suite.Require().NoError(suite.repo.Create(&validated))
	suite.Require().NoError(suite.repo.Create(&notValidated))

	err := suite.repo.DeleteValidatedSlotByEventId(event.Id)
	assert.NoError(suite.T(), err)

	var slots []model.Slot
	suite.db.Where("event_id = ?", event.Id).Find(&slots)
	assert.Len(suite.T(), slots, 1)
	assert.Equal(suite.T(), notValidated.Id, slots[0].Id)
}

func TestSlotRepoTestSuite(t *testing.T) {
	suite.Run(t, new(SlotRepoTestSuite))
}
