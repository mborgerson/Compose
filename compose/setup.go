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
    "net/http"
)

// SetupAlreadyComplete determines if setup has been completed or not. For now,
// consider setup complete if there is at least one user in the database.
func SetupAlreadyComplete() (bool) {
    count, _ := GetDatabaseHandle().C("users").Find(nil).Count()
    return count > 0
}

// Setup initializes the database to a usable state. Essentially, it just
// creates a user that can login.
func Setup() (error) {
    u, err := CreateUser()
    if err != nil {
        panic(err)
    }

    u.FirstName    = "Admin"
    u.LastName     = ""
    u.Email        = "admin@example.com"
    u.PasswordHash = u.GenPasswordHash("secret")

    _, err = u.Save()
    if err != nil {
        panic(err)
    }

    return nil
}

// SetupHandler is the handler for the /setup URI.
func SetupHandler(w http.ResponseWriter, r *http.Request) {
    if SetupAlreadyComplete() {
        http.NotFound(w, r)
        return
    }

    err := Setup()
    if err != nil {
        panic(err)
    }
}