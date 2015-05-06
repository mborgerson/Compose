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
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "io"
    "net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    db := GetDatabaseHandle()
    
    // Create JSON encoder for the response
    enc := json.NewEncoder(w)
    response := make(map[string]string)

    // Get handle to the file stream
    file, header, err := r.FormFile("file")
    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        enc.Encode(response)
        return
    }
    defer file.Close()

    // Create file in database
    out, err := db.GridFS("fs").Create(header.Filename)
    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        enc.Encode(response)
        return
    }
    defer out.Close()

    // Insert data into file
    _, err = io.Copy(out, file)
    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        enc.Encode(response)
        return
    }

    // Return response object
    response["status"] = "success"
    response["message"] = "file uploaded successfully"
    response["_id"] = bson.ObjectId(out.Id().(bson.ObjectId)).Hex()
    enc.Encode(response)
}

type FileInfo struct {
    Id         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
    Name       string        `json:"filename"      bson:"filename"` 
    UploadDate string        `json:"uploadDate"    bson:"uploadDate"`
    Size       int64         `json:"size"          bson:"size"` 
}

func GetFileInfoById(id bson.ObjectId) (*FileInfo, error) {
    db := GetDatabaseHandle()
    c := db.GridFS("fs")
    file, err := c.OpenId(id)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    info := &FileInfo{Id:         file.Id().(bson.ObjectId),
                      Name:       file.Name(),
                      UploadDate: file.UploadDate().String(),
                      Size:       file.Size()}

    return info, nil
}

func GetFileById(id bson.ObjectId) (*mgo.GridFile, error) {
    db := GetDatabaseHandle()
    c := db.GridFS("fs")
    file, err := c.OpenId(id)
    return file, err
}

func GetMultFileInfoById(ids []bson.ObjectId) (map[bson.ObjectId]*FileInfo, error) {
    out := make(map[bson.ObjectId]*FileInfo)
    for _, id := range ids {
        info, err := GetFileInfoById(id)
        if err != nil {
            out[id] = nil
            continue
        }
        out[id] = info
    }
    return out, nil
}

func (file *FileInfo) DeleteFile() (*FileInfo, error) {
    db := GetDatabaseHandle()
    c := db.GridFS("fs")
    err := c.RemoveId(file.Id)
    return file, err
}