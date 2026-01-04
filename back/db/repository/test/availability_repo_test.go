package test

import (
	model "app/db/models"
	"app/db/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AvailabilityRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.AvailabilityRepository
}

func (suite *AvailabilityRepoTestSuite) SetupSuite() {
	// Create in-memory SQLite database for testing
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = database

	// Auto-migrate the schema
	err = database.AutoMigrate(&model.Event{}, &model.Account{}, &model.Availability{})
	suite.Require().NoError(err)

	// Create repository with test DB
	suite.repo = repository.NewAvailabilityRepository(database)
}

func (suite *AvailabilityRepoTestSuite) SetupTest() {
	// Clean up tables before each test
	suite.db.Where("1 = 1").Delete(&model.Availability{})
	suite.db.Where("1 = 1").Delete(&model.Event{})
	suite.db.Where("1 = 1").Delete(&model.Account{})
}

func (suite *AvailabilityRepoTestSuite) TearDownSuite() {
	// Close database connection
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Helper function to create test account
func (suite *AvailabilityRepoTestSuite) createTestAccount() model.Account {
	username := "testuser"
	email := "test@example.com"
	account := model.Account{
		Id:       uuid.New(),
		UserName: &username,
		Email:    &email,
	}
	suite.db.Create(&account)
	return account
}

// Helper function to create test event
func (suite *AvailabilityRepoTestSuite) createTestEvent(accountId uuid.UUID, startsAt, endsAt time.Time) model.Event {
	event := model.Event{
		Id:          uuid.New(),
		Name:        "Test Event",
		Description: nil,
		Duration:    60,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		OwnerId:     accountId,
		Status:      "IN_DECISION",
	}
	suite.db.Create(&event)
	return event
}

// Helper function to create test availability
func (suite *AvailabilityRepoTestSuite) createTestAvailability(accountId, eventId uuid.UUID, startsAt, endsAt time.Time) model.Availability {
	availability := model.Availability{
		Id:        uuid.New(),
		AccountId: accountId,
		EventId:   eventId,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}
	suite.db.Create(&availability)
	return availability
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_NoOverlaps() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create availability completely within event range
	availability := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Availability: Jan 6-8 (within event)
		time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify availability still exists and unchanged
	var result model.Availability
	err = suite.db.Where("id = ?", availability.Id).First(&result).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availability.StartsAt, result.StartsAt)
	assert.Equal(suite.T(), availability.EndsAt, result.EndsAt)
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_LeftOverlap() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create availability that starts before event and ends within event (chevauche la date de début)
	availability := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), // Availability: Jan 3-7 (overlaps left)
		time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify availability start was adjusted to event start
	var result model.Availability
	err = suite.db.Where("id = ?", availability.Id).First(&result).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result.StartsAt)    // Should be adjusted to Jan 5
	assert.Equal(suite.T(), availability.EndsAt, result.EndsAt) // End should remain Jan 7
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_RightOverlap() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create availability that starts within event and ends after event (chevauche la date de fin)
	availability := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Availability: Jan 8-12 (overlaps right)
		time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify availability end was adjusted to event end
	var result model.Availability
	err = suite.db.Where("id = ?", availability.Id).First(&result).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availability.StartsAt, result.StartsAt) // Start should remain Jan 8
	assert.Equal(suite.T(), event.EndsAt, result.EndsAt)            // Should be adjusted to Jan 10
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_BothSidesOverlap() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create availability that encompasses entire event (chevauche la nouvelle date de fin et de début)
	availability := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Availability: Jan 1-15 (encompasses event)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify availability was adjusted to exactly match event range
	var result model.Availability
	err = suite.db.Where("id = ?", availability.Id).First(&result).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result.StartsAt) // Should be adjusted to Jan 5
	assert.Equal(suite.T(), event.EndsAt, result.EndsAt)     // Should be adjusted to Jan 10
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_CompletelyOutsideDeleted() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create availability completely before event
	availabilityBefore := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Availability: Jan 1-3 (completely before)
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC))

	// Create availability completely after event
	availabilityAfter := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC), // Availability: Jan 12-15 (completely after)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify both availabilities were deleted
	var count int64
	suite.db.Model(&model.Availability{}).Where("id = ?", availabilityBefore.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)

	suite.db.Model(&model.Availability{}).Where("id = ?", availabilityAfter.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_MixedScenarios() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Event: Jan 5-10
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// Create various types of availabilities
	avail1 := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Within event (should remain unchanged)
		time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC))

	avail2 := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), // Left overlap (should be adjusted)
		time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC))

	avail3 := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Right overlap (should be adjusted)
		time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC))

	avail4 := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Both sides overlap (should be adjusted)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	avail5 := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC), // Completely after (should be deleted)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify avail1 (within event) - unchanged
	var result1 model.Availability
	err = suite.db.Where("id = ?", avail1.Id).First(&result1).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), avail1.StartsAt, result1.StartsAt)
	assert.Equal(suite.T(), avail1.EndsAt, result1.EndsAt)

	// Verify avail2 (left overlap) - start adjusted
	var result2 model.Availability
	err = suite.db.Where("id = ?", avail2.Id).First(&result2).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result2.StartsAt)
	assert.Equal(suite.T(), avail2.EndsAt, result2.EndsAt)

	// Verify avail3 (right overlap) - end adjusted
	var result3 model.Availability
	err = suite.db.Where("id = ?", avail3.Id).First(&result3).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), avail3.StartsAt, result3.StartsAt)
	assert.Equal(suite.T(), event.EndsAt, result3.EndsAt)

	// Verify avail4 (both sides overlap) - both adjusted
	var result4 model.Availability
	err = suite.db.Where("id = ?", avail4.Id).First(&result4).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result4.StartsAt)
	assert.Equal(suite.T(), event.EndsAt, result4.EndsAt)

	// Verify avail5 (completely after) - deleted
	var count int64
	suite.db.Model(&model.Availability{}).Where("id = ?", avail5.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_DifferentEvents() {
	// Arrange
	account := suite.createTestAccount()

	// Create two different events
	event1 := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	event2 := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC))

	// Create availabilities for both events
	avail1 := suite.createTestAvailability(account.Id, event1.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Event1 availability (should be adjusted)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	avail2 := suite.createTestAvailability(account.Id, event2.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Event2 availability (should NOT be touched)
		time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC))

	// Act - only adjust availabilities for event1
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event1.Id, event1.StartsAt, event1.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify avail1 was adjusted
	var result1 model.Availability
	err = suite.db.Where("id = ?", avail1.Id).First(&result1).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event1.StartsAt, result1.StartsAt)
	assert.Equal(suite.T(), event1.EndsAt, result1.EndsAt)

	// Verify avail2 was NOT touched
	var result2 model.Availability
	err = suite.db.Where("id = ?", avail2.Id).First(&result2).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), avail2.StartsAt, result2.StartsAt)
	assert.Equal(suite.T(), avail2.EndsAt, result2.EndsAt)
}

func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_EmptyEvent() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))

	// No availabilities created

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err) // Should not error even with no availabilities
}

// Tests spécifiques pour tous les cas de chevauchement mentionnés
func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_AllOverlapCases() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), // Event: Jan 10-20
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC))

	// Cas 1: Chevauche la date de début seulement
	availLeftOverlap := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Jan 5-15 (chevauche début)
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	// Cas 2: Chevauche la date de fin seulement
	availRightOverlap := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // Jan 15-25 (chevauche fin)
		time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC))

	// Cas 3: Chevauche la nouvelle date de fin ET de début (englobe complètement)
	availBothOverlap := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Jan 1-30 (chevauche début et fin)
		time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC))

	// Cas 4: Complètement dans la plage (ne devrait pas changer)
	availWithin := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC), // Jan 12-18 (complètement dedans)
		time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC))

	// Cas 5: Complètement en dehors (devrait être supprimé)
	availOutside := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), // Jan 25-30 (complètement après)
		time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Vérifier cas 1: chevauche début - start ajusté à début d'événement
	var result1 model.Availability
	err = suite.db.Where("id = ?", availLeftOverlap.Id).First(&result1).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result1.StartsAt, "Availability qui chevauche le début devrait avoir son start ajusté")
	assert.Equal(suite.T(), availLeftOverlap.EndsAt, result1.EndsAt, "End ne devrait pas changer")

	// Vérifier cas 2: chevauche fin - end ajusté à fin d'événement
	var result2 model.Availability
	err = suite.db.Where("id = ?", availRightOverlap.Id).First(&result2).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availRightOverlap.StartsAt, result2.StartsAt, "Start ne devrait pas changer")
	assert.Equal(suite.T(), event.EndsAt, result2.EndsAt, "Availability qui chevauche la fin devrait avoir son end ajusté")

	// Vérifier cas 3: chevauche début ET fin - les deux ajustés
	var result3 model.Availability
	err = suite.db.Where("id = ?", availBothOverlap.Id).First(&result3).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, result3.StartsAt, "Start devrait être ajusté au début d'événement")
	assert.Equal(suite.T(), event.EndsAt, result3.EndsAt, "End devrait être ajusté à la fin d'événement")

	// Vérifier cas 4: complètement dedans - pas de changement
	var result4 model.Availability
	err = suite.db.Where("id = ?", availWithin.Id).First(&result4).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availWithin.StartsAt, result4.StartsAt, "Availability dans la plage ne devrait pas changer")
	assert.Equal(suite.T(), availWithin.EndsAt, result4.EndsAt, "Availability dans la plage ne devrait pas changer")

	// Vérifier cas 5: complètement dehors - supprimé
	var count int64
	suite.db.Model(&model.Availability{}).Where("id = ?", availOutside.Id).Count(&count)
	assert.Equal(suite.T(), int64(0), count, "Availability complètement en dehors devrait être supprimée")
}

// Test pour vérifier les cas limites de chevauchement
func (suite *AvailabilityRepoTestSuite) TestDeleteOutOfEventRangeAndAdjustOverlaps_EdgeCases() {
	// Arrange
	account := suite.createTestAccount()
	event := suite.createTestEvent(account.Id,
		time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC), // Event: Jan 10 12h - Jan 15 18h
		time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC))

	// Cas limite: availability qui finit exactement quand l'événement commence
	availTouchingStart := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)) // finit exactement quand event commence

	// Cas limite: availability qui commence exactement quand l'événement finit
	availTouchingEnd := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC), // commence exactement quand event finit
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC))

	// Cas limite: availability exactement à la même plage que l'événement
	availExactMatch := suite.createTestAvailability(account.Id, event.Id,
		time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC))

	// Act
	err := suite.repo.DeleteOutOfEventRangeAndAdjustOverlaps(event.Id, event.StartsAt, event.EndsAt)

	// Assert
	assert.NoError(suite.T(), err)

	// Availability qui touche le début devrait rester inchangée (elle n'est ni en dehors ni chevauchante)
	var resultTouchingStart model.Availability
	err = suite.db.Where("id = ?", availTouchingStart.Id).First(&resultTouchingStart).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availTouchingStart.StartsAt, resultTouchingStart.StartsAt)
	assert.Equal(suite.T(), availTouchingStart.EndsAt, resultTouchingStart.EndsAt)

	// Availability qui touche la fin devrait rester inchangée (elle n'est ni en dehors ni chevauchante)
	var resultTouchingEnd model.Availability
	err = suite.db.Where("id = ?", availTouchingEnd.Id).First(&resultTouchingEnd).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), availTouchingEnd.StartsAt, resultTouchingEnd.StartsAt)
	assert.Equal(suite.T(), availTouchingEnd.EndsAt, resultTouchingEnd.EndsAt)

	// Availability avec match exact devrait rester inchangée
	var resultExact model.Availability
	err = suite.db.Where("id = ?", availExactMatch.Id).First(&resultExact).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), event.StartsAt, resultExact.StartsAt)
	assert.Equal(suite.T(), event.EndsAt, resultExact.EndsAt)
}

// Run the test suite
func TestAvailabilityRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityRepoTestSuite))
}
