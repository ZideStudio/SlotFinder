package test

import (
	"app/commons/constants"
	model "app/db/models"
	"app/db/repository"
	"app/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type EventRepoTestSuite struct {
	suite.Suite
	db             *gorm.DB
	repo           *repository.EventRepository
	accountEventDb *repository.AccountEventRepository
}

func (suite *EventRepoTestSuite) SetupTest() {
	suite.db = testutils.TestDB(suite.T())
	suite.repo = repository.NewEventRepository(suite.db)
	suite.accountEventDb = repository.NewAccountEventRepository(suite.db)
}

func (suite *EventRepoTestSuite) createAccount(username string) model.Account {
	account := model.Account{Id: uuid.New(), UserName: &username}
	suite.db.Create(&account)
	return account
}

func (suite *EventRepoTestSuite) TestCreate_Success() {
	owner := suite.createAccount("owner")
	event := model.Event{
		Id:       uuid.New(),
		Name:     "New Event",
		Duration: 45,
		StartsAt: time.Now().Add(time.Hour),
		EndsAt:   time.Now().Add(2 * time.Hour),
		OwnerId:  owner.Id,
		Status:   "IN_DECISION",
	}

	err := suite.repo.Create(&event)
	assert.NoError(suite.T(), err)

	var found model.Event
	assert.NoError(suite.T(), suite.db.Where("id = ?", event.Id).First(&found).Error)
}

// Covers every constants.EventStatuses value, so one added without a
// matching migration fails here instead of only in production. SQLite never
// enforced this enum at all — only meaningful now that tests hit Postgres.
func (suite *EventRepoTestSuite) TestCreate_EachValidStatus_Accepted() {
	owner := suite.createAccount("owner")

	for _, status := range constants.EventStatuses {
		event := model.Event{
			Id: uuid.New(), Name: "Event " + string(status), Duration: 30,
			StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour),
			OwnerId: owner.Id, Status: status,
		}
		suite.Require().NoError(suite.repo.Create(&event), "status %q should be accepted by the event_status enum", status)

		var found model.Event
		suite.Require().NoError(suite.db.Where("id = ?", event.Id).First(&found).Error)
		assert.Equal(suite.T(), status, found.Status)
	}
}

// Asserts the event_status enum actually constrains the column at the
// database level, not just via application-side validation.
func (suite *EventRepoTestSuite) TestCreate_InvalidStatus_RejectedByEnum() {
	owner := suite.createAccount("owner")
	event := model.Event{
		Id: uuid.New(), Name: "Event", Duration: 30,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour),
		OwnerId: owner.Id, Status: constants.EventStatus("NOT_A_REAL_STATUS"),
	}

	err := suite.repo.Create(&event)
	assert.Error(suite.T(), err)
}

func (suite *EventRepoTestSuite) TestUpdates_Success() {
	owner := suite.createAccount("owner")
	event := model.Event{
		Id: uuid.New(), Name: "Original", Duration: 30,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour),
		OwnerId: owner.Id, Status: "IN_DECISION",
	}
	suite.Require().NoError(suite.repo.Create(&event))

	event.Name = "Updated"
	err := suite.repo.Updates(&event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated", event.Name)

	var found model.Event
	suite.Require().NoError(suite.db.Where("id = ?", event.Id).First(&found).Error)
	assert.Equal(suite.T(), "Updated", found.Name)
}

func (suite *EventRepoTestSuite) TestFindOneById_Found() {
	owner := suite.createAccount("owner")
	participant := suite.createAccount("participant")
	event := model.Event{
		Id: uuid.New(), Name: "Event", Duration: 30,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour),
		OwnerId: owner.Id, Status: "IN_DECISION",
	}
	suite.Require().NoError(suite.repo.Create(&event))
	suite.Require().NoError(suite.accountEventDb.Create(&model.AccountEvent{AccountId: owner.Id, EventId: event.Id}))
	suite.Require().NoError(suite.accountEventDb.Create(&model.AccountEvent{AccountId: participant.Id, EventId: event.Id}))

	var found model.Event
	err := suite.repo.FindOneById(event.Id, &found)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.Id, found.Id)
	assert.Len(suite.T(), found.Participants, 2)
	// Owner is replaced by its sanitized copy (Id intentionally stripped), so we
	// can only assert on the fields Sanitized() preserves.
	assert.NotNil(suite.T(), found.Owner.UserName)
	assert.Equal(suite.T(), *owner.UserName, *found.Owner.UserName)
}

func (suite *EventRepoTestSuite) TestFindOneById_NotFound() {
	var found model.Event
	err := suite.repo.FindOneById(uuid.New(), &found)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *EventRepoTestSuite) TestFindEventsByAccountId_SortOrderByStatus() {
	owner := suite.createAccount("owner")

	finished := model.Event{Id: uuid.New(), Name: "Finished Event", Duration: 30, StartsAt: time.Now(), EndsAt: time.Now(), OwnerId: owner.Id, Status: "FINISHED"}
	inDecision := model.Event{Id: uuid.New(), Name: "InDecision Event", Duration: 30, StartsAt: time.Now(), EndsAt: time.Now(), OwnerId: owner.Id, Status: "IN_DECISION"}
	upcoming := model.Event{Id: uuid.New(), Name: "Upcoming Event", Duration: 30, StartsAt: time.Now(), EndsAt: time.Now(), OwnerId: owner.Id, Status: "UPCOMING"}

	suite.Require().NoError(suite.repo.Create(&finished))
	suite.Require().NoError(suite.repo.Create(&inDecision))
	suite.Require().NoError(suite.repo.Create(&upcoming))

	for _, e := range []model.Event{finished, inDecision, upcoming} {
		suite.Require().NoError(suite.accountEventDb.Create(&model.AccountEvent{AccountId: owner.Id, EventId: e.Id}))
	}

	events, total, err := suite.repo.FindEventsByAccountId(owner.Id, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), events, 3)
	assert.Equal(suite.T(), inDecision.Id, events[0].Id)
	assert.Equal(suite.T(), upcoming.Id, events[1].Id)
	assert.Equal(suite.T(), finished.Id, events[2].Id)
}

func (suite *EventRepoTestSuite) TestFindEventsByAccountId_Pagination() {
	owner := suite.createAccount("owner")

	names := []string{"A Event", "B Event", "C Event"}
	for _, name := range names {
		event := model.Event{Id: uuid.New(), Name: name, Duration: 30, StartsAt: time.Now(), EndsAt: time.Now(), OwnerId: owner.Id, Status: "IN_DECISION"}
		suite.Require().NoError(suite.repo.Create(&event))
		suite.Require().NoError(suite.accountEventDb.Create(&model.AccountEvent{AccountId: owner.Id, EventId: event.Id}))
	}

	page1, total, err := suite.repo.FindEventsByAccountId(owner.Id, 2, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), page1, 2)
	assert.Equal(suite.T(), "A Event", page1[0].Name)
	assert.Equal(suite.T(), "B Event", page1[1].Name)

	page2, total, err := suite.repo.FindEventsByAccountId(owner.Id, 2, 2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), page2, 1)
	assert.Equal(suite.T(), "C Event", page2[0].Name)
}

func (suite *EventRepoTestSuite) TestFindEventsByAccountId_Empty() {
	owner := suite.createAccount("owner")

	events, total, err := suite.repo.FindEventsByAccountId(owner.Id, 10, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), total)
	assert.Len(suite.T(), events, 0)
}

func (suite *EventRepoTestSuite) TestDelete_Success() {
	owner := suite.createAccount("owner")
	event := model.Event{
		Id: uuid.New(), Name: "Event", Duration: 30,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour),
		OwnerId: owner.Id, Status: "IN_DECISION",
	}
	suite.Require().NoError(suite.repo.Create(&event))

	err := suite.repo.Delete(event.Id)
	assert.NoError(suite.T(), err)

	var count int64
	suite.db.Model(&model.Event{}).Where("id = ?", event.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func TestEventRepoTestSuite(t *testing.T) {
	suite.Run(t, new(EventRepoTestSuite))
}

// Uses its own dedicated DB (not the shared suite one) since it drops the
// event table to force the paginated id lookup (Pluck) to fail.
func TestFindEventsByAccountId_PluckQueryFails(t *testing.T) {
	database := testutils.TestDB(t)
	repo := repository.NewEventRepository(database)

	require.NoError(t, database.Migrator().DropTable(&model.Event{}))

	_, _, err := repo.FindEventsByAccountId(uuid.New(), 10, 0)
	assert.Error(t, err)
}
