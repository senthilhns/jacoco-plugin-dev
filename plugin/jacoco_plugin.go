package plugin

import (
	"fmt"
	"os"
)

type JacocoPlugin struct {
	CoveragePluginArgs
	JacocoPluginParams
	JacocoPluginStateStore
}

type JacocoPluginStateStore struct {
	BuildRootPath     string
	ExecFilePathsList []string
	ClassFilesList    []string

	ClassesCompletePathsList []string
	ClassesRelativePathsList []string
}

type JacocoPluginParams struct {
	ExecPattern string `envconfig:"PLUGIN_REPORTS_PATH_PATTERN"`

	ClassPatterns          string `envconfig:"PLUGIN_CLASS_DIRECTORIES"`
	ClassInclusionPatterns string `envconfig:"PLUGIN_CLASS_INCLUSION_PATTERN"`
	ClassExclusionPatterns string `envconfig:"PLUGIN_CLASS_EXCLUSION_PATTERN"`

	SourcePattern          string `envconfig:"PLUGIN_SOURCE_DIRECTORIES"`
	SourceInclusionPattern string `envconfig:"PLUGIN_SOURCE_INCLUSION_PATTERN"`
	SourceExclusionPattern string `envconfig:"PLUGIN_SOURCE_EXCLUSION_PATTERN"`

	SkipCopyOfSrcFiles bool `envconfig:"PLUGIN_SKIP_SOURCE_COPY"`

	MinimumInstructionCoverage float64 `envconfig:"PLUGIN_THRESHOLD_INSTRUCTION"`
	MinimumBranchCoverage      float64 `envconfig:"PLUGIN_THRESHOLD_BRANCH"`
	MinimumComplexityCoverage  int     `envconfig:"PLUGIN_THRESHOLD_COMPLEXITY"`
	MinimumLineCoverage        float64 `envconfig:"PLUGIN_THRESHOLD_LINE"`
	MinimumMethodCoverage      float64 `envconfig:"PLUGIN_THRESHOLD_METHOD"`
	MinimumClassCoverage       float64 `envconfig:"PLUGIN_THRESHOLD_CLASS"`
}

func (p *JacocoPlugin) Init() error {
	LogPrintln(p, "JacocoPlugin Init")

	err := p.SetBuildRoot("")
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}

	return nil
}

func (p *JacocoPlugin) InspectProcessArgs(
	argNamesList []string) (map[string]interface{}, error) {

	m := map[string]interface{}{}
	for _, argName := range argNamesList {
		switch argName {
		case ClassFilesListParamKey:
			m["ClassesCompletePathsList"] = p.ClassesCompletePathsList
			m["ClassesRelativePathsList"] = p.ClassesRelativePathsList
		}
	}
	return m, nil
}

func (p *JacocoPlugin) SetBuildRoot(buildRootPath string) error {

	if buildRootPath == "" {
		buildRootPath = os.Getenv(BuildRootPathKeyStr)
	}

	ok, err := IsDirExists(buildRootPath)

	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in SetBuildRoot: "+err.Error())
		return err
	}

	if !ok {
		LogPrintln(p, "JacocoPlugin Error in SetBuildRoot: Build root path does not exist")
		return GetNewError("Error in SetBuildRoot: Build root path does not exist")
	}

	p.BuildRootPath = buildRootPath
	return nil
}

func (p *JacocoPlugin) DeInit() error {
	LogPrintln(p, "JacocoPlugin DeInit")
	return nil
}

func (p *JacocoPlugin) ValidateAndProcessArgs(args Args) error {
	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	err := p.IsExecFileArgOk(args)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in ValidateAndProcessArgs: "+err.Error())
		return err
	}

	err = p.IsClassArgOk(args)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in ValidateAndProcessArgs: "+err.Error())
		return err
	}
	return nil
}

func (p *JacocoPlugin) GetClassPatternsStrArray() []string {
	return ToStringArrayFromCsvString(p.ClassPatterns)
}

func (p *JacocoPlugin) IsClassArgOk(args Args) error {

	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	if args.ClassPatterns == "" {
		return GetNewError("Error in IsClassArgOk: ClassPatterns is empty")
	}
	p.ClassPatterns = args.ClassPatterns
	p.ClassInclusionPatterns = args.ClassInclusionPatterns
	p.ClassExclusionPatterns = args.ClassExclusionPatterns

	classesCompletePathsList, classesRelativePathsList, err :=
		FilterFileOrDirUsingGlobPatterns(p.BuildRootPath, p.GetClassPatternsStrArray(),
			p.ClassInclusionPatterns, p.ClassExclusionPatterns)

	fmt.Println("CCCCCCCCCCCCC", classesRelativePathsList)

	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in IsClassArgOk: "+err.Error())
		return GetNewError("Error in IsClassArgOk: " + err.Error())
	}

	p.ClassesCompletePathsList, p.ClassesRelativePathsList = classesCompletePathsList, classesRelativePathsList
	return nil
}

func (p *JacocoPlugin) IsExecFileArgOk(args Args) error {
	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	if args.ExecFilesPathPattern == "" {
		return GetNewError("Error in IsExecFileArgOk: ExecFilesPathPattern is empty")
	}

	execFilesPathList, err := GetAllEntriesFromGlobPattern(p.BuildRootPath, args.ExecFilesPathPattern)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in IsExecFileArgOk: "+err.Error())
		return GetNewError("Error in IsExecFileArgOk: " + err.Error())
	}

	p.ExecFilePathsList = execFilesPathList

	if len(p.ExecFilePathsList) < 1 {
		LogPrintln(p, "JacocoPlugin Error in IsExecFileArgOk: No jacoco exec files found")
		return GetNewError("Error in IsExecFileArgOk: No jacoco exec files found")
	}

	return nil
}

func (p *JacocoPlugin) GetExecFilesList() []string {
	return p.ExecFilePathsList
}

func (p *JacocoPlugin) Run() error {
	LogPrintln(p, "JacocoPlugin Run")
	return nil
}

func (p *JacocoPlugin) PersistResults() error {
	LogPrintln(p, "JacocoPlugin StoreResults")
	return nil
}

func (p *JacocoPlugin) WriteOutputVariables() error {
	LogPrintln(p, "JacocoPlugin WriteOutputVariables")
	return nil
}

// Attr methods follow

func (p *JacocoPlugin) IsQuiet() bool {
	return false
}

func (p *JacocoPlugin) GetPluginType() string {
	return JacocoPluginType
}

const (
	BuildRootPathKeyStr    = "BUILD_ROOT_PATH"
	ClassFilesListParamKey = "ClassFilesList"
)

//
//
