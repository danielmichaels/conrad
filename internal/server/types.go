package server

import (
	"database/sql"
	"github.com/danielmichaels/conrad/internal/config"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/smtp"
	"github.com/danielmichaels/conrad/internal/validator"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"strings"
	"sync"
)

type Application struct {
	Config   *config.Conf
	Logger   zerolog.Logger
	wg       sync.WaitGroup
	Mailer   *smtp.Mailer
	Db       *repository.Queries
	Tx       *sql.DB
	Sessions *sessions.CookieStore
}

type notificationForm struct {
	Name          string   `form:"Name"`
	Enabled       int64    `form:"Enabled"`
	IgnoreDrafts  int64    `form:"IgnoreDrafts"`
	RemindAuthors int64    `form:"RemindAuthors"`
	MinAge        int64    `form:"MinAge"`
	MinStaleness  int64    `form:"MinStaleness"`
	IgnoreTerms   string   `form:"IgnoreTerms"`
	IgnoreLabels  string   `form:"IgnoreLabels"`
	RequireLabels int64    `form:"RequireLabels"`
	Days          []string `form:"Days"`
	// schedule
	ScheduleTimes []string `form:"ScheduleTimes"`
	// mattermost specific
	WebhookURL string              `form:"WebhookURL"`
	Channel    string              `form:"Channel"`
	Validator  validator.Validator `form:"-"`
}

func (f *notificationForm) fill(n *repository.GetNotificationByIDRow) {
	f.Name = n.Name
	f.Enabled = n.Enabled
	f.IgnoreDrafts = n.IgnoreDrafts
	f.RemindAuthors = n.RemindAuthors
	f.MinAge = n.MinAge
	f.MinStaleness = n.MinStaleness
	f.IgnoreTerms = n.IgnoreTerms.String
	f.IgnoreLabels = n.IgnoreLabels.String
	f.RequireLabels = n.RequireLabels
	f.ScheduleTimes = strings.Split(n.ScheduledTime.String, ",")
	f.WebhookURL = n.WebhookUrl.String
	f.Channel = n.MattermostChannel.String
	f.Days = strings.Split(n.Days, ",")
}
