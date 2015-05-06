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
    "os"
)

const (
    ConfigDefaultFilename = "compose.json"
)

type Config struct {
    DatabaseHost       string
    DatabaseName       string
    AssetsPath         string
    TemplatesPath      string
    AdminAssetsPath    string
    AdminTemplatesPath string
    IndexPostsPerPage  int
}

var config *Config = nil

// GetDefaultConfig returns the default configuration settings.
func GetDefaultConfig() (*Config, error) {
    return &Config{
        DatabaseHost:       "127.0.0.1",
        DatabaseName:       "compose",
        AssetsPath:         "theme_site/dist/assets",
        TemplatesPath:      "theme_site/dist/templates",
        AdminAssetsPath:    "theme_admin/dist/assets",
        AdminTemplatesPath: "theme_admin/dist/templates",
        IndexPostsPerPage:  5,
    }, nil
}

// FileExists returns a bool indicating whether the path exists or not.
func FileExists(path string) (bool) {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// LoadConfig loads the configuration file.
func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Decode the config
    decoder := json.NewDecoder(file)
    config := &Config{}
    err = decoder.Decode(config)
    if err != nil {
        return nil, err
    }

    return config, nil
}

// GetConfig gets the global configuration. It will load the config upon first
// call. 
func GetConfig() (*Config, error) {
    if config == nil {
        lconfig, err := LoadConfig(ConfigDefaultFilename)
        if err != nil {
            // Config could not be loaded
            return nil, err
        }

        // Assign to global config ptr
        config = lconfig
    }
    return config, nil
}

// Save saves the current configuration.
func (c *Config) Save(filename string) (error) {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Encode the config
    encoding, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    file.Write(encoding)
    return err
}