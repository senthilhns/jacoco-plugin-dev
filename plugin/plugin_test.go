package plugin

import (
	"context"
	"os"
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
	ClassesCompletePathsList []string `json:"ClassesCompletePathsList"`
	ClassesRelativePathsList []string `json:"ClassesRelativePathsList"`
}

func (c *ClassesTestInfo) isallrequiredpathsPresent(requiredPaths []string) bool {
	pathSet := make(map[string]struct{})

	for _, path := range c.ClassesRelativePathsList {
		pathSet[path] = struct{}{}
	}

	for _, required := range requiredPaths {
		if _, found := pathSet[required]; !found {
			return false
		}
	}
	return true
}

func TestClassPathWithIncludeExcludeVariations(t *testing.T) {

	classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
		"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	classInclusionPatterns := "**/*.class"
	classExclusionPatterns := "**/controllers/*.class"

	expectedPaths := []string{"gameoflife-build/target/classes",
		"gameoflife-core/target/classes", "gameoflife-web/target/classes"}

	CheckClassPathWithIncludeExcludeVariation(classPatterns, classInclusionPatterns,
		classExclusionPatterns, expectedPaths, t)

}

func CheckClassPathWithIncludeExcludeVariation(classPatterns, classInclusionPatterns,
	classExclusionPatterns string, expectedPaths []string, t *testing.T) {

	//classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
	//	"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	//classInclusionPatterns := "**/*.class"
	//classExclusionPatterns := "**/controllers/*.class"

	_, _ = classExclusionPatterns, classInclusionPatterns

	classesMapList, err := CheckClassPaths(classPatterns, "", "", t)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	classesJsonStr, err := ToJsonStringFromMap[map[string]interface{}](classesMapList)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	cti, err := ToStructFromJsonString[ClassesTestInfo](classesJsonStr)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	isAllOk := cti.isallrequiredpathsPresent(expectedPaths)

	if !isAllOk {
		t.Errorf("Error in TestClassPathWithIncludeExclude: Expected paths not found")
	}
}

func CheckClassPaths(classPattern, classInclusionPattern,
	classExclusionPattern string, t *testing.T) (map[string]interface{}, error) {

	args := GetTestNewArgs()
	args.ClassPatterns = classPattern
	if len(classInclusionPattern) > 0 {
		args.ClassInclusionPatterns = classInclusionPattern
	}
	if len(classExclusionPattern) > 0 {
		args.ClassExclusionPatterns = classExclusionPattern
	}

	plugin, err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestClassPathWithIncludeExclude: %s", err.Error())
	}

	classesListMap, err := plugin.InspectProcessArgs([]string{ClassFilesListParamKey})
	return classesListMap, err
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
)

//
