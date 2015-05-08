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
    "time"
)

const HttpDateTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// CheckModifiedHandler will check the request for an If-Modified-Since header.
// If present, it will compare If-Modified-Since to the lastModified parameter.
// If lastModified is equal or greater than If-Modified-Since, the "304 Not
// Modified" response is sent and true is returned. Otherwise, false is
// returned.
func CheckModifiedHandler(w http.ResponseWriter, r *http.Request, lastModified time.Time) (bool) {
    // Check for If-Modified-Since in header
    modifiedSince, err := http.ParseTime(r.Header.Get("If-Modified-Since"))
    if err != nil {
        // Error parsing time
        return false
    }

    // Make sure lastModified is in UTC and with Second resolution
    lastModified = lastModified.UTC().Truncate(time.Second)

    // Check to see if the content was modified
    if lastModified.After(modifiedSince) {
        // Content has been modified
        return false
    } else {
        // Content has not been modified
        w.WriteHeader(http.StatusNotModified)
        return true
    }
}
