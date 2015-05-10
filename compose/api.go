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
    "gopkg.in/mgo.v2/bson"
    "net/http"
    "encoding/json"
    "github.com/zenazn/goji/web"
)

func WriteJson(w http.ResponseWriter, obj interface{}) (error) {
    // Encode the config
    encoding, err := json.MarshalIndent(obj, "", "  ")
    if err != nil {
        return err
    }
    w.Write(encoding)
    return err
}



// LoadConfig loads the configuration file.
func DecodeJsonPayload(r *http.Request, obj interface{}) (error) {
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(obj)
    return err
}

// ApiListPosts is a handler to list posts.
func ApiListPosts(c web.C, w http.ResponseWriter, r *http.Request) {
    posts, err := ListPostHeaders(0, 0, true)
    if err != nil {
        panic(err)
    }

    WriteJson(w, posts)
}

// ApiGetPost is a handler to get a post given an id.
func ApiGetPost(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostById(bson.ObjectIdHex(c.URLParams["id"]))
    if err != nil {
        http.NotFound(w, r)
        return
    }
    WriteJson(w, post)
}

// ApiCreatePost is a handler to create a new post.
func ApiCreatePost(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := CreatePost()
    if err != nil {
        panic(err)
    }

    post.Title = "New Post"
    post.Save()

    WriteJson(w, post)
}

// ApiDeletePost is a handler to delete a post.
func ApiDeletePost(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostById(bson.ObjectIdHex(c.URLParams["id"]))
    if err != nil {
        http.NotFound(w, r)
        return
    }
    post.Delete()
}

// ApiUpdatePost is a handler to update an existing post.
func ApiUpdatePost(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostById(bson.ObjectIdHex((c.URLParams["id"])))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    post = &Post{}
    err = DecodeJsonPayload(r, post)
    if err != nil {
        panic(err)
    }
    post, err = post.Save()
    if err != nil {
        panic(err)
    }
}

// ApiGetFileInfo is a handler to get info for a single file, given a file id.
func ApiGetFileInfo(c web.C, w http.ResponseWriter, r *http.Request) {
    id := bson.ObjectIdHex(c.URLParams["id"])
    info, err := GetFileInfoById(id)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    WriteJson(w, info)
}

// ApiGetFileInfoList is a handler to get the file info for one or more files,
// given their ids.
func ApiGetFileInfoList(c web.C, w http.ResponseWriter, r *http.Request) {
    ids := []string{}
    err := DecodeJsonPayload(r, &ids)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    bson_ids := make([]bson.ObjectId, len(ids), len(ids))
    for i, id := range(ids) {
        bson_ids[i] = bson.ObjectIdHex(id)
    }

    info, err := GetMultFileInfoById(bson_ids)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    WriteJson(w, info)
}

// ApiDeleteFile is a handler to delete a file, given the id.
func ApiDeleteFile(c web.C, w http.ResponseWriter, r *http.Request) {
    file, err := GetFileInfoById(bson.ObjectIdHex(c.URLParams["id"]))
    if err != nil {
        http.NotFound(w, r)
        return
    }
    file.DeleteFile()
}

// ApiUpdateSettings is a handler to update the settings.
func ApiUpdateSettings(c web.C, w http.ResponseWriter, r *http.Request) {
    updates := &struct{
        Email    string `json:"email"`
        Password string `json:"password"`
    }{}

    // Get session
    cookie, err := r.Cookie(CookieName)
    session, err := FindSessionByToken(cookie.Value)

    // Get User
    user, err := FindUserById(session.User)

    err = DecodeJsonPayload(r, updates)
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
func ApiGetSettings(c web.C, w http.ResponseWriter, r *http.Request) {
    settings := &struct{
        Email    string `json:"email"`
    }{}

    // Get session
    cookie, err := r.Cookie(CookieName)
    session, err := FindSessionByToken(cookie.Value)

    // Get User
    user, err := FindUserById(session.User)
    if err != nil {
        panic(err)
    }

    settings.Email = user.Email

    WriteJson(w, settings)
}
