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
    "github.com/zenazn/goji/web"
    "net/http"
    "strconv"
)

// IndexHandler is the handler for Index page(s).
func IndexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
    page := 1
    if c.URLParams["page"] != "" {
        desiredPage, err := strconv.Atoi(c.URLParams["page"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        page = desiredPage
    }

    // Get total number of posts that are not drafts
    total, err := CountPosts(false)
    if err != nil {
        panic(err)
    }
    numPages := total/config.IndexPostsPerPage
    if total % config.IndexPostsPerPage != 0 {
        numPages += 1
    }

    // Is this page valid?
    if page < 1 || (page != 1 && page > numPages) {
        http.NotFound(w, r)
        return
    }

    v := map[string]interface{}{}
    posts, err := ListPosts((page-1) * config.IndexPostsPerPage,
                            config.IndexPostsPerPage,
                            false)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    v["Posts"] = posts
    v["CurrentPage"] = page
    v["TotalPages"] = numPages

    err = SiteTemplates.ExecuteTemplate(w, "index.html", v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}