package mail

import (
	"app/commons/constants"
	model "app/db/models"
	"fmt"
	"maps"
	"time"

	"github.com/rs/zerolog/log"
)

// formatTimeForLocale formats time according to language
func formatTimeForLocale(t time.Time, lang constants.AccountLanguage) string {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		if t.Minute() == 0 {
			return fmt.Sprintf("%dh", t.Hour())
		}
		return fmt.Sprintf("%dh%02d", t.Hour(), t.Minute())
	default: // EN
		if t.Minute() == 0 {
			return t.Format("3PM")
		}
		return t.Format("3:04PM")
	}
}

// formatDateTimeRangeForLocale formats a date/time range according to language
func formatDateTimeRangeForLocale(start, end time.Time, lang constants.AccountLanguage) string {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		startDate := start.Format("2 January 2006")
		endDate := end.Format("2 January 2006")
		startTime := formatTimeForLocale(start, lang)
		endTime := formatTimeForLocale(end, lang)

		// French month names
		startDate = translateMonthToFrench(startDate)
		endDate = translateMonthToFrench(endDate)

		if start.Format("2006-01-02") == end.Format("2006-01-02") {
			// Same date: "10 janvier 2025 10h à 11h15"
			return fmt.Sprintf("%s %s à %s", startDate, startTime, endTime)
		} else {
			// Different dates: "10 janvier 2025 10h à 11 janvier 2025 11h15"
			return fmt.Sprintf("%s %s à %s %s", startDate, startTime, endDate, endTime)
		}
	default: // EN
		startDate := start.Format("January 2, 2006")
		endDate := end.Format("January 2, 2006")
		startTime := formatTimeForLocale(start, lang)
		endTime := formatTimeForLocale(end, lang)

		if start.Format("2006-01-02") == end.Format("2006-01-02") {
			// Same date: "January 10, 2026 9AM - 10:15AM"
			return fmt.Sprintf("%s %s - %s", startDate, startTime, endTime)
		} else {
			// Different dates: "January 10, 2026 9AM - January 10, 2026 10:15AM"
			return fmt.Sprintf("%s %s - %s %s", startDate, startTime, endDate, endTime)
		}
	}
}

// translateMonthToFrench replaces English month names with French ones
func translateMonthToFrench(dateStr string) string {
	months := map[string]string{
		"January":   "janvier",
		"February":  "février",
		"March":     "mars",
		"April":     "avril",
		"May":       "mai",
		"June":      "juin",
		"July":      "juillet",
		"August":    "août",
		"September": "septembre",
		"October":   "octobre",
		"November":  "novembre",
		"December":  "décembre",
	}

	for en, fr := range months {
		for i := 0; i <= len(dateStr)-len(en); i++ {
			if i+len(en) <= len(dateStr) && dateStr[i:i+len(en)] == en {
				dateStr = dateStr[:i] + fr + dateStr[i+len(en):]
				break
			}
		}
	}
	return dateStr
}

// formatDateRangeForLocale formats a date range for event periods
func formatDateRangeForLocale(start, end time.Time, lang constants.AccountLanguage) (string, string) {
	switch lang {
	case constants.ACCOUNT_LANGUAGE_FR:
		startStr := start.Format("2 January 2006")
		endStr := end.Format("2 January 2006")
		return translateMonthToFrench(startStr), translateMonthToFrench(endStr)
	default: // EN
		return start.Format("January 2, 2006"), end.Format("January 2, 2006")
	}
}

// SendEventNotificationEmails sends notification emails to participants based on template type
func (s *MailService) SendEventNotificationEmails(templateType constants.MailTemplate, event *model.Event, targetSlot *model.Slot) {
	// Prepare base email parameters
	ownerName := "Event Organizer"
	if event.Owner.UserName != nil {
		ownerName = *event.Owner.UserName
	}

	baseParams := map[string]string{
		"EventName": event.Name,
		"Owner":     ownerName,
		"EventUrl":  fmt.Sprintf("%s/event/%s", s.Config.Origin, event.Id.String()),
	}

	if event.Description != nil {
		baseParams["EventDescription"] = *event.Description
		baseParams["Description"] = *event.Description
	}

	// Date and time information will be added per participant based on their language

	// Determine subject and whether to exclude owner
	var subject string
	excludeOwner := false

	switch templateType {
	case constants.MAIL_TEMPLATE_EVENT_CONFIRMATION:
		subject = fmt.Sprintf("Event Confirmed: %s", event.Name)
	case constants.MAIL_TEMPLATE_EVENT_CANCELLATION:
		subject = fmt.Sprintf("Event Cancelled: %s", event.Name)
		excludeOwner = true // Owner should not receive cancellation emails
	default:
		log.Error().
			Str("templateType", string(templateType)).
			Str("eventId", event.Id.String()).
			Msg("Unknown template type for event notification email")
		return
	}

	// Send email to each participant
	for _, accountEvent := range event.AccountEvents {
		if accountEvent.Account.Email == nil {
			continue
		}

		// Skip owner for cancellation emails
		if excludeOwner && accountEvent.Account.Id == event.OwnerId {
			continue
		}

		participantParams := make(map[string]string)
		maps.Copy(participantParams, baseParams)

		if accountEvent.Account.Id == event.OwnerId {
			participantParams["IsOwner"] = "true"
		} else {
			participantParams["IsOwner"] = "false"
		}

		if accountEvent.Account.UserName != nil {
			participantParams["Username"] = *accountEvent.Account.UserName
		}

		// Add localized date/time information based on participant's language
		userLang := accountEvent.Account.Language
		if targetSlot != nil {
			// Format the confirmed/cancelled slot date and time
			formattedDateTime := formatDateTimeRangeForLocale(targetSlot.StartsAt, targetSlot.EndsAt, userLang)
			participantParams["FormattedDateTime"] = formattedDateTime

			// For cancellation emails, also add cancelled slot fields
			if templateType == constants.MAIL_TEMPLATE_EVENT_CANCELLATION {
				participantParams["CancelledFormattedDateTime"] = formattedDateTime
			}
		} else {
			// Format the event date range
			startStr, endStr := formatDateRangeForLocale(event.StartsAt, event.EndsAt, userLang)
			participantParams["DateRangeStart"] = startStr
			participantParams["DateRangeEnd"] = endStr
		}

		err := s.SendMail(EmailParams{
			Template: templateType,
			To:       *accountEvent.Account.Email,
			Subject:  subject,
			Params:   participantParams,
		})

		if err != nil {
			log.Error().
				Err(err).
				Str("eventId", event.Id.String()).
				Str("participantEmail", *accountEvent.Account.Email).
				Str("templateType", string(templateType)).
				Msg("Failed to send event notification email to participant")
		} else {
			log.Info().
				Str("eventId", event.Id.String()).
				Str("participantEmail", *accountEvent.Account.Email).
				Str("templateType", string(templateType)).
				Msg("Event notification email sent successfully to participant")
		}
	}
}
