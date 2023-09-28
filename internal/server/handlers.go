package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/danielmichaels/conrad/internal/providers"
	"github.com/danielmichaels/conrad/internal/render"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var form struct {
		Email     string              `form:"Email"`
		Password  string              `form:"Password"`
		Validator validator.Validator `form:"-"`
	}
	data := app.newTemplateData(r)
	data["Form"] = form
	switch r.Method {
	case http.MethodGet:
		err := render.Page(w, http.StatusOK, data, "auth/signup.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		existingUser, err := app.Db.GetUserByEmail(ctx, form.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				app.serverError(w, r, err)
				return
			}
		}
		form.Validator.CheckField(existingUser.ID == 0, "Email", "Email is already in use")
		form.Validator.CheckField(form.Email != "", "Email", "Email is required")
		form.Validator.CheckField(validator.Matches(form.Email, validator.RgxEmail), "Email", "Must be a valid email address")
		form.Validator.CheckField(form.Password != "", "Password", "Password is required")
		form.Validator.CheckField(len(form.Password) >= 8, "Password", "Password is too short")
		form.Validator.CheckField(len(form.Password) <= 72, "Password", "Password is too long")
		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form

			err := render.Page(w, http.StatusUnprocessableEntity, data, "auth/signup.tmpl")
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}

		hashedPassword, err := Hash(form.Password)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		id, err := app.Db.InsertNewUser(ctx, repository.InsertNewUserParams{
			Email:          form.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		session, err := app.Sessions.Get(r, "session")
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		session.Values["userID"] = id

		err = session.Save(r, w)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Email     string              `form:"Email"`
		Password  string              `form:"Password"`
		Validator validator.Validator `form:"-"`
	}
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data["Form"] = form

		err := render.Page(w, http.StatusOK, data, "auth/login.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		user, err := app.Db.GetUserByEmail(context.Background(), form.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				app.serverError(w, r, err)
				return
			}
			form.Validator.AddFieldError("Email", "Invalid credentials")
			form.Validator.AddFieldError("Password", "Invalid credentials")
		}

		form.Validator.CheckField(form.Email != "", "Email", "Email is required")
		if user.ID != 0 {
			passwordMatches, err := Matches(form.Password, user.HashedPassword)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			form.Validator.CheckField(form.Password != "", "Password", "Password is required")
			form.Validator.CheckField(passwordMatches, "Password", "Invalid credentials")
		}

		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form

			err := render.Page(w, http.StatusUnprocessableEntity, data, "auth/login.tmpl")
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}

		session, err := app.Sessions.Get(r, "session")
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		session.Values["userID"] = user.ID

		redirectPath, ok := session.Values["redirectPathAfterLogin"].(string)
		if ok {
			delete(session.Values, "redirectPathAfterLogin")
		} else {
			redirectPath = "/dashboard"
		}

		err = session.Save(r, w)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	}
}
func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) {
	session, err := app.Sessions.Get(r, "session")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	delete(session.Values, "userID")
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	err := render.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *Application) dashboard(w http.ResponseWriter, r *http.Request) {
	clients, err := app.Db.GetAllClients(context.Background())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.Logger.Error().Err(err).Msg("provider-not-found")
			app.serverError(w, r, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data["Clients"] = clients
	data["FormURL"] = "/dashboard/clients"
	err = render.Page(w, http.StatusOK, data, "pages/dashboard.tmpl")
	if err != nil {
		app.Logger.Error().Err(err).Msg("provider-not-found")
		app.serverError(w, r, err)
		return
	}
}
func (app *Application) clients(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Name        string              `form:"Name"`
		GitLabURL   string              `form:"GitLabURL"`
		ClientToken string              `form:"ClientToken"`
		Insecure    string              `form:"Insecure"`
		Validator   validator.Validator `form:"-"`
	}

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data["Form"] = form
		data["FormURL"] = "/dashboard/clients"

		err := render.Page(w, http.StatusOK, data, "pages/create-clients.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		form.GitLabURL = formatURL(form.GitLabURL)

		insecure := false
		if form.Insecure == "on" {
			insecure = true
		}
		glab, err := providers.NewGitlab(
			form.ClientToken,
			form.GitLabURL,
			insecure,
			providers.GitlabClientDefaultTimeout,
		)
		if err != nil {
			app.Logger.Error().Err(err).Msg("gitlab_client_err")
			app.serverError(w, r, err)
		}
		projects, _, err := glab.Client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Membership: Ptr(true),
		})
		if err != nil {
			var ge *gitlab.ErrorResponse
			if errors.As(err, &ge) {
				form.Validator.AddFieldError("ClientToken", fmt.Sprintf("Invalid token: %q", ge.Response.Status))
			}
			app.Logger.Error().Err(err).
				Interface("gitlab_error", Map{"message": ge.Message, "status": ge.Response.Status}).
				Msg("gitlab-auth-failure")
		}

		form.Validator.CheckField(form.Name != "", "Name", "Name is required")
		form.Validator.CheckField(validator.IsURL(form.GitLabURL), "GitLabURL", "Value is not a valid URL")
		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form
			data["FormURL"] = "/dashboard/clients"

			err := render.Page(w, http.StatusUnprocessableEntity, data, "pages/create-clients.tmpl")
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}
		id, err := app.Db.InsertNewClient(context.Background(), repository.InsertNewClientParams{
			Name:        form.Name,
			CreatedBy:   contextGetAuthenticatedUser(r).ID,
			AccessToken: form.ClientToken,
			GitlabUrl:   form.GitLabURL,
			Insecure:    isOnCheckbox(form.Insecure),
		})
		if err != nil {
			app.Logger.Error().Err(err).Msg("new_client_insert_error")
			app.serverError(w, r, err)
			return
		}

		for _, p := range projects {
			err := app.Db.InsertClientRepo(context.Background(), repository.InsertClientRepoParams{
				Name:       p.Name,
				RepoID:     int64(p.ID),
				ClientID:   id,
				RepoWebUrl: p.WebURL,
			})
			if err != nil {
				app.Logger.Error().Err(err).Interface("data", Map{
					"repo_id": p.ID,
				}).Send()
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/dashboard/clients/%d", id), http.StatusSeeOther)
	}
}
func (app *Application) clientHome(w http.ResponseWriter, r *http.Request) {
	var form struct {
		RepoID    []int64             `form:"RepoID"`
		Validator validator.Validator `form:"-"`
	}
	id, err := urlIDParam(r, nil)
	if err != nil {
		app.notFound(w, r)
		return
	}
	client, err := app.Db.GetClientById(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.serverError(w, r, err)
			return
		}
		app.notFound(w, r)
		return
	}
	repos, err := app.Db.GetAllClientRepos(context.Background(), id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.serverError(w, r, err)
			return
		}
		app.notFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data["Form"] = form
		data["Clients"] = client
		data["Repos"] = repos
		data["FormURL"] = fmt.Sprintf("/dashboard/clients/%d", id)

		err = render.Page(w, http.StatusOK, data, "pages/client-home.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.badRequest(w, r, err)
			return
		}
		found := make([]int64, 0, len(form.RepoID))
		found = append(found, form.RepoID...)
		for _, v := range repos {
			var tracked int64 = 0
			if slices.Contains(found, v.RepoID) {
				tracked = 1
			}
			err := app.Db.UpdateTrackedRepoStatus(context.Background(), repository.UpdateTrackedRepoStatusParams{
				Tracked: tracked,
				RepoID:  v.RepoID,
			})
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
		data := app.newTemplateData(r)
		data["Form"] = form
		data["Clients"] = client
		data["Repos"] = repos
		data["FormURL"] = fmt.Sprintf("/dashboard/clients/%d", id)
		http.Redirect(w, r, fmt.Sprintf("/dashboard/clients/%d", id), http.StatusSeeOther)
	}
}

func convStrBoolToBool(v string) bool {
	return v == "true"
}
func (app *Application) clientGitlab(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var form struct {
		Name        string              `form:"Name"`
		GitLabURL   string              `form:"GitLabURL"`
		ClientToken string              `form:"ClientToken"`
		Insecure    string              `form:"Insecure"`
		Validator   validator.Validator `form:"-"`
	}

	id, err := urlIDParam(r, nil)
	if err != nil {
		app.notFound(w, r)
		return
	}
	pageURL := fmt.Sprintf("/dashboard/clients/%d/gitlab", id)
	client, err := app.Db.GetClientById(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.Logger.Error().Err(err).Msg("client_by_id")
			app.serverError(w, r, err)
			return
		}
		app.notFound(w, r)
		return
	}

	form.ClientToken = client.AccessToken
	form.Insecure = isOnCheckbox(client.Insecure)
	form.Name = client.Name
	form.GitLabURL = client.GitlabUrl

	data := app.newTemplateData(r)
	data["Form"] = form
	data["Clients"] = client
	data["FormURL"] = pageURL

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data["Form"] = form
		data["Clients"] = client
		data["FormURL"] = pageURL
		err = render.Page(w, http.StatusOK, data, "pages/client-gitlab.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.badRequest(w, r, err)
			return
		}
		form.Validator.CheckField(form.Name != "", "Name", "Name is required")
		form.Validator.CheckField(validator.IsURL(form.GitLabURL), "GitLabURL", "Value is not a valid URL")
		form.GitLabURL = formatURL(form.GitLabURL)

		glab, err := providers.NewGitlab(
			form.ClientToken,
			form.GitLabURL,
			convStrBoolToBool(form.Insecure),
			providers.GitlabClientDefaultTimeout,
		)
		if err != nil {
			app.Logger.Error().Err(err).Msg("gitlab_client_err")
			form.Validator.AddFieldError("ClientToken", "Unable to contact Gitlab")
		}
		projects, _, err := glab.Client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Membership: Ptr(true),
		})
		if err != nil {
			form.Validator.AddFieldError("ClientToken", "Unable to contact Gitlab")
			form.Validator.AddFieldError("GitLabURL", "Unable to contact Gitlab")
			app.Logger.Error().Err(err).Msg("gitlab-auth-failure")
		}
		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form
			data["Clients"] = client
			data["FormURL"] = pageURL
			//err := render.Page(w, http.StatusUnprocessableEntity, data, "pages/client-gitlab.tmpl")
			err := render.Page(w, http.StatusOK, data, "pages/client-gitlab.tmpl")
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}
		id, err := app.Db.UpdateExistingClient(ctx, repository.UpdateExistingClientParams{
			Name:        form.Name,
			CreatedBy:   contextGetAuthenticatedUser(r).ID,
			AccessToken: form.ClientToken,
			GitlabUrl:   form.GitLabURL,
			Insecure:    isOnCheckbox(form.Insecure),
			ID:          id,
		})
		if err != nil {
			app.Logger.Error().Err(err).Msg("new_client_insert_error")
			app.serverError(w, r, err)
			return
		}

		// Check if any of the existing tracked repositories have changed
		for _, p := range projects {
			params := repository.UpsertClientRepoParams{
				Name:       p.Name,
				RepoID:     int64(p.ID),
				ClientID:   id,
				RepoWebUrl: p.WebURL,
			}
			err := app.Db.UpsertClientRepo(ctx, params)
			if err != nil {
				app.Logger.Error().Err(err).Interface("data", Map{"repo_id": p.ID}).Send()
				app.serverError(w, r, err)
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/dashboard/clients/%d", id), http.StatusSeeOther)
	}
}

func (app *Application) clientNotifications(w http.ResponseWriter, r *http.Request) {
	var form struct {
		NotificationID []int64             `form:"NotificationID"`
		Validator      validator.Validator `form:"-"`
	}
	ctx := context.Background()
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	client, err := app.Db.GetClientById(ctx, int64(id))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.Logger.Error().Err(err).Msg("client_by_id")
			app.serverError(w, r, err)
			return
		}
		app.notFound(w, r)
		return
	}
	notifications, err := app.Db.GetAllNotifications(ctx)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	pageURL := fmt.Sprintf("/dashboard/clients/%d/notifications", id)
	data := app.newTemplateData(r)
	data["Form"] = form
	data["FormURL"] = pageURL
	data["Clients"] = client
	data["Notifications"] = notifications
	switch r.Method {
	case http.MethodGet:
		err = render.Page(w, http.StatusOK, data, "pages/client-notifications.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.badRequest(w, r, err)
			return
		}
		found := make([]int64, 0, len(form.NotificationID))
		found = append(found, form.NotificationID...)
		for _, v := range notifications {
			var enabled int64 = 0
			if slices.Contains(found, v.ID) {
				enabled = 1
			}
			err := app.Db.UpdateEnabledNotificationStatus(ctx, repository.UpdateEnabledNotificationStatusParams{
				Enabled: enabled,
				ID:      v.ID,
			})
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}
		data := app.newTemplateData(r)
		data["Form"] = form
		http.Redirect(w, r, pageURL, http.StatusSeeOther)
	}
}

func (app *Application) mattermostNotification(w http.ResponseWriter, r *http.Request) {
	var form struct {
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
	ctx := context.Background()
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	pageURL := fmt.Sprintf("/dashboard/clients/%d/notifications/mattermost", id)
	backURL := fmt.Sprintf("/dashboard/clients/%d/notifications", id)

	data := app.newTemplateData(r)
	data["Form"] = form
	data["FormURL"] = pageURL
	data["BackURL"] = backURL
	switch r.Method {
	case http.MethodGet:
		err = render.Page(w, http.StatusOK, data, "pages/create-notifications.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}
	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.badRequest(w, r, err)
			return
		}

		tx, err := app.Tx.Begin()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				app.Logger.Error().Err(err).Msg("defer_transaction_err")
			}
		}(tx)

		qtx := app.Db.WithTx(tx)
		notification, err := qtx.InsertNotification(ctx, repository.InsertNotificationParams{
			Enabled:        1,
			Name:           form.Name,
			ClientID:       int64(id),
			IgnoreApproved: 0,
			IgnoreDrafts:   0,
			RemindAuthors:  0,
			MinAge:         0,
			MinStaleness:   0,
			IgnoreTerms:    sql.NullString{},
			IgnoreLabels:   sql.NullString{},
			RequireLabels:  0,
			Days:           strings.Join(form.Days, ","),
		})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		times, err := qtx.UpsertNotificationTimes(ctx, repository.UpsertNotificationTimesParams{
			NotificationID: notification,
			ScheduledTime:  strings.Join(form.ScheduleTimes, ","),
		})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		_, err = qtx.UpsertNotificationMattermost(ctx, repository.UpsertNotificationMattermostParams{
			NotificationID:    times,
			MattermostChannel: form.Channel,
			WebhookUrl:        form.WebhookURL,
		})
		if err != nil {
			return
		}
		err = tx.Commit()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data := app.newTemplateData(r)
		data["Form"] = form

		http.Redirect(w, r, backURL, http.StatusSeeOther)
	}
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

func (app *Application) notificationDetail(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := urlIDParam(r, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	nid, err := urlIDParam(r, Ptr("nid"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	pageURL := fmt.Sprintf("/dashboard/clients/%d/notifications/%d", id, nid)
	backURL := fmt.Sprintf("/dashboard/clients/%d/notifications", id)

	var form notificationForm
	notification, err := app.Db.GetNotificationByID(ctx, nid)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		form.fill(&notification)
		data := app.newTemplateData(r)
		data["Form"] = form
		data["FormURL"] = pageURL
		data["BackURL"] = backURL
		err = render.Page(w, http.StatusOK, data, "pages/create-notifications.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}
	case http.MethodDelete:
		fmt.Println("DELETE")
		return
	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.Logger.Error().Err(err).Msg("post_decode_error")
			app.badRequest(w, r, err)
			return
		}
		form.Validator.CheckField(form.Name != "", "Name", "Name is required")
		form.Validator.CheckField(form.Channel != "", "Channel", "Channel is required")
		form.Validator.CheckField(validator.IsURL(form.WebhookURL), "WebhookURL", "Value is not a valid URL")
		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form
			err := render.Page(w, http.StatusUnprocessableEntity, data, "pages/create-notifications.tmpl")
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}

		tx, err := app.Tx.Begin()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				app.Logger.Error().Err(err).Msg("defer_transaction_err")
			}
		}(tx)
		qtx := app.Db.WithTx(tx)
		notification, err := qtx.UpsertNotification(ctx, repository.UpsertNotificationParams{
			//notification, err := app.Db.UpsertNotification(ctx, repository.UpsertNotificationParams{
			ID:             nid,
			Enabled:        1,
			Name:           form.Name,
			ClientID:       id,
			IgnoreApproved: 0,
			IgnoreDrafts:   0,
			RemindAuthors:  0,
			MinAge:         0,
			MinStaleness:   0,
			IgnoreTerms:    sql.NullString{},
			IgnoreLabels:   sql.NullString{},
			RequireLabels:  0,
			Days:           strings.Join(form.Days, ","),
		})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		notificationId := notification
		fmt.Println("noID", notificationId)
		times, err := qtx.UpsertNotificationTimes(ctx, repository.UpsertNotificationTimesParams{
			//times, err := app.Db.UpdateNotificationTimes(ctx, repository.UpdateNotificationTimesParams{
			NotificationID: notificationId,
			ScheduledTime:  strings.Join(form.ScheduleTimes, ","),
		})
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.serverError(w, r, err)
			return
		}
		fmt.Println("noTimesID", times)
		_, err = qtx.UpsertNotificationMattermost(ctx, repository.UpsertNotificationMattermostParams{
			//_, err = app.Db.UpdateNotificationMattermost(ctx, repository.UpdateNotificationMattermostParams{
			NotificationID:    notificationId,
			MattermostChannel: form.Channel,
			WebhookUrl:        form.WebhookURL,
		})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		err = tx.Commit()
		if err != nil {
			app.Logger.Error().Err(err).Send()
			app.serverError(w, r, err)
			return
		}
		data := app.newTemplateData(r)
		data["Form"] = form

		http.Redirect(w, r, backURL, http.StatusSeeOther)
	}
}
