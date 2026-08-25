package mail

import (
	"app/commons/constants"
	"app/config"
	model "app/db/models"
	"app/testutils"
	"errors"
	"html/template"
	"net/smtp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestMailService builds the struct directly, keeping most tests independent of package-level config.
func newTestMailService(t *testing.T) *MailService {
	t.Helper()

	s := &MailService{
		Config: config.Config{
			Origin: "https://slotfinder.test",
			Email: config.EmailConfiguration{
				Host:     "smtp.test",
				Port:     "587",
				Address:  "noreply@slotfinder.test",
				Password: "password",
			},
		},
		templates:    make(map[constants.MailTemplate]*template.Template),
		translations: make(map[string]map[constants.MailTemplate]map[string]any),
	}
	require.NoError(t, s.loadTemplates())
	require.NoError(t, s.loadTranslations())
	return s
}

func TestNewMailService_ReusesProvidedInstance(t *testing.T) {
	t.Parallel()
	existing := &MailService{}
	assert.Same(t, existing, NewMailService(existing))
}

func TestNewMailService_Nil_BuildsRealDependencies(t *testing.T) {
	t.Parallel()
	s := NewMailService(nil)
	assert.NotNil(t, s.SendMailFunc)
	assert.NotEmpty(t, s.templates)
	assert.NotEmpty(t, s.translations)
}

func TestLoadTemplates_AllTemplatesLoaded(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	for _, tmpl := range []constants.MailTemplate{
		constants.MAIL_TEMPLATE_WELCOME,
		constants.MAIL_TEMPLATE_PASSWORD_RESET,
		constants.MAIL_TEMPLATE_PASSWORD_RESET_CONFIRMATION,
		constants.MAIL_TEMPLATE_EVENT_CONFIRMATION,
		constants.MAIL_TEMPLATE_EVENT_CANCELLATION,
	} {
		_, exists := s.templates[tmpl]
		assert.True(t, exists, "expected template %s to be loaded", tmpl)
	}
}

func TestLoadTranslations_EnAndFrLoaded(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	assert.NotEmpty(t, s.translations["en"])
	assert.NotEmpty(t, s.translations["fr"])
}

func TestGetTranslations_FallbackToEnglish(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	translations := s.getTranslations(constants.MAIL_TEMPLATE_WELCOME, constants.AccountLanguage("de"))
	assert.Equal(t, s.translations[string(constants.ACCOUNT_LANGUAGE_EN)][constants.MAIL_TEMPLATE_WELCOME], translations)
}

func TestGetTranslations_UnknownTemplate(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	translations := s.getTranslations(constants.MailTemplate("does-not-exist"), constants.ACCOUNT_LANGUAGE_EN)
	assert.Empty(t, translations)
}

func TestRenderTemplate_UnknownTemplate(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	_, err := s.renderTemplate(constants.MailTemplate("does-not-exist"), nil, constants.ACCOUNT_LANGUAGE_EN)
	assert.Error(t, err)
}

func TestRenderTemplate_Success(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	html, err := s.renderTemplate(constants.MAIL_TEMPLATE_WELCOME, map[string]string{"LoginUrl": "https://slotfinder.test/login"}, constants.ACCOUNT_LANGUAGE_EN)
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
}

func TestEventUrl(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	eventId := uuid.New()

	assert.Equal(t, "https://slotfinder.test/event/"+eventId.String(), s.eventUrl(eventId))
}

func TestEventEmailCommonParams(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	event := model.Event{Name: "Team Sync"}
	eventId := uuid.New()
	startsAt := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	endsAt := startsAt.Add(time.Hour)

	params := s.eventEmailCommonParams(event, eventId, startsAt, endsAt, constants.ACCOUNT_LANGUAGE_EN, "UTC")
	assert.Equal(t, "Team Sync", params["eventName"])
	assert.Equal(t, s.eventUrl(eventId), params["eventUrl"])
	assert.NotEmpty(t, params["whenFormattedDateTime"])
}

func TestEventEmailCommonParams_InvalidTimeZoneFallsBackToUTC(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	event := model.Event{Name: "Team Sync"}

	params := s.eventEmailCommonParams(event, uuid.New(), time.Now(), time.Now().Add(time.Hour), constants.ACCOUNT_LANGUAGE_EN, "Not/AZone")
	assert.NotEmpty(t, params["whenFormattedDateTime"])
}

func TestEventEmailEnrichOptionalFields(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	username := "alice"
	ownerUsername := "bob"
	description := "Weekly sync"
	participant := model.Account{Username: &username}
	event := model.Event{Owner: model.Account{Username: &ownerUsername}, Description: &description}

	params := map[string]string{}
	s.eventEmailEnrichOptionalFields(params, participant, event)

	assert.Equal(t, "alice", params["username"])
	assert.Equal(t, "bob", params["owner"])
	assert.Equal(t, "Weekly sync", params["eventDescription"])
}

func TestEventEmailEnrichOptionalFields_NilFieldsSkipped(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	params := map[string]string{}
	s.eventEmailEnrichOptionalFields(params, model.Account{}, model.Event{})

	_, hasUsername := params["username"]
	_, hasOwner := params["owner"]
	assert.False(t, hasUsername)
	assert.False(t, hasOwner)
}

func TestSendMail_ValidationErrors(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	assert.Error(t, s.SendMail(EmailParams{To: "a@b.com", Subject: "Hi"}))
	assert.Error(t, s.SendMail(EmailParams{Template: constants.MAIL_TEMPLATE_WELCOME, Subject: "Hi"}))
	assert.Error(t, s.SendMail(EmailParams{Template: constants.MAIL_TEMPLATE_WELCOME, To: "a@b.com"}))
}

func TestSendMail_RenderError(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	err := s.SendMail(EmailParams{Template: constants.MailTemplate("unknown"), To: "a@b.com", Subject: "Hi"})
	assert.Error(t, err)
}

func TestSendMail_Success(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	var capturedTo []string
	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedTo = to
		return nil
	}

	err := s.SendMail(EmailParams{
		Template: constants.MAIL_TEMPLATE_WELCOME,
		To:       "a@b.com",
		Subject:  "Welcome",
		Language: constants.ACCOUNT_LANGUAGE_FR,
		Params:   map[string]string{"LoginUrl": "https://slotfinder.test/login"},
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{"a@b.com"}, capturedTo)
}

func TestSendMail_DefaultsLanguageAndParams(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	testutils.StubSMTP(t, &s.SendMailFunc)

	err := s.SendMail(EmailParams{Template: constants.MAIL_TEMPLATE_WELCOME, To: "a@b.com", Subject: "Welcome"})
	assert.NoError(t, err)
}

func TestSendMail_Disabled_SkipsSending(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	s.Config.Email.Disabled = true

	called := false
	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called = true
		return nil
	}

	err := s.SendMail(EmailParams{Template: constants.MAIL_TEMPLATE_WELCOME, To: "a@b.com", Subject: "Welcome"})

	assert.NoError(t, err)
	assert.False(t, called)
}

func TestSendMail_SmtpError(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errors.New("connection refused")
	}

	err := s.SendMail(EmailParams{Template: constants.MAIL_TEMPLATE_WELCOME, To: "a@b.com", Subject: "Welcome"})
	assert.Error(t, err)
}

func TestSendEventConfirmationEmail_NoEmailOrUsername_NoOp(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	// No email/username set -> should return without attempting to send (and thus without needing the SendMailFunc stub).
	s.SendEventConfirmationEmail(model.Account{}, model.Event{}, uuid.New(), uuid.New(), time.Now(), time.Now())
}

func TestSendEventCancellationEmail_NoEmailOrUsername_NoOp(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)
	s.SendEventCancellationEmail(model.Account{}, model.Event{}, uuid.New(), uuid.New(), time.Now(), time.Now())
}

func TestSendEventConfirmationEmail_Success(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	sent := make(chan struct{}, 1)
	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		sent <- struct{}{}
		return nil
	}

	email := "participant@example.com"
	username := "participant"
	ownerId := uuid.New()
	participant := model.Account{Id: ownerId, Email: &email, Username: &username, Language: constants.ACCOUNT_LANGUAGE_FR}
	event := model.Event{Name: "Sync"}

	s.SendEventConfirmationEmail(participant, event, uuid.New(), ownerId, time.Now(), time.Now().Add(time.Hour))

	select {
	case <-sent:
	case <-time.After(2 * time.Second):
		t.Fatal("expected SendMail to be called asynchronously")
	}
}

func TestSendEventCancellationEmail_Success(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	sent := make(chan struct{}, 1)
	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		sent <- struct{}{}
		return nil
	}

	email := "account@example.com"
	username := "account"
	ownerId := uuid.New()
	account := model.Account{Id: uuid.New(), Email: &email, Username: &username, Language: constants.ACCOUNT_LANGUAGE_EN}
	event := model.Event{Name: "Sync"}

	s.SendEventCancellationEmail(account, event, uuid.New(), ownerId, time.Now(), time.Now().Add(time.Hour))

	select {
	case <-sent:
	case <-time.After(2 * time.Second):
		t.Fatal("expected SendMail to be called asynchronously")
	}
}

func TestSendEventCancellationEmail_FrenchSubject(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	sent := make(chan struct{}, 1)
	s.SendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		sent <- struct{}{}
		return nil
	}

	email := "account-fr@example.com"
	username := "account"
	ownerId := uuid.New()
	account := model.Account{Id: uuid.New(), Email: &email, Username: &username, Language: constants.ACCOUNT_LANGUAGE_FR}
	event := model.Event{Name: "Sync"}

	s.SendEventCancellationEmail(account, event, uuid.New(), ownerId, time.Now(), time.Now().Add(time.Hour))

	select {
	case <-sent:
	case <-time.After(2 * time.Second):
		t.Fatal("expected SendMail to be called asynchronously")
	}
}

func TestBuildEmailMessage(t *testing.T) {
	t.Parallel()
	s := newTestMailService(t)

	msg := s.buildEmailMessage("to@example.com", "Subject Line", "<p>Body</p>")
	assert.Contains(t, string(msg), "To: to@example.com")
	assert.Contains(t, string(msg), "Subject: Subject Line")
	assert.Contains(t, string(msg), "<p>Body</p>")
}
