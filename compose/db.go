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
    "errors"
    "gopkg.in/mgo.v2"
)

var MongoSession *mgo.Session

// SetupDatabaseSession initializes the session.
func SetupDatabaseSession() error { 
    mongoSession, err := mgo.Dial(config.DatabaseHost)
    if err != nil {
        return err
    }
    mongoSession.SetMode(mgo.Monotonic, true)
    MongoSession = mongoSession
    return nil
}

// GetDatabaseHandle gets the current database session handle.
func GetDatabaseHandle() (*mgo.Database) {
    if MongoSession == nil {
        panic(errors.New("No session available"))
    }
    return MongoSession.DB(config.DatabaseName)
}

// CleanupDatabaseSession closes the current session.
func CleanupDatabaseSession() error {
    MongoSession.Close()
    return nil
}