package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"src.userspace.com.au/sws"
)

func handleUsers(db sws.UserStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authedUser, ok := ctx.Value("user").(*sws.User)
		if !ok {
			httpError(w, 422, "no user in context")
			return
		}
		email := chi.URLParam(r, "email")
		if email == "" {
			httpError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		user, err := db.GetUserByEmail(email)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if *authedUser.ID != *user.ID && !authedUser.Admin {
			log("failed attempt to edit user", *user.ID, *authedUser.ID)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			user.FirstName = stringPtr(r.FormValue("first_name"))
			user.LastName = stringPtr(r.FormValue("last_name"))
			user.Email = stringPtr(r.FormValue("email"))
			user.PasswordConfirm = r.FormValue("password_confirmation")
			user.Password = r.FormValue("password")
			if errs := user.Validate(); len(errs) > 0 {
				log("invalid user:", errs)
				r = flashSet(r, flashError, strings.Join(errs, "<br>"))
			} else {
				if err := db.SaveUser(user); err != nil {
					log("failed to update user:", err)
					r = flashSet(r, flashError, err.Error())
				} else {
					r = flashSet(r, flashSuccess, "successfully updated")
				}
			}
		}

		payload := newTemplateData(r)
		payload.User = user

		if err := rndr.Render(w, "user", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}
}
