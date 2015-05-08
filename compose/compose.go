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
    "fmt"
    "html/template"
    "io"
    "mime"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

var templates *template.Template
var AdminTemplates *template.Template
var IndexRegExp = regexp.MustCompile("^/([0-9]+)?$")
var PostRegExp = regexp.MustCompile("^/([A-Za-z0-9-]+)(/.*)?$")

// MainHandler is the handler which will determine if the request is for the
// index or for a post. It will then call the appropriate handler with required
// parameters. 
func MainHandler(w http.ResponseWriter, r *http.Request) {
    // Get slug
    m := IndexRegExp.FindStringSubmatch(r.URL.Path)
    if m != nil {
        // Index page
        page := 1
        if m[1] != "" {
            pageConv, err := strconv.ParseInt(m[1], 10, 0)
            if err == nil {
                page = int(pageConv)
            }
        }

        IndexHandler(w, r, page)
        return
    }

    // Get slug
    m = PostRegExp.FindStringSubmatch(r.URL.Path)
    if m != nil {
        if m[2] == "" {
            // Enforce trailing slash
            http.Redirect(w, r, strings.Join([]string{r.URL.Path, "/"}, ""), 301)
            return
        }

        ViewHandler(w, r, m[1], m[2])
        return
    }
}

// ViewHandler is the handler for viewing a post.
func ViewHandler(w http.ResponseWriter, r *http.Request, slug string, filename string) {
    post, err := FindPostBySlug(slug)
    if post == nil {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        panic(err)
    }

    if filename == "/" {
        err = templates.ExecuteTemplate(w, "post.html", post)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    } else {
        // Get the list of files for this post and see if this request
        // matches any of those files
        file_infos, err := GetMultFileInfoById(post.Files)
        if err != nil {
            http.NotFound(w, r)
            return
        }
        for _, info := range file_infos {
            if info.Name == filename[1:len(filename)] {
                // Guess mime from extension
                ext := filepath.Ext(info.Name)
                mime_type := mime.TypeByExtension(ext)

                // Set header
                if mime_type != "" {
                    w.Header().Set("Content-Type", mime_type)
                }

                // Send File
                file, _ := GetFileById(info.Id)
                defer file.Close()
                io.Copy(w, file)
                return
            }
        }
        http.NotFound(w, r)
        return
    }
}

// IndexHandler is the handler for Index pages.
func IndexHandler(w http.ResponseWriter, r *http.Request, page int) {
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

    err = templates.ExecuteTemplate(w, "index.html", v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// AssetsHandler is the handler that will serve assets.
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
    if config.AssetsPath != "" {
        fs := http.StripPrefix("/assets/", http.FileServer(http.Dir(config.AssetsPath)))
        fs.ServeHTTP(w, r)
        return
    }
    http.NotFound(w, r)
}

// BuildTemplates builds all the required site templates.
func BuildTemplates() error {
    funcMap := template.FuncMap {
        "add": func(a, b int) int { return a+b },
        "sub": func(a, b int) int { return a-b },
    }

    files := []string{
        "index.html",
        "post.html",
    }

    for i, file := range files {
        files[i] = filepath.Join(config.TemplatesPath, file)
    }

    _templates := template.New("base")
    _templates.Delims("<%", "%>")
    _templates.Funcs(funcMap)
    _templates, err := _templates.ParseFiles(files...)
    if err != nil {
        return err
    }
    templates = _templates
    return nil
}

// BuildAdminTemplates builds all the required admin templates.
func BuildAdminTemplates() error {
    funcMap := template.FuncMap {
        "add": func(a, b int) int { return a+b },
        "sub": func(a, b int) int { return a-b },
    }

    files := []string{
        "index.html", 
        "edit.html",
        "posts.html",
        "settings.html",
        "login.html",
    }

    for i, file := range files {
        files[i] = filepath.Join(config.AdminTemplatesPath, file)
    }

    _templates := template.New("base")
    _templates.Delims("<%", "%>")
    _templates.Funcs(funcMap)
    _templates, err := _templates.ParseFiles(files...)
    if err != nil {
        return err
    }
    AdminTemplates = _templates
    return nil
}

// main is the entry point. Loads the program resources and begins waiting for
// connections.
func main() {
    // Create a config file with the defaults
    if !FileExists(ConfigDefaultFilename) {
        config, _ := GetDefaultConfig()
        err := config.Save(ConfigDefaultFilename)
        if err != nil {
            fmt.Println("I tried to create default config file but failed. Check directory permissions.")
            os.Exit(1)
        }
        fmt.Println("The config file could not be found, so I created a config file at '", ConfigDefaultFilename, "'. Please ensure this file contains the correct values and relaunch.")
        os.Exit(0)
    }

    // Load config
    _, err := GetConfig()
    if err != nil {
        fmt.Println("Failed to load the config file:", err.Error())
        os.Exit(1)
    }

    // Build Templates
    err = BuildTemplates()
    if err != nil {
        fmt.Println("Failed to build templates:", err.Error())
        os.Exit(1)
    }
    err = BuildAdminTemplates()
    if err != nil {
        fmt.Println("Failed to build admin templates:", err.Error())
        os.Exit(1)
    }

    // Connect to the database
    err = SetupDatabaseSession()
    if err != nil {
        panic(err)
    }
    defer CleanupDatabaseSession()

    // Setup the router
    http.Handle("/api/",                         GetApiHandler())
    http.HandleFunc("/admin/assets/",            AdminAssetsHandler)
    http.HandleFunc("/admin/partials/edit/",     MakeRestrictedHttpHandler(AdminEditHandler))
    http.HandleFunc("/admin/partials/posts/",    MakeRestrictedHttpHandler(AdminPostsHandler))
    http.HandleFunc("/admin/partials/settings/", MakeRestrictedHttpHandler(AdminSettingsHandler))
    http.HandleFunc("/admin/",                   MakeRestrictedHttpHandler(AdminHandler))
    http.HandleFunc("/assets/",                  AssetsHandler)
    http.HandleFunc("/login/",                   LoginHandler)
    http.HandleFunc("/logout/",                  LogoutHandler)
    http.HandleFunc("/setup/",                   SetupHandler)
    http.HandleFunc("/upload/",                  MakeRestrictedHttpHandler(UploadHandler))
    http.HandleFunc("/",                         MainHandler)

    // Begin serving
    http.ListenAndServe(fmt.Sprintf("%s:%d", config.BindHost, config.BindPort), nil)
}