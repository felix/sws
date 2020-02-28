package main

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"src.userspace.com.au/sws"
)

const (
	loginURL  = "/login"
	logoutURL = "/logout"
)

func handleLogin(db sws.UserStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		if email == "" || password == "" {
			//httpError(w, 406, "bad auth")
			r = flashSet(r, flashError, "invalid credentials")
			http.Redirect(w, flashQuery(r), loginURL, http.StatusSeeOther)
			return
		}

		debug("authing email", email)

		user, err := db.GetUserByEmail(email)
		if err != nil || user == nil {
			//httpError(w, 404, err.Error())
			r = flashSet(r, flashError, "invalid user")
			http.Redirect(w, flashQuery(r), loginURL, http.StatusSeeOther)
			return
		}

		if !user.Enabled {
			debug("user", email, "is disabled")
			//httpError(w, 403, "forbidden")
			r = flashSet(r, flashError, "access denied")
			http.Redirect(w, flashQuery(r), loginURL, http.StatusSeeOther)
			return
		}

		if err := user.ValidPassword(password); err != nil {
			//httpError(w, 401, err.Error())
			r = flashSet(r, flashError, "authentication failed")
			http.Redirect(w, flashQuery(r), loginURL, http.StatusSeeOther)
			return
		}
		debug("user", email, "is authed")

		expiry := time.Now().Add(time.Hour)

		_, t, err := tokenAuth.Encode(jwt.MapClaims{
			"user_id": *user.ID,
			"exp":     expiry.Unix(),
		})
		if err != nil {
			httpError(w, 500, err.Error())
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
