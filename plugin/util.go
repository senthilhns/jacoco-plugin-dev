// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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
	IncludedPathsListWithPrefix []PathWithPrefix
	ExcludedPathsListWithPrefix []PathWithPrefix

	CompleteClassPathPrefix string
	RelativeClassPath       string

	IncludeClassesRelativePathsList []string
	IncludeClassesCompletePathsList []string

	ExcludeClassesRelativePathsList []string
	ExcludeClassesCompletePathsList []string
}

type IncludeExcludesMerged struct {
	CompletePathPrefix string
	RelativePathsList  []string
	CompletePathsList  []string

	CompletePathsWithPrefixList []PathWithPrefix
}

func (i *IncludeExcludesMerged) CopyTo(toDstPathPrefix, buildRootPath string) error {

	err := i.CreateUniqueDirs(toDstPathPrefix)
	if err != nil {
		LogPrintln(nil, "Error in CreateUniqueDirs: ", err.Error())
		return err
	}

	for _, pathWithPrefix := range i.CompletePathsWithPrefixList {
		prefix := pathWithPrefix.CompletePathPrefix
		relPath := pathWithPrefix.RelativePath

		srcPath := filepath.Join(prefix, relPath)
		dstPath := filepath.Join(toDstPathPrefix, relPath)

		err := CopyFile(srcPath, dstPath)
		if err != nil {
			fmt.Println("** Error **:  in copying file: ", err.Error())
		}
	}

	return nil
}

func CopyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Flush content to disk to ensure everything is written
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func (i *IncludeExcludesMerged) CreateUniqueDirs(toDstPathPrefix string) error {
	uniqueDirs := i.GetAllUniqueDirs(toDstPathPrefix)
	for _, dir := range uniqueDirs {
		newDir := filepath.Join(toDstPathPrefix, dir)
		err := CreateDir(newDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *IncludeExcludesMerged) GetAllUniqueDirs(toDstPathPrefix string) []string {

	uniqueDirsMap := map[string]bool{}
	_ = uniqueDirsMap

	for _, pathWithPrefix := range i.CompletePathsWithPrefixList {
		curDir := filepath.Dir(pathWithPrefix.RelativePath)
		uniqueDirsMap[curDir] = true
	}

	uniqDirsList := []string{}
	for dir := range uniqueDirsMap {
		uniqDirsList = append(uniqDirsList, dir)
	}

	return uniqDirsList
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

func MergeIncludeExcludeFileCompletePaths(filesInfoStore []FilesInfoStore) []IncludeExcludesMerged {

	validFileList := []string{}
	excludeMap := make(map[string]bool)
	var validFilesListWithPrefix []PathWithPrefix

	for _, fileInfo := range filesInfoStore {
		for _, excludedPathWithPrefix := range fileInfo.ExcludedPathsListWithPrefix {
			excludeFileCompletePath := filepath.Join(excludedPathWithPrefix.CompletePathPrefix,
				excludedPathWithPrefix.RelativePath)
			excludeMap[excludeFileCompletePath] = true
		}
	}

	for _, fileInfo := range filesInfoStore {
		for _, includedPathWithPrefix := range fileInfo.IncludedPathsListWithPrefix {
			includedFileCompletePath := filepath.Join(includedPathWithPrefix.CompletePathPrefix,
				includedPathWithPrefix.RelativePath)
			if _, excluded := excludeMap[includedFileCompletePath]; !excluded {
				validFileList = append(validFileList, includedFileCompletePath)
				validFilesListWithPrefix = append(validFilesListWithPrefix, includedPathWithPrefix)
			}
		}
	}

	var result []IncludeExcludesMerged
	includeExcludesMerged := IncludeExcludesMerged{
		CompletePathPrefix:          "",
		CompletePathsList:           validFileList,
		CompletePathsWithPrefixList: validFilesListWithPrefix,
	}
	result = append(result, includeExcludesMerged)
	return result
}

type PathWithPrefix struct {
	CompletePathPrefix string
	RelativePath       string
}

func WalkDir2(completePath, relativePath, completePathPrefix string,
	includeGlobPatternStrList, excludeGlobPatternStrList []string) (FilesInfoStore, error) {

	relativeIncludePathsList := []string{}
	relativeExcludePathsList := []string{}

	includedCompletePathsList := []string{}
	excludedCompletePathsList := []string{}

	var includedPathsListWithPrefix []PathWithPrefix
	var excludedPathsListWithPrefix []PathWithPrefix

	for _, includeGlobPatternStr := range includeGlobPatternStrList {

		rootSearchDirFS := os.DirFS(completePath)
		relPattern := strings.TrimPrefix(includeGlobPatternStr, completePath+"/")
		matchedFiles, err := doublestar.Glob(rootSearchDirFS, relPattern)

		if err != nil {
			fmt.Println("Error in doublestar.Glob: ", err.Error())
			return FilesInfoStore{}, err
		}

		for _, matchedFile := range matchedFiles {
			matchedFileCompletePath := filepath.Join(completePath, matchedFile)
			includedCompletePathsList = append(includedCompletePathsList, matchedFileCompletePath)
			pathWithPrefix := PathWithPrefix{
				CompletePathPrefix: completePath,
				RelativePath:       matchedFile,
			}
			includedPathsListWithPrefix = append(includedPathsListWithPrefix, pathWithPrefix)
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

		for _, matchedFile := range matchedFiles {
			matchedFileCompletePath := filepath.Join(completePath, matchedFile)
			excludedCompletePathsList = append(excludedCompletePathsList, matchedFileCompletePath)
			pathWithPrefix := PathWithPrefix{
				CompletePathPrefix: completePath,
				RelativePath:       matchedFile,
			}
			excludedPathsListWithPrefix = append(excludedPathsListWithPrefix, pathWithPrefix)
		}

		relativeExcludePathsList = append(relativeExcludePathsList, matchedFiles...)
	}

	classInfoStore := FilesInfoStore{
		CompleteClassPathPrefix:         completePathPrefix,
		RelativeClassPath:               relativePath,
		IncludeClassesRelativePathsList: relativeIncludePathsList,
		IncludeClassesCompletePathsList: includedCompletePathsList,
		ExcludeClassesRelativePathsList: relativeExcludePathsList,
		ExcludeClassesCompletePathsList: excludedCompletePathsList,
		IncludedPathsListWithPrefix:     includedPathsListWithPrefix,
		ExcludedPathsListWithPrefix:     excludedPathsListWithPrefix,
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

func GetRandomDir(absolutePrefixPath, workspacePrefixStr string) (string, error) {
	randomStr, err := generateRandomString(8)
	if err != nil {
		return "", err
	}

	dirName := fmt.Sprintf("%s-%s", workspacePrefixStr, randomStr)
	completePath := filepath.Join(absolutePrefixPath, dirName)

	if _, err := os.Stat(completePath); os.IsNotExist(err) {
		if err := os.MkdirAll(completePath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return completePath, nil
	} else if err != nil {
		return "", fmt.Errorf("failed to check directory: %w", err)
	}

	return GetRandomDir(absolutePrefixPath, workspacePrefixStr)
}

func GetRandomJacocoWorkspaceDir(absolutePrefixPath string) (string, error) {
	return GetRandomDir(absolutePrefixPath, "jacoco-workspace-")
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func CreateDir(absolutePath string) error {
	if absolutePath == "" || absolutePath == "." || absolutePath == ".." {
		return nil
	}

	err := os.MkdirAll(absolutePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", absolutePath, err)
	}
	return nil
}

//
//
