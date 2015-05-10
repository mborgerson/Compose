// Copyright (C) 2015  Matt Borgerson
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
    "errors"
    "net/http"
    "github.com/zenazn/goji/web"
)

const (
    CookieName = "session_token"
)

// MakeRestrictedHttpHandler creates a wrapper that requires the user to be
// logged in to access the handler.
func MakeRestrictedHttpHandler(handler func(web.C, http.ResponseWriter, *http.Request)) (func(web.C, http.ResponseWriter, *http.Request)) {
    return func(c web.C, w http.ResponseWriter, r *http.Request) {
        // Get cookie
        cookie, err := r.Cookie(CookieName)
        if err == nil && IsSessionTokenValid(cookie.Value) {
            // Valid session. Continue to handler.
            handler(c, w, r)
        } else {
            // No. Redirect to login page.
            http.Redirect(w, r, "/login", http.StatusUnauthorized)
        }
    }
}

// Login will create a new session, given an e-mail and a password.
func Login(email, password string) (*Session, error) {
    // Lookup user
    user, err := FindUserByEmail(email)

    if err != nil {
        // Error occurred. Probably bad email.
        return nil, err
    }

    // Check password
    if !user.TestPassword(password) {
        return nil, errors.New("Invalid password")
    }

    // Create session
    session, err := CreateSession(user)
    if err != nil {
        return nil, err
    }
    _, err = session.Save()
    if err != nil {
        return nil, err
    }

    return session, nil
}

// LoginHandler is the handler for the login page.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Already logged in?
    c, err := r.Cookie(CookieName)
    if err == nil && IsSessionTokenValid(c.Value) {
        // Yes. Redirect to admin page.
        http.Redirect(w, r, "/admin", http.StatusSeeOther)
        return
    }

    if r.Method == "GET" {
        // Want the login page
        err := AdminTemplates.ExecuteTemplate(w, "login.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else if r.Method == "POST" {
        // Trying to login
        r.ParseForm()
        session, err := Login(r.FormValue("email"), r.FormValue("password"))

        if err == nil {
            // Send cookie 
            http.SetCookie(w, &http.Cookie{Name:  CookieName,
                                           Value: session.Token,
                                           Path:  "/"})

            // Continue to admin page.
            http.Redirect(w, r, "/admin", http.StatusSeeOther)
        } else {
            // Bad credentials!
            http.Redirect(w, r, "/login", http.StatusUnauthorized)
        }
    } else {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
    }
}

// LogoutHandler is the handler for the logout page.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Already logged in?
    c, err := r.Cookie(CookieName)
    if err == nil {
        session, err := FindSessionByToken(c.Value)
        if err != nil {
            session.Destroy()
        }
    }

    // Send cookie 
    http.SetCookie(w, &http.Cookie{Name:  CookieName,
                                   Value: "",
                                   Path:  "/"})
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}