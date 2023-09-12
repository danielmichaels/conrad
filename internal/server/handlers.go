package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/danielmichaels/conrad/internal/render"
	"github.com/danielmichaels/conrad/internal/repository"
	"github.com/danielmichaels/conrad/internal/validator"
	"github.com/danielmichaels/conrad/internal/version"
	"github.com/go-playground/webhooks/v6/gitlab"
	"net/http"
)

func (app *Application) status(w http.ResponseWriter, r *http.Request) {
	data := Map{
		"Status":  "OK",
		"Version": version.Get(),
	}
	err := render.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	err := render.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
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
			if errors.Is(err, sql.ErrNoRows) {
				form.Validator.AddFieldError("email", "invalid credentials")
				form.Validator.AddFieldError("password", "invalid credentials")
			} else {
				app.serverError(w, r, err)
				return
			}
		}

		form.Validator.CheckField(form.Email != "", "Email", "Email is required")
		form.Validator.CheckField(user.Email != "", "Email", "Email address could not be found")

		if user.ID != 0 {
			passwordMatches, err := Matches(form.Password, user.HashedPassword)
			if err != nil {
				app.serverError(w, r, err)
				return
			}

			form.Validator.CheckField(form.Password != "", "Password", "Password is required")
			form.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
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

func (app *Application) gitlabWebhook(w http.ResponseWriter, r *http.Request) {
	hook, err := gitlab.New(gitlab.Options.Secret(app.Config.Secrets.GitlabWHVerifySecret))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	payload, err := hook.Parse(r, gitlab.CommentEvents, gitlab.MergeRequestEvents)
	if err != nil {
		if errors.Is(err, gitlab.ErrEventNotFound) {
			app.errorHookEventNotFound(w, r)
		}
		app.Logger.Error().Err(err).Send()
		return
	}
	switch payload.(type) {
	case gitlab.MergeRequestEventPayload:
		mr := payload.(gitlab.MergeRequestEventPayload)
		fmt.Printf("%+v\n", mr)
		_ = render.JSON(w, http.StatusOK, nil)
	case gitlab.CommentEventPayload:
		comment := payload.(gitlab.CommentEventPayload)
		fmt.Printf("%+v\n", comment)
		_ = render.JSON(w, http.StatusOK, nil)
	case gitlab.PushEventPayload:
		comment := payload.(gitlab.PushEventPayload)
		fmt.Printf("%+v\n", comment)
		_ = render.JSON(w, http.StatusOK, nil)
	default:
		app.Logger.Warn().Msgf("unknown hook: %v", hook)
		_ = render.JSON(w, http.StatusOK, nil)
	}
}
