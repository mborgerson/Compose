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
    "crypto/sha256"
    "errors"
    "fmt"
    "gopkg.in/mgo.v2/bson"
    "time"
)

type Session struct {
    Id    bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
    Token string        `json:"token"         bson:"token"`
    User  bson.ObjectId `json:"userId"        bson:"userId"`
}

// CreateSession creates a new session object. Save() should be called on the
// session once it should be written to the database.
func CreateSession(user *User) (*Session, error) {
    // Valid user?
    if user == nil {
        return nil, errors.New("Invalid user")
    }

    // Create session object
    session := &Session{}
    session.Id = bson.NewObjectId()
    session.User = user.Id

    // Create the session token ... sha256(email:timestamp)
    x := fmt.Sprintf("%s:%s", user.Email, time.Now().String())
    session.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(x)))

    return session, nil
}

// FindSessionById looks up a session by the session id.
func FindSessionById(id bson.ObjectId) (*Session, error) {
    c := GetDatabaseHandle().C("sessions")

    // Find session by id
    session := &Session{}
    err := c.FindId(id).One(session)
    if err != nil {
        return nil, err
    }

    return session, nil
}

// FindSessionByToken looks up a session by the session token.
func FindSessionByToken(token string) (*Session, error) {
    c := GetDatabaseHandle().C("sessions")

    // Find session by token
    session := &Session{}
    err := c.Find(bson.M{"token": token}).One(session)
    if err != nil {
        return nil, err
    }

    return session, nil
}

// IsSessionTokenValid determines if a session token is valid.
func IsSessionTokenValid(token string) (bool) {
    _, err := FindSessionByToken(token)
    return err == nil
}

// Destroy destroys a session.
func (s *Session) Destroy() (error) {
    // ...
    return nil
}

// Save updates or creates a session in the database.
func (s *Session) Save() (*Session, error) {
    c := GetDatabaseHandle().C("sessions")
    _, err := c.UpsertId(s.Id, s)
    return s, err
}
