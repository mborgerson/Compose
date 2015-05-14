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
    "github.com/zenazn/goji"
    "html/template"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
)

var SiteTemplates *template.Template
var AdminTemplates *template.Template

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

    tmpl := template.New("base")
    tmpl.Delims("<%", "%>")
    tmpl.Funcs(funcMap)
    tmpl, err := tmpl.ParseFiles(files...)
    if err != nil {
        return err
    }
    SiteTemplates = tmpl

    files = []string{
        "index.html", 
        "edit.html",
        "posts.html",
        "settings.html",
        "login.html",
    }

    for i, file := range files {
        files[i] = filepath.Join(config.AdminTemplatesPath, file)
    }

    tmpl = template.New("base")
    tmpl.Delims("<%", "%>")
    tmpl.Funcs(funcMap)
    tmpl, err = tmpl.ParseFiles(files...)
    if err != nil {
        return err
    }
    AdminTemplates = tmpl
    return nil
}

func MakeStaticHandler(prefix, dir string) (http.HandlerFunc) {
    return http.StripPrefix(prefix, http.FileServer(http.Dir(dir))).ServeHTTP
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

    // Connect to the database
    err = SetupDatabaseSession()
    if err != nil {
        panic(err)
    }
    defer CleanupDatabaseSession()

    // Setup the router
    goji.Get(    "/setup",                   SetupHandler)
    goji.Get(    "/admin/assets/*",          MakeStaticHandler("/admin/assets/", config.AdminAssetsPath))
    goji.Get(    "/admin/partials/edit",     MakeRestrictedHttpHandler(AdminEditHandler))
    goji.Get(    "/admin/partials/posts",    MakeRestrictedHttpHandler(AdminPostsHandler))
    goji.Get(    "/admin/partials/settings", MakeRestrictedHttpHandler(AdminSettingsHandler))
    goji.Get(    "/admin/",                  http.RedirectHandler("/admin", http.StatusMovedPermanently))
    goji.Get(    "/admin",                   MakeRestrictedHttpHandler(AdminHandler))
    goji.Get(    "/admin/*",                 MakeRestrictedHttpHandler(AdminHandler))
    goji.Post(   "/upload",                  MakeRestrictedHttpHandler(UploadHandler))
    goji.Get(    "/api/posts",               MakeRestrictedHttpHandler(ApiListPosts))
    goji.Post(   "/api/posts",               MakeRestrictedHttpHandler(ApiCreatePost))
    goji.Get(    "/api/post/:id",            MakeRestrictedHttpHandler(ApiGetPost))
    goji.Put(    "/api/post/:id",            MakeRestrictedHttpHandler(ApiUpdatePost))
    goji.Delete( "/api/post/:id",            MakeRestrictedHttpHandler(ApiDeletePost))
    goji.Post(   "/api/file",                MakeRestrictedHttpHandler(ApiGetFileInfoList))
    goji.Get(    "/api/file/:id",            MakeRestrictedHttpHandler(ApiGetFileInfo))
    goji.Delete( "/api/file/:id",            MakeRestrictedHttpHandler(ApiDeleteFile))
    goji.Get(    "/api/settings",            MakeRestrictedHttpHandler(ApiGetSettings))
    goji.Post(   "/api/settings",            MakeRestrictedHttpHandler(ApiUpdateSettings))
    goji.Get(    "/assets/*",                MakeStaticHandler("/assets/", config.AssetsPath))
    goji.Get(    "/login",                   LoginHandler)
    goji.Post(   "/login",                   LoginHandler)
    goji.Get(    "/logout",                  LogoutHandler)
    indexRegexp := regexp.MustCompile("^/(?P<page>[0-9]*)$")
    goji.Get(    indexRegexp,                IndexHandler)
    goji.Get(    "/:slug",                   ViewHandler)
    goji.Get(    "/:slug/",                  ViewHandlerRemoveTrailingSlash)
    goji.Get(    "/:slug/:file",             ViewFileHandler)

    // Begin serving
    goji.Serve()
}