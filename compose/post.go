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
    "github.com/mborgerson/GoTruncateHtml/truncatehtml"
    "github.com/russross/blackfriday"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "html/template"
    "net/http"
    "time"
    "github.com/zenazn/goji/web"
)

type PostHeader struct {
    Id           bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
    Title        string          `json:"title"         bson:"title"`
    Date         time.Time       `json:"date"          bson:"date"`
    LastModified time.Time       `json:"last_modified" bson:"last_modified"`
    Slug         string          `json:"slug"          bson:"slug"`
    Draft        bool            `json:"draft"         bson:"draft"`
    Files        []bson.ObjectId `json:"files"         bson:"files"`
}

type Post struct {
                 PostHeader      `json:",inline"       bson:",inline"`
    Body         string          `json:"body"          bson:"body"`
}

// FindPostBySlug finds a post by the slug. An error is returned if the post
// for the given slug could not be found.
func FindPostBySlug(slug string) (*Post, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    post := &Post{}
    err := c.Find(bson.M{"slug":slug}).One(post)
    if err != nil {
        return nil, err
    }
    return post, nil
}

// FindPostById finds a post given a post id. An error is ruterned if the post
// for the given id could not be found.
func FindPostById(id bson.ObjectId) (*Post, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    post := &Post{}
    err := c.FindId(id).One(post)
    if err != nil {
        return nil, err
    }
    return post, nil
}

// ListPosts will return a slice of limit reverse-chronologicaly orderded posts,
// starting from start and optionally including drafts.
func ListPosts(start int, limit int, includeDrafts bool) ([]Post, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    var posts []Post
    posts = nil
    var q *mgo.Query
    if includeDrafts {
        q = c.Find(nil)
    } else {
        q = c.Find(bson.M{"draft":false})
    }
    err := q.Sort("-date").Skip(start).Limit(limit).All(&posts)
    return posts, err
}

// ListPostHeaders will return a slice of limit reverse-chronologicaly orderded posts,
// starting from start and optionally including drafts.
func ListPostHeaders(start int, limit int, includeDrafts bool) ([]PostHeader, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    var posts []PostHeader
    posts = nil
    var q *mgo.Query
    if includeDrafts {
        q = c.Find(nil)
    } else {
        q = c.Find(bson.M{"draft":false})
    }
    err := q.Sort("-date").Skip(start).Limit(limit).All(&posts)
    return posts, err
}

// CreatePost creates a new post object. Call Save() on the post to write it
// to the database.
func CreatePost() (*Post, error) {
    id := bson.NewObjectId()
    newPost := &Post{
        PostHeader: PostHeader{
            Id:    id,
            Title: "",
            Date:  time.Now(),
            Slug:  "",
            Draft: true,
            Files: []bson.ObjectId{}},
        Body:  ""}
    return newPost, nil
}

// CountPosts counts the total number of posts, optionally including drafts.
func CountPosts(includeDrafts bool) (int, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    var q *mgo.Query
    if includeDrafts {
        q = c.Find(nil)
    } else {
        q = c.Find(bson.M{"draft":false})
    }
    return q.Count() 
}

// Save writes the post to the database.
func (post *Post) Save() (*Post, error) {
    db := GetDatabaseHandle()
    c := db.C("posts")
    post.LastModified = time.Now()
    _, err := c.UpsertId(post.Id, post)
    return post, err
}

// Delete removes the post and all of its files from the database.
func (post *Post) Delete() (*Post, error) {
    db := GetDatabaseHandle()
    // Delete post files
    file_infos, err := GetMultFileInfoById(post.Files)
    if err != nil {
        panic(err)
    }
    for _, info := range file_infos {
        info.DeleteFile()
    }

    c := db.C("posts")
    err = c.RemoveId(post.Id)
    return post, err
}

// RenderBody renders the Markdown formatted body of the post into HTML.
func (post *Post) RenderBody() (template.HTML, error) {
    return template.HTML(blackfriday.MarkdownCommon([]byte(post.Body))), nil
}

// RenderBodySnippet renders a truncated HTML version of the Markdown formatted
// body.
func (post *Post) RenderBodySnippet(maxlen int, ellipsis string) (template.HTML, error) {
    body := blackfriday.MarkdownCommon([]byte(post.Body))

    truncated, err := truncatehtml.TruncateHtml(body, maxlen, ellipsis)
    if err != nil {
        return template.HTML(""), err
    }

    return template.HTML(truncated), nil
}

// ViewHandler is the handler for viewing a post.
func ViewHandler(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostBySlug(c.URLParams["slug"])
    if post == nil {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        panic(err)
    }

    if CheckModifiedHandler(w, r, post.LastModified) {
        // Not modified
        return
    }

    w.Header().Set("Last-Modified", post.LastModified.UTC().Format(HttpDateTimeFormat))
    err = SiteTemplates.ExecuteTemplate(w, "post.html", post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
// ViewHandlerRemoveTrailingSlash will remove the trailing slash from a valid
// post url.
func ViewHandlerRemoveTrailingSlash(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostBySlug(c.URLParams["slug"])
    if post == nil {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        panic(err)
    }

    http.Redirect(w, r, "/"+c.URLParams["slug"], http.StatusMovedPermanently)
}

// ViewFileHandler is the handler for viewing a post file.
func ViewFileHandler(c web.C, w http.ResponseWriter, r *http.Request) {
    post, err := FindPostBySlug(c.URLParams["slug"])
    if post == nil {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        panic(err)
    }

    // Get the list of files for this post and see if the filename in this
    // request matches any of those files
    file_infos, err := GetMultFileInfoById(post.Files)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    for _, info := range file_infos {
        if info.Name == c.URLParams["file"] {
            DownloadHandler(w, r, info.Id)
            return
        }
    }
    http.NotFound(w, r)
}