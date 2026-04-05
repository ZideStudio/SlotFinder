package lib

import (
	"fmt"
	"time"

	"app/commons/constants"

	"github.com/goodsign/monday"
)

// FormatLocalizedDate formats a date time range into a localized human-friendly string
func FormatLocalizedDate(start, end time.Time, lang constants.AccountLanguage) string {
	loc := start.Location()
	endInLoc := end.In(loc)

	if sameLocalDay(start, endInLoc) {
		return formatSameDay(start, endInLoc, lang)
	}

	return formatMultiDay(start, endInLoc, lang)
}

func formatSameDay(start, end time.Time, lang constants.AccountLanguage) string {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		// "Lundi 06 décembre de 20h00 à 23h00"
		return fmt.Sprintf(
			"%s %s de %s à %s",
			formatWeekday(start, lang),
			formatDayMonth(start, lang),
			formatTime(start, lang),
			formatTime(end, lang),
		)
	default:
		// Fallback to English
		// "Thursday, May 14, 17:00–18:00"
		return fmt.Sprintf(
			"%s, %s, %s–%s",
			formatWeekday(start, constants.ACCOUNT_LANGUAGE_EN),
			formatEnglishMonthDay(start),
			formatTime(start, constants.ACCOUNT_LANGUAGE_EN),
			formatTime(end, constants.ACCOUNT_LANGUAGE_EN),
		)
	}
}

func formatMultiDay(start, end time.Time, lang constants.AccountLanguage) string {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		// "Du Lundi 06 décembre à 20h00 au Mardi 07 décembre à 23h00"
		return fmt.Sprintf(
			"Du %s %s à %s au %s %s à %s",
			formatWeekday(start, lang),
			formatDayMonth(start, lang),
			formatTime(start, lang),
			formatWeekday(end, lang),
			formatDayMonth(end, lang),
			formatTime(end, lang),
		)
	default:
		// Fallback to English
		// Keep existing multi-day English format (only same-day format was requested to change)
		// "From Monday 06 December at 20:00 to Tuesday 07 December at 23:00"
		return fmt.Sprintf(
			"From %s %s at %s to %s %s at %s",
			formatWeekday(start, constants.ACCOUNT_LANGUAGE_EN),
			formatDayMonth(start, constants.ACCOUNT_LANGUAGE_EN),
			formatTime(start, constants.ACCOUNT_LANGUAGE_EN),
			formatWeekday(end, constants.ACCOUNT_LANGUAGE_EN),
			formatDayMonth(end, constants.ACCOUNT_LANGUAGE_EN),
			formatTime(end, constants.ACCOUNT_LANGUAGE_EN),
		)
	}
}

func mondayLocale(lang constants.AccountLanguage) monday.Locale {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		return monday.LocaleFrFR
	default:
		return monday.LocaleEnUS
	}
}

func sameLocalDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func formatWeekday(t time.Time, lang constants.AccountLanguage) string {
	return monday.Format(t, "Monday", mondayLocale(lang))
}

func formatDayMonth(t time.Time, lang constants.AccountLanguage) string {
	return monday.Format(t, "02 January", mondayLocale(lang))
}

func formatEnglishMonthDay(t time.Time) string {
	return monday.Format(t, "January 2", mondayLocale(constants.ACCOUNT_LANGUAGE_EN))
}

func formatTime(t time.Time, lang constants.AccountLanguage) string {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		return t.Format("15h04")
	default:
		// "20:00"
		return t.Format("15:04")
	}
}
