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
    "fmt"
    "gopkg.in/mgo.v2/bson"
)

const (
    PasswordSalt = "goblog"
)

type User struct {
    Id           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
    FirstName    string        `json:"firstName"     bson:"firstName"`
    LastName     string        `json:"lastName"      bson:"lastName"`
    Email        string        `json:"email"         bson:"email"`
    PasswordHash string        `json:"password"      bson:"password"`
}


func CreateUser() (*User, error) {
    // Create User object
    user := &User{}
    user.Id = bson.NewObjectId()

    return user, nil
}

func FindUserById(id bson.ObjectId) (*User, error) {
    db := GetDatabaseHandle()
    c := db.C("users")
    user := &User{}
    err := c.FindId(id).One(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func FindUserByEmail(email string) (*User, error) {
    db := GetDatabaseHandle()
    c := db.C("users")
    user := &User{}
    err := c.Find(bson.M{"email": email}).One(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (s *User) Destroy() (error) {
    // ...
    return nil
}

func (u *User) Save() (*User, error) {
    db := GetDatabaseHandle()
    // Insert User to database
    c := db.C("users")
    _, err := c.UpsertId(u.Id, u)
    return u, err
}

func (u *User) TestPassword(password string) (bool) {
    return u.PasswordHash == u.GenPasswordHash(password)
}

func (u *User) GenPasswordHash(password string) (string) {
    x := fmt.Sprintf("%s:%s:%s", PasswordSalt, u.Id.Hex(), password)
    return fmt.Sprintf("%x", sha256.Sum256([]byte(x)))
}