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

func FilterFileOrDirUsingGlobPatterns(rootSearchDir string, dirsGlobList []string,
	includeGlobPattern, excludeGlobPattern string) ([]string, []string, error) {

	var completePathsList []string
	var relativePathList []string

	for _, dirPattern := range dirsGlobList {

		rootSearchDirFS := os.DirFS(rootSearchDir)

		relPattern := strings.TrimPrefix(dirPattern, rootSearchDir+"/")

		matchedDirs, err := doublestar.Glob(rootSearchDirFS, relPattern)
		if err != nil {
			return nil, nil, err
		}

		for _, dir := range matchedDirs {
			completePath := filepath.Join(rootSearchDir, dir)
			completePathsList = append(completePathsList, completePath)
			relativePathList = append(relativePathList, dir)
		}

		//for _, dir := range matchedDirs {
		//	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//		if err != nil {
		//			return err
		//		}
		//		// Apply inclusion pattern if provided
		//		if includeGlobPattern != "" {
		//			matches, err := doublestar.Match(includeGlobPattern, path)
		//			if err != nil || !matches {
		//				return nil
		//			}
		//		}
		//		// Apply exclusion pattern if provided
		//		if excludeGlobPattern != "" {
		//			excluded, err := doublestar.Match(excludeGlobPattern, path)
		//			if err != nil {
		//				return err
		//			}
		//			if excluded {
		//				return nil
		//			}
		//		}
		//		completePathsList = append(completePathsList, path)
		//		return nil
		//	})
		//	if err != nil {
		//		return nil, nil, err
		//	}
		//}
	}

	//// Generate relative paths
	//for _, path := range completePathsList {
	//	relPath, err := filepath.Rel(".", path)
	//	if err != nil {
	//		return nil, nil, err
	//	}
	//	relativePathList = append(relativePathList, relPath)
	//}

	return completePathsList, relativePathList, nil
}

func ToStringArrayFromCsvString(input string) []string {
	return strings.Split(input, ",")
}

func IsMapHasAllStrings(m map[string]interface{}, strList []string) bool {
	for _, str := range strList {
		found := false
		for _, val := range m {
			if valStr, ok := val.(string); ok && valStr == str {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ToJsonStringFromMap[T any](m T) (string, error) {
	outBytes, err := json.Marshal(m)
	if err == nil {
		return string(outBytes), nil
	}
	return "", err
}

func ToStructFromJsonString[T any](jsonStr string) (T, error) {
	var v T
	err := json.Unmarshal([]byte(jsonStr), &v)
	return v, err
}
