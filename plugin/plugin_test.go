package plugin

import (
	"context"
	"os"
	"strings"
	"testing"
)

/*
https://github.com/syamv/game-of-life

java -jar jacoco.jar \
    report   ./gameoflife-core/target/jacoco.exec   ./gameoflife-web/target/jacoco.exec   \
    --classfiles ./gameoflife-core/target/classes   \
    --sourcefiles ./gameoflife-core/src/main/java   \
    --html ./gameoflife-core/target/site/jacoco_html   \
    --xml ./gameoflife-core/target/site/jacoco.xml
*/

// rm -rf /opt/hns/test-resources/game-of-life-master/jacoco-workspace--* && BUILD_ROOT_PATH=/opt/hns/test-resources/game-of-life-master go test -count=1 -run ^TestSourcePathWithIncludeAndExclude$
func TestSourcePathWithIncludeAndExclude(t *testing.T) {

	classPatterns := "/opt/hns/test-resources/game-of-life-master/**/target/classes," + " " +
		"/opt/hns/test-resources/game-of-life-master/**/WEB-INF/classes"
	classInclusionPatterns := "**/*.class, **/*.xml"
	classExclusionPatterns := "**/controllers/*.class"

	sourcePatterns := "**/src/main/java"
	sourceInclusionPatterns := "**/*.java, *.groovy"
	sourceExclusionPatterns := "**/controllers/*.java"

	CheckSourceAndClassPathsWithIncludeExcludeVariations(sourcePatterns, sourceInclusionPatterns, sourceExclusionPatterns,
		classPatterns, classInclusionPatterns, classExclusionPatterns, t)

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
	classPatterns, classInclusionPatterns, classExclusionPatterns string, t *testing.T) {

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
		_, err := os.Stat(completePath)
		if err != nil {
			t.Errorf("Error in CheckFilesCopiedToWorkSpace: %s", err.Error())
		}
	}

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
