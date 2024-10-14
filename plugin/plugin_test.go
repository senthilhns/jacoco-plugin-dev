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
	args := GetTestNewArgs()

	err := Exec(context.TODO(), args)
	if err != nil {
		t.Errorf("Error in TestExecPathPatterns: %s", err.Error())
	}

}

func TestEmptyExecPathPattern(t *testing.T) {
	args := GetTestNewArgs()
	args.ExecFilesPathPattern = ""
	err := Exec(context.TODO(), args)
	if err == nil {
		t.Errorf("Error in TestEmptyExecPathPattern is accepted")
	}
}

func GetTestNewArgs() Args {
	args := Args{
		Pipeline:           Pipeline{},
		CoveragePluginArgs: CoveragePluginArgs{PluginToolType: JacocoPluginType},
		EnvPluginInputArgs: EnvPluginInputArgs{ExecFilesPathPattern: TestBuildRootPath},
	}

	return args
}

const (
	TestBuildRootPath = "/opt/hns/test-resources/game-of-life-master/gameoflife-core/target/jacoco.exec"
)

//
