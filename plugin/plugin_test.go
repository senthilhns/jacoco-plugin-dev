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

	err := Exec(context.TODO(), args)
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
	err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestEmptyExecPathPattern is accepted")
	}
}

func CheckExecPathPattern(globPattern string, t *testing.T) {
	args := GetTestNewArgs()
	args.ExecFilesPathPattern = globPattern
	err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("CheckExecPathPattern for globPattern: %s" + globPattern + " err == " + err.Error())
	}
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
