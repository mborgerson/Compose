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
    "github.com/ant0ine/go-json-rest/rest"
    "gopkg.in/mgo.v2/bson"
    "net/http"
)

// MakeRestrictedApiHandler wraps an existing API handler s.t. the user must
// be logged in before using this API.
func MakeRestrictedApiHandler(handler func(rest.ResponseWriter, *rest.Request)) (func(rest.ResponseWriter, *rest.Request)) {
    return func(w rest.ResponseWriter, r *rest.Request) {
        // Get cookie
        c, err := r.Cookie(CookieName)
        if err == nil && IsSessionTokenValid(c.Value) {
            // Yes. Continue to handler.
            handler(w, r)
        } else {
            // No. Redirect to login page.
            rest.Error(w, "Unauthorized", http.StatusUnauthorized)
        }
    }
}

// ApiListPosts is a handler to list posts.
func ApiListPosts(w rest.ResponseWriter, r *rest.Request) {
    posts, err := ListPosts(0, 0, true)
    if err != nil {
        panic(err)
    }

    w.WriteJson(posts)
}

// ApiGetPost is a handler to get a post given an id.
func ApiGetPost(w rest.ResponseWriter, r *rest.Request) {
    post, err := FindPostById(bson.ObjectIdHex(r.PathParam("id")))
    if err != nil {
        rest.NotFound(w, r)
        return
    }
    w.WriteJson(post)
}

// ApiCreatePost is a handler to create a new post.
func ApiCreatePost(w rest.ResponseWriter, r *rest.Request) {
    post, err := CreatePost()
    if err != nil {
        panic(err)
    }

    post.Title = "New Post"
    post.Save()

    w.WriteJson(post)
}

// ApiDeletePost is a handler to delete a post.
func ApiDeletePost(w rest.ResponseWriter, r *rest.Request) {
    post, err := FindPostById(bson.ObjectIdHex(r.PathParam("id")))
    if err != nil {
        rest.NotFound(w, r)
        return
    }
    post.Delete()
}

// ApiUpdatePost is a handler to update an existing post.
func ApiUpdatePost(w rest.ResponseWriter, r *rest.Request) {
    post, err := FindPostById(bson.ObjectIdHex((r.PathParam("id"))))
    if err != nil {
        rest.NotFound(w, r)
        return
    }

    post = &Post{}
    err = r.DecodeJsonPayload(post)
    if err != nil {
        panic(err)
    }
    post, err = post.Save()
    if err != nil {
        panic(err)
    }
}

// ApiGetFileInfo is a handler to get info for a single file, given a file id.
func ApiGetFileInfo(w rest.ResponseWriter, r *rest.Request) {
    id := bson.ObjectIdHex(r.PathParam("id"))
    info, err := GetFileInfoById(id)
    if err != nil {
        rest.NotFound(w, r)
        return
    }

    w.WriteJson(info)
}

// ApiGetFileInfoList is a handler to get the file info for one or more files,
// given their ids.
func ApiGetFileInfoList(w rest.ResponseWriter, r *rest.Request) {
    ids := []string{}
    err := r.DecodeJsonPayload(&ids)
    if err != nil {
        rest.NotFound(w, r)
        return
    }

    bson_ids := make([]bson.ObjectId, len(ids), len(ids))
    for i, id := range(ids) {
        bson_ids[i] = bson.ObjectIdHex(id)
    }

    info, err := GetMultFileInfoById(bson_ids)
    if err != nil {
        rest.NotFound(w, r)
        return
    }

    w.WriteJson(info)
}

// ApiDeleteFile is a handler to delete a file, given the id.
func ApiDeleteFile(w rest.ResponseWriter, r *rest.Request) {
    file, err := GetFileInfoById(bson.ObjectIdHex(r.PathParam("id")))
    if err != nil {
        rest.NotFound(w, r)
        return
    }
    file.DeleteFile()
}

// ApiUpdateSettings is a handler to update the settings.
func ApiUpdateSettings(w rest.ResponseWriter, r *rest.Request) {
    updates := &struct{
        Email    string `json:"email"`
        Password string `json:"password"`
    }{}

    // Get session
    c, err := r.Cookie(CookieName)
    session, err := FindSessionByToken(c.Value)

    // Get User
    user, err := FindUserById(session.User)

    err = r.DecodeJsonPayload(updates)
    if err != nil {
        panic(err)
    }

    if updates.Email != "" {
        user.Email = updates.Email
        user.Save()
    }

    if updates.Password != "" {
        user.PasswordHash = user.GenPasswordHash(updates.Password)
        user.Save()
    }
}

// ApiGetSettings is a handler to get the current settings.
func ApiGetSettings(w rest.ResponseWriter, r *rest.Request) {
    settings := &struct{
        Email    string `json:"email"`
    }{}

    // Get session
    c, err := r.Cookie(CookieName)
    session, err := FindSessionByToken(c.Value)

    // Get User
    user, err := FindUserById(session.User)
    if err != nil {
        panic(err)
    }

    settings.Email = user.Email

    w.WriteJson(settings)
}

// GetApiHandler returns the API router.
func GetApiHandler() http.Handler {
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    router, err := rest.MakeRouter(
        &rest.Route{"GET",    "/api/posts",    MakeRestrictedApiHandler(ApiListPosts)},
        &rest.Route{"POST",   "/api/posts",    MakeRestrictedApiHandler(ApiCreatePost)},
        &rest.Route{"GET",    "/api/post/:id", MakeRestrictedApiHandler(ApiGetPost)},
        &rest.Route{"PUT",    "/api/post/:id", MakeRestrictedApiHandler(ApiUpdatePost)},
        &rest.Route{"DELETE", "/api/post/:id", MakeRestrictedApiHandler(ApiDeletePost)},
        &rest.Route{"POST",   "/api/file",     MakeRestrictedApiHandler(ApiGetFileInfoList)},
        &rest.Route{"GET",    "/api/file/:id", MakeRestrictedApiHandler(ApiGetFileInfo)},
        &rest.Route{"DELETE", "/api/file/:id", MakeRestrictedApiHandler(ApiDeleteFile)},
        &rest.Route{"GET",    "/api/settings", MakeRestrictedApiHandler(ApiGetSettings)},
        &rest.Route{"POST",   "/api/settings", MakeRestrictedApiHandler(ApiUpdateSettings)},
    )
    if err != nil {
        panic(err)
    }
    api.SetApp(router)
    return api.MakeHandler()
}
