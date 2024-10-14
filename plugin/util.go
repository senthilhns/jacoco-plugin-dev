// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/bmatcuk/doublestar/v4"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func writeCard(path, schema string, card interface{}) {
	data, _ := json.Marshal(map[string]interface{}{
		"schema": schema,
		"data":   card,
	})
	switch {
	case path == "/dev/stdout":
		writeCardTo(os.Stdout, data)
	case path == "/dev/stderr":
		writeCardTo(os.Stderr, data)
	case path != "":
		ioutil.WriteFile(path, data, 0644)
	}
}

func writeCardTo(out io.Writer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	io.WriteString(out, "\u001B]1338;")
	io.WriteString(out, encoded)
	io.WriteString(out, "\u001B]0m")
	io.WriteString(out, "\n")
}

func GetNewError(s string) error {
	return errors.New(s)
}

func LogPrintln(p Plugin, args ...interface{}) {
	if p != nil {
		if p.IsQuiet() {
			return
		}
	}

	log.Println(append([]interface{}{"Plugin Info:"}, args...)...)
}

func IsDirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, err
	}
	return info.IsDir(), nil
}

// GetAllEntriesFromGlobPattern searches for files matching Maven-style glob patterns.
func GetAllEntriesFromGlobPattern(rootDir, globPatterns string) ([]string, error) {
	var matches []string
	patterns := strings.Split(globPatterns, ",")

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			match, err := doublestar.Match(pattern, path)
			if err != nil {
				return err
			}
			if match {
				matches = append(matches, path)
				break
			}
		}
		return nil
	})

	return matches, err
}
