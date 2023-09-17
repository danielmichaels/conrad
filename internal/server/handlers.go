package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/danielmichaels/conrad/internal/hooks"
	"github.com/danielmichaels/conrad/internal/render"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"strconv"
)

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

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Email     string              `form:"Email"`
		Password  string              `form:"Password"`
		Validator validator.Validator `form:"-"`
	}

	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data["Form"] = form

		err := render.Page(w, http.StatusOK, data, "pages/signup.tmpl")
		if err != nil {
			app.serverError(w, r, err)
		}

	case http.MethodPost:
		err := render.DecodePostForm(r, &form)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		form.Validator.CheckField(form.Email != "", "Email", "Email is required")
		form.Validator.CheckField(validator.Matches(form.Email, validator.RgxEmail), "Email", "Must be a valid email address")
		form.Validator.CheckField(form.Password != "", "Password", "Password is required")
		form.Validator.CheckField(len(form.Password) >= 8, "Password", "Password is too short")
		form.Validator.CheckField(len(form.Password) <= 72, "Password", "Password is too long")

		if form.Validator.HasErrors() {
			data := app.newTemplateData(r)
			data["Form"] = form

			err := render.Page(w, http.StatusUnprocessableEntity, data, "pages/signup.tmpl")
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

		id, err := app.Db.InsertNewUser(context.Background(), repository.InsertNewUserParams{
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

		err := render.Page(w, http.StatusOK, data, "pages/login.tmpl")
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

			err := render.Page(w, http.StatusUnprocessableEntity, data, "pages/login.tmpl")
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

func (app *Application) clients(w http.ResponseWriter, r *http.Request) {
	var form struct {
		Name                string              `form:"Name"`
		WebhookURL          string              `form:"WebhookURL"`
		GitLabURL           string              `form:"GitLabURL"`
		ClientToken         string              `form:"ClientToken"`
		TrackedRepositories []string            `form:"TrackedRepositories"`
		Validator           validator.Validator `form:"-"`
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

		glab, err := hooks.NewGitlab(form.ClientToken, form.GitLabURL, true)
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
		form.Validator.CheckField(validator.IsURL(form.WebhookURL), "WebhookURL", "Value is not a valid URL")
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
			WebhookUrl:  form.WebhookURL,
		})
		if err != nil {
			app.Logger.Error().Err(err).Msg("new_client_insert_error")
			app.serverError(w, r, err)
			return
		}

		for _, p := range projects {
			err := app.Db.InsertClientRepo(context.Background(), repository.InsertClientRepoParams{
				Name:     p.Name,
				RepoID:   int64(p.ID),
				ClientID: id,
			})
			if err != nil {
				app.Logger.Error().Err(err).Interface("data", Map{
					"repo_id": p.ID,
				}).Send()
			}
		}
		if err != nil {
			// fixme
			app.serverError(w, r, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/clients/%d", id), http.StatusTemporaryRedirect)
	}
}
func (app *Application) clientDetails(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		// fixme
		app.serverError(w, r, err)
		return
	}
	client, err := app.Db.GetClientById(context.Background(), int64(id))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.Logger.Error().Err(err).Msg("repos_not_found")
			app.serverError(w, r, err)
			return
		}
		app.notFound(w, r)
		return
	}

	data := app.newTemplateData(r)
	data["Clients"] = client
	data["FormURL"] = fmt.Sprintf("/clients/%d", id)
	err = render.Page(w, http.StatusOK, data, "pages/client.tmpl")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

//
//func (app *Application) gitlabWebhook(w http.ResponseWriter, r *http.Request) {
//	hook, err := gitlab.New(gitlab.Options.Secret(app.Config.Secrets.GitlabWHVerifySecret))
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//	payload, err := hook.Parse(r, gitlab.CommentEvents, gitlab.MergeRequestEvents)
//	if err != nil {
//		if errors.Is(err, gitlab.ErrEventNotFound) {
//			app.errorHookEventNotFound(w, r)
//		}
//		app.Logger.Error().Err(err).Send()
//		return
//	}
//	switch payload.(type) {
//	case gitlab.MergeRequestEventPayload:
//		mr := payload.(gitlab.MergeRequestEventPayload)
//		fmt.Printf("%+v\n", mr)
//		_ = render.JSON(w, http.StatusOK, nil)
//	case gitlab.CommentEventPayload:
//		comment := payload.(gitlab.CommentEventPayload)
//		fmt.Printf("%+v\n", comment)
//		_ = render.JSON(w, http.StatusOK, nil)
//	case gitlab.PushEventPayload:
//		comment := payload.(gitlab.PushEventPayload)
//		fmt.Printf("%+v\n", comment)
//		_ = render.JSON(w, http.StatusOK, nil)
//	default:
//		app.Logger.Warn().Msgf("unknown hook: %v", hook)
//		_ = render.JSON(w, http.StatusOK, nil)
//	}
//}
