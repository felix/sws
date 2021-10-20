package main

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"src.userspace.com.au/render"

	"src.userspace.com.au/sws"
)

const (
	loginURL  = "/login"
	logoutURL = "/logout"
)

func handleLogin(db sws.UserStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var user *sws.User
			r, user = authUser(db, r)
			if user != nil {
				expiry := time.Now().Add(time.Hour)

				_, t, err := tokenAuth.Encode(jwt.MapClaims{
					"user_id": *user.ID,
					"exp":     expiry.Unix(),
				})
				if err != nil {
					log("failed to encode claims:", err)
					r = flashSet(r, flashError, "internal error")
					http.Redirect(w, r, flashURL(r, "/"), http.StatusSeeOther)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "jwt",
					Value:    t,
					HttpOnly: true,
					Path:     "/",
					//Secure: true,
					Expires: expiry,
				})
				r = r.WithContext(context.WithValue(r.Context(), "user", user))
				r = flashSet(r, flashSuccess, "authenticated successfully")
				qs := r.URL.Query()
				if returnPath := qs.Get("return_to"); returnPath != "" {
					qs.Del("return_to")
					r.URL.RawQuery = qs.Encode()
					debug("redirecting to", returnPath)
					http.Redirect(w, r, flashURL(r, returnPath), http.StatusSeeOther)
				}
				http.Redirect(w, r, flashURL(r, "/sites"), http.StatusSeeOther)
			}
		}

		payload := newTemplateData(r)
		rndr.HTML(
			w, 200, payload,
			render.Template("login.tmpl"),
			render.Layout("layout.tmpl"),
		)
	}
}

func handleLogout(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			//Secure: true,
			Expires: time.Time{},
		})
		r = flashSet(r, flashSuccess, "de-authenticated successfully")
		http.Redirect(w, r, flashURL(r, "/"), http.StatusSeeOther)
	}
}

func authRedirect(w http.ResponseWriter, r *http.Request, msg string) {
	flashSet(r, flashError, msg)
	log(msg)
	qs := r.URL.Query()
	qs.Set("return_to", r.URL.Path)
	r.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, flashURL(r, loginURL), http.StatusSeeOther)
}

func authUser(db sws.UserStore, r *http.Request) (*http.Request, *sws.User) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if email == "" || password == "" {
		return flashSet(r, flashError, "invalid credentials"), nil
	}

	debug("authing email", email)

	user, err := db.GetUserByEmail(email)
	if err != nil || user == nil {
		return flashSet(r, flashError, "invalid user"), nil
	}

	if !user.Enabled {
		debug("user", email, "is disabled")
		return flashSet(r, flashError, "access denied"), nil
	}

	if err := user.ValidPassword(password); err != nil {
		return flashSet(r, flashError, "authentication failed"), nil
	}
	debug("user", email, "is authed")
	return r, user
}
