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

// AdminHandler is the main handler for all other admin URLs. Because the admin
// pages use Angular routing, just return the index page and let the JS side
// determine what content to show.
func AdminHandler(w http.ResponseWriter, r *http.Request) {
    err := AdminTemplates.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// AdminEditHandler is a handler for the post edit partial.
func AdminEditHandler(w http.ResponseWriter, r *http.Request) {
    err := AdminTemplates.ExecuteTemplate(w, "edit.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// AdminSettingsHandler is a handler for the settings partial.
func AdminSettingsHandler(w http.ResponseWriter, r *http.Request) {
    err := AdminTemplates.ExecuteTemplate(w, "settings.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// AdminPostsHandler is a handler for the posts partial.
func AdminPostsHandler(w http.ResponseWriter, r *http.Request) {
    err := AdminTemplates.ExecuteTemplate(w, "posts.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// AdminAssetsHandler is a handler to serve admin assets (js, css, etc)
func AdminAssetsHandler(w http.ResponseWriter, r *http.Request) {
    if config.AdminAssetsPath != "" {
        fs := http.StripPrefix("/admin/assets/", http.FileServer(http.Dir(config.AdminAssetsPath)))
        fs.ServeHTTP(w, r)
        return
    }
    http.NotFound(w, r)
}