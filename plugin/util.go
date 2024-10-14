// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	includeGlobPatternCsvStr, excludeGlobPatternCsvStr string, autoFillIncludePattern string) ([]FilesInfoStore, error) {

	var classesInfoStoreList []FilesInfoStore

	if len(includeGlobPatternCsvStr) == 0 {
		if len(autoFillIncludePattern) > 0 {
			includeGlobPatternCsvStr = autoFillIncludePattern
		}
	}

	includeGlobPatternStrList := ToStringArrayFromCsvString(includeGlobPatternCsvStr)
	excludeGlobPatternStrList := ToStringArrayFromCsvString(excludeGlobPatternCsvStr)

	for _, dirPattern := range dirsGlobList {

		rootSearchDirFS := os.DirFS(rootSearchDir)

		relPattern := strings.TrimPrefix(dirPattern, rootSearchDir+"/")

		matchedDirs, err := doublestar.Glob(rootSearchDirFS, relPattern)
		if err != nil {
			return classesInfoStoreList, err
		}

		for _, relativePath := range matchedDirs {
			completePath := filepath.Join(rootSearchDir, relativePath)
			classesInfoStore, err := WalkDir2(completePath, relativePath, rootSearchDir+"/",
				includeGlobPatternStrList, excludeGlobPatternStrList)
			if err != nil {
				fmt.Println("Error in WalkDir: ", err.Error())
				return classesInfoStoreList, err
			}
			classesInfoStoreList = append(classesInfoStoreList, classesInfoStore)
		}

	}

	return classesInfoStoreList, nil
}

type FilesInfoStore struct {
	CompleteClassPathPrefix         string
	RelativeClassPath               string
	IncludeClassesRelativePathsList []string
	ExcludeClassesRelativePathsList []string
}

type IncludeExcludesMerged struct {
	CompletePathPrefix string
	RelativePathsList  []string
}

func MergeIncludeExcludeFilePaths(filesInfoStore []FilesInfoStore) []IncludeExcludesMerged {
	mergedResults := make(map[string]map[string]struct{})

	for _, fileInfo := range filesInfoStore {
		if _, exists := mergedResults[fileInfo.CompleteClassPathPrefix]; !exists {
			mergedResults[fileInfo.CompleteClassPathPrefix] = make(map[string]struct{})
		}
		excludeSet := make(map[string]struct{})
		for _, excludePath := range fileInfo.ExcludeClassesRelativePathsList {
			excludeSet[excludePath] = struct{}{}
		}

		for _, includePath := range fileInfo.IncludeClassesRelativePathsList {
			if _, excluded := excludeSet[includePath]; !excluded {
				mergedResults[fileInfo.CompleteClassPathPrefix][includePath] = struct{}{}
			}
		}
	}

	var result []IncludeExcludesMerged
	for prefix, paths := range mergedResults {
		var relativePaths []string
		for path := range paths {
			relativePaths = append(relativePaths, path)
		}
		result = append(result, IncludeExcludesMerged{
			CompletePathPrefix: prefix,
			RelativePathsList:  relativePaths,
		})
	}

	return result
}

func WalkDir2(completePath, relativePath, completePathPrefix string,
	includeGlobPatternStrList, excludeGlobPatternStrList []string) (FilesInfoStore, error) {

	relativeIncludePathsList := []string{}
	relativeExcludePathsList := []string{}

	for _, includeGlobPatternStr := range includeGlobPatternStrList {

		rootSearchDirFS := os.DirFS(completePath)
		relPattern := strings.TrimPrefix(includeGlobPatternStr, completePath+"/")
		matchedFiles, err := doublestar.Glob(rootSearchDirFS, relPattern)

		if err != nil {
			fmt.Println("Error in doublestar.Glob: ", err.Error())
		}

		relativeIncludePathsList = append(relativeIncludePathsList, matchedFiles...)
	}

	for _, excludeGlobPatternStr := range excludeGlobPatternStrList {

		rootSearchDirFS := os.DirFS(completePath)
		relPattern := strings.TrimPrefix(excludeGlobPatternStr, completePath+"/")
		matchedFiles, err := doublestar.Glob(rootSearchDirFS, relPattern)

		if err != nil {
			fmt.Println("Error in doublestar.Glob: ", err.Error())
		}

		relativeExcludePathsList = append(relativeExcludePathsList, matchedFiles...)
	}

	classInfoStore := FilesInfoStore{
		CompleteClassPathPrefix:         completePathPrefix,
		RelativeClassPath:               relativePath,
		IncludeClassesRelativePathsList: relativeIncludePathsList,
		ExcludeClassesRelativePathsList: relativeExcludePathsList,
	}

	return classInfoStore, nil
}

func TrimStrings(input []string) []string {
	var trimmed []string
	for _, str := range input {
		trimmed = append(trimmed, strings.TrimSpace(str))
	}
	return trimmed
}

func ToStringArrayFromCsvString(input string) []string {
	sl := strings.Split(input, ",")

	return TrimStrings(sl)
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
