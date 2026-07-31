package lib

import (
	"app/commons/constants"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatLocalizedDate_SameDay_English(t *testing.T) {
	start := time.Date(2026, time.May, 14, 17, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.May, 14, 18, 0, 0, 0, time.UTC)

	result := FormatLocalizedDate(start, end, constants.ACCOUNT_LANGUAGE_EN)
	assert.Equal(t, "Thursday, May 14, 17:00–18:00", result)
}

func TestFormatLocalizedDate_SameDay_French(t *testing.T) {
	start := time.Date(2026, time.December, 7, 20, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.December, 7, 23, 0, 0, 0, time.UTC)

	result := FormatLocalizedDate(start, end, constants.ACCOUNT_LANGUAGE_FR)
	assert.Equal(t, "lundi 07 décembre de 20h00 à 23h00", result)
}

func TestFormatLocalizedDate_MultiDay_English(t *testing.T) {
	start := time.Date(2026, time.December, 7, 20, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.December, 8, 23, 0, 0, 0, time.UTC)

	result := FormatLocalizedDate(start, end, constants.ACCOUNT_LANGUAGE_EN)
	assert.Equal(t, "From Monday 07 December at 20:00 to Tuesday 08 December at 23:00", result)
}

func TestFormatLocalizedDate_MultiDay_French(t *testing.T) {
	start := time.Date(2026, time.December, 7, 20, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.December, 8, 23, 0, 0, 0, time.UTC)

	result := FormatLocalizedDate(start, end, constants.ACCOUNT_LANGUAGE_FR)
	assert.Equal(t, "Du lundi 07 décembre à 20h00 au mardi 08 décembre à 23h00", result)
}

func TestFormatLocalizedDate_UnknownLanguageFallsBackToEnglish(t *testing.T) {
	start := time.Date(2026, time.May, 14, 17, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.May, 14, 18, 0, 0, 0, time.UTC)

	result := FormatLocalizedDate(start, end, constants.AccountLanguage("de"))
	assert.Contains(t, result, "Thursday")
}
