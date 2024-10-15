package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func JacocoFilesExistTest() error {

	paths := []string{
		"/opt/hns/test-resources/game-of-life-master/gameoflife-core/target/jacoco.exec",
		"/opt/hns/test-resources/game-of-life-master/gameoflife-web/target/jacoco.exec",
	}

	var err error
	for _, path := range paths {
		_, err = os.Stat(path)

		if err != nil {
			LogPrintln(nil, "Error in JacocoFilesExistTest: "+err.Error())
			return err
		}

	}

	return nil

}

func TestMain(m *testing.M) {

	err := JacocoFilesExistTest()
	if err != nil {
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestIsBuildRootExists(t *testing.T) {
	args := GetTestNewArgs()

	_, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestIsBuildRootExists: %s", err.Error())
	}
}

func TestExecPathPatterns(t *testing.T) {
	CheckExecPathPattern(TestExecPathPattern01, t)
	CheckExecPathPattern(TestExecPathPattern02, t)
	CheckExecPathPattern(TestExecPathPattern03, t)
	CheckExecPathPattern(TestExecPathPattern04, t)
}

func TestEmptyExecPathPattern(t *testing.T) {
	args := GetTestNewArgs()
	args.ExecFilesPathPattern = ""
	_, err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestEmptyExecPathPattern is accepted")
	}
}

func CheckExecPathPattern(globPattern string, t *testing.T) {
	args := GetTestNewArgs()
	args.ExecFilesPathPattern = globPattern
	_, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("CheckExecPathPattern for globPattern: %s" + globPattern + " err == " + err.Error())
	}
}

type ClassesTestInfo struct {
	ClassesInfoStoreList []struct {
		CompleteClassPathPrefix         string   `json:"CompleteClassPathPrefix"`
		RelativeClassPath               string   `json:"RelativeClassPath"`
		IncludeClassesRelativePathsList []string `json:"IncludeClassesRelativePathsList"`
		ExcludeClassesRelativePathsList []string `json:"ExcludeClassesRelativePathsList"`
	} `json:"ClassesInfoStoreList"`
}

func (c *ClassesTestInfo) isAllRequiredIncludePathsPresent(requiredIncludePaths []string) bool {
	pathSet := make(map[string]struct{})

	for _, classInfo := range c.ClassesInfoStoreList {
		for _, classPath := range classInfo.IncludeClassesRelativePathsList {
			pathSet[classPath] = struct{}{}
		}
	}

	for _, required := range requiredIncludePaths {
		if _, found := pathSet[required]; !found {
			fmt.Println("Required path not found: ", required)
			return false
		}
	}
	return true
}

func (c *ClassesTestInfo) isAllRequiredExcludePathsPresent(requiredExcludePaths []string) bool {
	excludePathSet := make(map[string]struct{})

	for _, classInfo := range c.ClassesInfoStoreList {
		for _, excludePath := range classInfo.ExcludeClassesRelativePathsList {
			excludePathSet[excludePath] = struct{}{}
		}
	}

	for _, required := range requiredExcludePaths {
		if _, found := excludePathSet[required]; !found {
			return false
		}
	}
	return true
}

func TestClassPathWithNoIncludeNoExclude(t *testing.T) {

	classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
		"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	classInclusionPatterns := ""
	classExclusionPatterns := ""

	expectedIncludePaths := []string{
		"com/wakaleo/gameoflife/domain/Cell.class",
		"com/wakaleo/gameoflife/domain/Grid.class",
		"com/wakaleo/gameoflife/domain/GridReader.class",
		"com/wakaleo/gameoflife/domain/GridWriter.class",
		"com/wakaleo/gameoflife/domain/Universe.class",
		"com/wakaleo/gameoflife/webtests/controllers/GameController.class",
		"com/wakaleo/gameoflife/webtests/controllers/HomePageController.class",
	}

	expectedExcludePaths := []string{}

	CheckClassPathWithIncludeExcludeVariation(classPatterns, classInclusionPatterns,
		classExclusionPatterns, expectedIncludePaths, expectedExcludePaths, t)

}

func TestSourcePathWithIncludeAndExclude(t *testing.T) {

	classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
		"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	classInclusionPatterns := "**/*.class, **/*.xml"
	classExclusionPatterns := "**/controllers/*.class"

	sourcePatterns := "**/src/main/java"
	sourceInclusionPatterns := "**/*.java, *.groovy"
	sourceExclusionPatterns := "**/controllers/*.java"

	expectedIncludePaths := []string{
		"com/wakaleo/gameoflife/domain/Cell.java",
		"com/wakaleo/gameoflife/domain/Grid.java",
		"com/wakaleo/gameoflife/domain/GridReader.java",
		"com/wakaleo/gameoflife/domain/GridWriter.java",
		"com/wakaleo/gameoflife/domain/Universe.java",
	}
	expectedExcludePaths := []string{
		"com/wakaleo/gameoflife/webtests/controllers/GameController.java",
		"com/wakaleo/gameoflife/webtests/controllers/HomePageController.java",
	}

	CheckSourceAndClassPathsWithIncludeExcludeVariations(sourcePatterns, sourceInclusionPatterns, sourceExclusionPatterns,
		classPatterns, classInclusionPatterns, classExclusionPatterns, expectedIncludePaths, expectedExcludePaths, t)

}

type SourceTestInfo struct {
	FinalizedSourcesList []struct {
		CompleteClassPathPrefix         string   `json:"CompleteClassPathPrefix"`
		RelativeClassPath               string   `json:"RelativeClassPath"`
		IncludeClassesRelativePathsList []string `json:"IncludeClassesRelativePathsList"`
		ExcludeClassesRelativePathsList []string `json:"ExcludeClassesRelativePathsList"`
	} `json:"FinalizedSourcesList"`
}

type WorkSpaceInfo struct {
	WorkSpaceCompletePathKeyStr struct {
		Classes   string `json:"classes"`
		ExecFiles string `json:"execFiles"`
		Sources   string `json:"sources"`
		Workspace string `json:"workspace"`
	} `json:"WorkSpaceCompletePathKeyStr"`
}

func CheckSourceAndClassPathsWithIncludeExcludeVariations(
	sourcePattern, sourceInclusionPattern, sourceExclusionPattern,
	classPatterns, classInclusionPatterns, classExclusionPatterns string,
	expectedIncludePaths, expectedExcludePaths []string, t *testing.T) {

	plugin, _, err := CheckSourcePathsWithClassPaths(classPatterns, classInclusionPatterns,
		classExclusionPatterns, sourcePattern, sourceInclusionPattern, sourceExclusionPattern, t)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	workSpaceInfoMap, err := plugin.InspectProcessArgs([]string{WorkSpaceCompletePathKeyStr})
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	js, err := ToJsonStringFromMap[map[string]interface{}](workSpaceInfoMap)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	wsi, err := ToStructFromJsonString[WorkSpaceInfo](js)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}
	CheckFilesCopiedToWorkSpace(wsi, t)
}

func CheckFilesCopiedToWorkSpace(wsi WorkSpaceInfo, t *testing.T) {
	expectedFilesList := []string{
		"$WORKSPACE/sources/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Universe.java",
		"$WORKSPACE/sources/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Grid.java",
		"$WORKSPACE/sources/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/Cell.java",
		"$WORKSPACE/sources/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/GridReader.java",
		"$WORKSPACE/sources/gameoflife-core/src/main/java/com/wakaleo/gameoflife/domain/GridWriter.java",
		"$WORKSPACE/classes/pmd-rules.xml",
		"$WORKSPACE/classes/com/wakaleo/gameoflife/domain/Universe.class",
		"$WORKSPACE/classes/com/wakaleo/gameoflife/domain/Cell.class",
		"$WORKSPACE/classes/com/wakaleo/gameoflife/domain/GridReader.class",
		"$WORKSPACE/classes/com/wakaleo/gameoflife/domain/GridWriter.class",
		"$WORKSPACE/classes/com/wakaleo/gameoflife/domain/Grid.class",
		"$WORKSPACE/classes/custom-checkstyle.xml",
		"$WORKSPACE/execFiles/gameoflife-core/target/jacoco.exec",
		"$WORKSPACE/execFiles/gameoflife-web/target/jacoco.exec",
	}

	for _, expectedFile := range expectedFilesList {
		completePath := strings.ReplaceAll(expectedFile, "$WORKSPACE", wsi.WorkSpaceCompletePathKeyStr.Workspace+"/")
		fi, err := os.Stat(completePath)
		if err != nil {
			t.Errorf("Error in CheckFilesCopiedToWorkSpace: %s", err.Error())
		}
	}

}

func (s *SourceTestInfo) isAllRequiredIncludePathsPresent(requiredIncludePaths []string) bool {
	pathSet := make(map[string]struct{})

	for _, sourceInfo := range s.FinalizedSourcesList {
		for _, sourcePath := range sourceInfo.IncludeClassesRelativePathsList {
			pathSet[sourcePath] = struct{}{}
		}
	}

	for _, required := range requiredIncludePaths {
		if _, found := pathSet[required]; !found {
			fmt.Println("Required path not found: ", required)
			return false
		}
	}
	return true
}

func (s *SourceTestInfo) isAllRequiredExcludePathsPresent(requiredExcludePaths []string) bool {
	excludePathSet := make(map[string]struct{})

	for _, sourceInfo := range s.FinalizedSourcesList {
		for _, excludePath := range sourceInfo.ExcludeClassesRelativePathsList {
			excludePathSet[excludePath] = struct{}{}
		}
	}

	for _, required := range requiredExcludePaths {
		if _, found := excludePathSet[required]; !found {
			fmt.Println("Required exclude path not found: ", required)
			return false
		}
	}
	return true
}

func CheckSourcePathsWithClassPaths(classPattern, classInclusionPattern, classExclusionPattern,
	sourcePattern, sourceInclusionPattern, sourceExclusionPattern string,
	t *testing.T) (Plugin, map[string]interface{}, error) {

	args := GetTestNewArgs()
	args.ClassPatterns = classPattern
	args.ClassInclusionPatterns = classInclusionPattern
	args.ClassExclusionPatterns = classExclusionPattern

	args.SourcePattern = sourcePattern
	args.SourceInclusionPattern = sourceInclusionPattern
	args.SourceExclusionPattern = sourceExclusionPattern

	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	sourcesInfo, err := plugin.InspectProcessArgs([]string{FinalizedSourcesListParamKey})
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}
	return plugin, sourcesInfo, err
}

func TestClassPathWithIncludeAndExclude(t *testing.T) {

	classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
		"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	classInclusionPatterns := "**/*.class, **/*.xml"
	classExclusionPatterns := "**/controllers/*.class"

	expectedIncludePaths := []string{
		"com/wakaleo/gameoflife/domain/Cell.class",
		"com/wakaleo/gameoflife/domain/Grid.class",
		"com/wakaleo/gameoflife/domain/GridReader.class",
		"com/wakaleo/gameoflife/domain/GridWriter.class",
		"com/wakaleo/gameoflife/domain/Universe.class",
	}

	expectedExcludePaths := []string{
		"com/wakaleo/gameoflife/webtests/controllers/GameController.class",
		"com/wakaleo/gameoflife/webtests/controllers/HomePageController.class",
	}

	CheckClassPathWithIncludeExcludeVariation(classPatterns, classInclusionPatterns,
		classExclusionPatterns, expectedIncludePaths, expectedExcludePaths, t)

}

func CheckClassPathWithIncludeExcludeVariation(classPatterns, classInclusionPatterns,
	classExclusionPatterns string, expectedIncludePaths, expectedExcludePaths []string, t *testing.T) {

	classesInfo, err := CheckClassPaths(classPatterns, classInclusionPatterns, classExclusionPatterns, t)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	classesJsonStr, err := ToJsonStringFromMap[map[string]interface{}](classesInfo)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	cti, err := ToStructFromJsonString[ClassesTestInfo](classesJsonStr)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	isAllOk := cti.isAllRequiredIncludePathsPresent(expectedIncludePaths)

	if !isAllOk {
		t.Errorf("Error in TestClassPathWithIncludeExclude: Expected paths not found")
	}

	isAllOk = cti.isAllRequiredExcludePathsPresent(expectedExcludePaths)
	if !isAllOk {
		t.Errorf("Error in TestClassPathWithIncludeExclude: Expected exclude paths not found")
	}
}

func CheckClassPaths(classPattern, classInclusionPattern,
	classExclusionPattern string, t *testing.T) (map[string]interface{}, error) {

	args := GetTestNewArgs()
	args.ClassPatterns = classPattern
	args.ClassInclusionPatterns = classInclusionPattern
	args.ClassExclusionPatterns = classExclusionPattern

	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	classesInfo, err := plugin.InspectProcessArgs([]string{ClassesInfoStoreListParamKey})
	return classesInfo, err
}

func GetTestNewArgs() Args {
	args := Args{
		Pipeline:           Pipeline{},
		CoveragePluginArgs: CoveragePluginArgs{PluginToolType: JacocoPluginType},
		EnvPluginInputArgs: EnvPluginInputArgs{ExecFilesPathPattern: TestBuildRootPath},
	}
	args.ExecFilesPathPattern = TestExecPathPattern01
	return args
}

const (
	TestBuildRootPath     = "/opt/hns/test-resources/game-of-life-master/gameoflife-core/target/jacoco.exec"
	TestExecPathPattern01 = "**/target/jacoco.exec"
	TestExecPathPattern02 = "**/target/**.exec"
	TestExecPathPattern03 = "**/jacoco.exec"
	TestExecPathPattern04 = "**/target/jacoco.exec, **/target/**.exec, **/jacoco.exec"

	SourceIncludePathPattern01 = "**/*.class, *.groovy"
	SourceExcludePathPattern01 = "**/src/test/java/**/*.java"
)

//
