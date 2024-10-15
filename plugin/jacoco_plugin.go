package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type JacocoPlugin struct {
	CoveragePluginArgs
	JacocoPluginParams
	JacocoPluginStateStore
}

type JacocoPluginStateStore struct {
	BuildRootPath string

	ExecFilePathsWithPrefixList []PathWithPrefix

	ClassesInfoStoreList []FilesInfoStore
	FinalizedClassesList []IncludeExcludesMerged

	SourcesInfoStoreList []FilesInfoStore
	FinalizedSourcesList []IncludeExcludesMerged

	JacocoWorkSpaceDir string

	ExecFilesFinalCompletePath []string

	JacocoJarPath string
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

	err = p.CreateNewWorkspace()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}

	err = p.SetJarPath()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}

	return nil
}

func (p *JacocoPlugin) SetJarPath() error {
	p.JacocoJarPath = os.Getenv("JACOCO_JAR_PATH")
	if p.JacocoJarPath == "" {
		p.JacocoJarPath = DefaultJacocoJarPath
	}

	_, err := os.Stat(p.JacocoJarPath)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in SetJarPath: "+err.Error())
		return GetNewError("Error in SetJarPath: " + err.Error())
	}

	return nil
}

func (p *JacocoPlugin) CreateNewWorkspace() error {

	buildRootPath, err := p.GetBuildRootPath()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in CopyClassesToWorkspace: "+err.Error())
		return GetNewError("Error in CopyClassesToWorkspace: " + err.Error())
	}

	jacocoWorkSpaceDir, err := GetRandomJacocoWorkspaceDir(buildRootPath)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}
	p.JacocoWorkSpaceDir = jacocoWorkSpaceDir

	err = CreateDir(p.JacocoWorkSpaceDir)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}

	err = CreateDir(p.GetOutputReportsWorkSpaceDir())
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Init: "+err.Error())
		return err
	}

	return nil
}

func (p *JacocoPlugin) GetWorkspaceDir() string {
	return p.JacocoWorkSpaceDir
}

func (p *JacocoPlugin) InspectProcessArgs(argNamesList []string) (map[string]interface{}, error) {

	m := map[string]interface{}{}
	for _, argName := range argNamesList {
		switch argName {
		case ClassesInfoStoreListParamKey:
			m[argName] = p.ClassesInfoStoreList
		case FinalizedSourcesListParamKey:
			m[argName] = p.SourcesInfoStoreList
		case WorkSpaceCompletePathKeyStr:
			nm := map[string]string{}
			nm["classes"] = p.GetClassesWorkSpaceDir()
			nm["sources"] = p.GetSourcesWorkSpaceDir()
			nm["execFiles"] = p.GetExecFilesWorkSpaceDir()
			nm["workspace"] = p.GetWorkspaceDir()
			m[argName] = nm
		}

	}
	return m, nil
}

func (p *JacocoPlugin) GetBuildRootPath() (string, error) {
	buildRootPath := os.Getenv(BuildRootPathKeyStr)

	if buildRootPath == "" {
		LogPrintln(p, "JacocoPlugin Error in GetBuildRootPath: Build root path is empty")
		return "", GetNewError("Error in GetBuildRootPath: Build")
	}

	return buildRootPath, nil
}

func (p *JacocoPlugin) SetBuildRoot(buildRootPath string) error {

	var err error

	if buildRootPath == "" {
		buildRootPath, err = p.GetBuildRootPath()
		if err != nil {
			LogPrintln(p, "JacocoPlugin Error in SetBuildRoot: "+err.Error())
			return err
		}
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

	err = p.IsSourceArgOk(args)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in ValidateAndProcessArgs: "+err.Error())
		return err
	}

	return nil
}

func (p *JacocoPlugin) GetClassesList() []IncludeExcludesMerged {
	return p.FinalizedClassesList
}

func (p *JacocoPlugin) GetSourcesList() []IncludeExcludesMerged {
	return p.FinalizedSourcesList
}

func (p *JacocoPlugin) DoPostArgsValidationSetup(args Args) error {
	LogPrintln(p, "JacocoPlugin DoPostArgsValidationSetup")

	err := p.CopyClassesToWorkspace()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in DoPostArgsValidationSetup: "+err.Error())
		return err
	}

	err = p.CopySourcesToWorkspace()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in DoPostArgsValidationSetup: "+err.Error())
		return err
	}

	err = p.CopyJacocoExecFilesToWorkspace()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in DoPostArgsValidationSetup: "+err.Error())
		return err
	}

	return nil
}

func (p *JacocoPlugin) GetExecFilesWorkSpaceDir() string {
	return filepath.Join(p.GetWorkspaceDir(), "execFiles")
}

func (p *JacocoPlugin) GetClassesWorkSpaceDir() string {
	return filepath.Join(p.GetWorkspaceDir(), "classes")
}

func (p *JacocoPlugin) GetOutputReportsWorkSpaceDir() string {
	return filepath.Join(p.GetWorkspaceDir(), "reports_dir")
}

func (p *JacocoPlugin) GetSourcesWorkSpaceDir() string {
	return filepath.Join(p.GetWorkspaceDir(), "sources")
}

func (p *JacocoPlugin) CopyJacocoExecFilesToWorkspace() error {
	uniqueDirs, err := p.GetJacocoExecFilesUniqueDirs()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in CopyJacocoExecFilesToWorkspace: "+err.Error())
		return err
	}

	execFilesDir := p.GetExecFilesWorkSpaceDir()
	LogPrintln(p, "JacocoPlugin Copying Exec files to workspace: "+execFilesDir)
	err = CreateDir(execFilesDir)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in CopyJacocoExecFilesToWorkspace: "+err.Error())
		return GetNewError("Error in CopyJacocoExecFilesToWorkspace: " + err.Error())
	}

	for _, dir := range uniqueDirs {
		newDir := filepath.Join(execFilesDir, dir)
		err = CreateDir(newDir)
		if err != nil {
			LogPrintln(p, "JacocoPlugin Error in CopyJacocoExecFilesToWorkspace: "+err.Error())
			return GetNewError("Error in CopyJacocoExecFilesToWorkspace: " + err.Error())
		}
	}

	for _, execFilePathsWithPrefix := range p.ExecFilePathsWithPrefixList {
		relPath := execFilePathsWithPrefix.RelativePath
		srcFilePath := filepath.Join(execFilePathsWithPrefix.CompletePathPrefix, execFilePathsWithPrefix.RelativePath)
		dstFilePath := filepath.Join(execFilesDir, relPath)
		err = CopyFile(srcFilePath, dstFilePath)
		if err != nil {
			LogPrintln(p, "JacocoPlugin Error in CopyJacocoExecFilesToWorkspace: "+err.Error())
			return GetNewError("Error in CopyJacocoExecFilesToWorkspace: " + err.Error())
		}

		p.ExecFilesFinalCompletePath = append(p.ExecFilesFinalCompletePath, dstFilePath)
	}

	return nil
}

func (p *JacocoPlugin) GetJacocoExecFilesUniqueDirs() ([]string, error) {

	uniqueDirMap := map[string]bool{}

	for _, execFilePathsWithPrefix := range p.ExecFilePathsWithPrefixList {
		dir := filepath.Dir(execFilePathsWithPrefix.RelativePath)
		uniqueDirMap[dir] = true
	}

	execFilesDirList := []string{}

	for dir, _ := range uniqueDirMap {
		execFilesDirList = append(execFilesDirList, dir)
	}

	return execFilesDirList, nil
}

func (p *JacocoPlugin) CopyClassesToWorkspace() error {

	classesList := p.GetClassesList()
	if len(classesList) < 1 {
		LogPrintln(p, "JacocoPlugin Error in CopyClassesToWorkspace: No class files to copy")
		return GetNewError("Error in CopyClassesToWorkspace: No class files to copy")
	}

	dstClassesDir := p.GetClassesWorkSpaceDir()
	LogPrintln(p, "JacocoPlugin Copying classes to workspace: "+dstClassesDir)

	err := CreateDir(dstClassesDir)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in CopyClassesToWorkspace: "+err.Error())
		return GetNewError("Error in CopyClassesToWorkspace: " + err.Error())
	}

	for _, classInfo := range classesList {

		err := classInfo.CopyTo(dstClassesDir, p.BuildRootPath)
		if err != nil {
			continue
		}
	}

	return nil
}

func (p *JacocoPlugin) CopySourcesToWorkspace() error {

	if p.SkipCopyOfSrcFiles {
		LogPrintln(p, "JacocoPlugin Skipping copying of source files")
		return nil
	}

	dstSourcesDir := p.GetSourcesWorkSpaceDir()
	LogPrintln(p, "JacocoPlugin Copying sources to workspace: "+dstSourcesDir)
	err := CreateDir(dstSourcesDir)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in CopySourcesToWorkspace: "+err.Error())
		return GetNewError("Error in CopySourcesToWorkspace: " + err.Error())
	}

	sourcesList := p.GetSourcesList()
	for _, sourceInfo := range sourcesList {
		err := sourceInfo.CopySourceTo(dstSourcesDir, p.BuildRootPath)
		if err != nil {
			continue
		}
	}

	return nil
}

func (p *JacocoPlugin) GetClassPatternsStrArray() []string {
	return ToStringArrayFromCsvString(p.ClassPatterns)
}

func (p *JacocoPlugin) GetSourcePatternsStrArray() []string {
	return ToStringArrayFromCsvString(p.SourcePattern)
}

func (p *JacocoPlugin) IsSourceArgOk(args Args) error {
	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	if p.SkipCopyOfSrcFiles {
		LogPrintln(p, "JacocoPlugin Skipping copying of source files")
		return nil
	}

	if args.SourcePattern == "" {
		return GetNewError("Error in IsSourceArgOk: SourcePattern is empty")
	}
	p.SourcePattern = args.SourcePattern
	p.SourceInclusionPattern = args.SourceInclusionPattern
	p.SourceExclusionPattern = args.SourceExclusionPattern

	sourcesInfoStoreList, err :=
		FilterFileOrDirUsingGlobPatterns(p.BuildRootPath, p.GetSourcePatternsStrArray(),
			p.SourceInclusionPattern, p.SourceExclusionPattern, AllSourcesAutoFillGlob)

	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in IsSourceArgOk: "+err.Error())
		return GetNewError("Error in IsSourceArgOk: " + err.Error())
	}

	p.SourcesInfoStoreList = sourcesInfoStoreList
	p.FinalizedSourcesList = MergeIncludeExcludeFileCompletePaths(p.SourcesInfoStoreList)

	return nil

}

func (p *JacocoPlugin) IsClassArgOk(args Args) error {

	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	if args.ClassPatterns == "" {
		return GetNewError("Error in IsClassArgOk: ClassPatterns is empty")
	}
	p.ClassPatterns = args.ClassPatterns
	p.ClassInclusionPatterns = args.ClassInclusionPatterns
	p.ClassExclusionPatterns = args.ClassExclusionPatterns

	classesInfoStoreList, err :=
		FilterFileOrDirUsingGlobPatterns(p.BuildRootPath, p.GetClassPatternsStrArray(),
			p.ClassInclusionPatterns, p.ClassExclusionPatterns, AllClassesAutoFillGlob)

	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in IsClassArgOk: "+err.Error())
		return GetNewError("Error in IsClassArgOk: " + err.Error())
	}

	p.ClassesInfoStoreList = classesInfoStoreList
	p.FinalizedClassesList = MergeIncludeExcludeFileCompletePaths(p.ClassesInfoStoreList)

	if len(p.FinalizedClassesList) < 1 {
		LogPrintln(p, "Error in IsClassArgOk: No class inferred from class patterns")
		return GetNewError("Error in IsClassArgOk: No class inferred from class patterns")
	}
	return nil
}

func (p *JacocoPlugin) IsExecFileArgOk(args Args) error {

	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")

	if args.ExecFilesPathPattern == "" {
		return GetNewError("Error in IsExecFileArgOk: ExecFilesPathPattern is empty")
	}

	execFilesPathList, err := GetAllJacocoExecFilesFromGlobPattern(p.BuildRootPath, args.ExecFilesPathPattern)
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in IsExecFileArgOk: "+err.Error())
		return GetNewError("Error in IsExecFileArgOk: " + err.Error())
	}

	p.ExecFilePathsWithPrefixList = execFilesPathList

	if len(p.ExecFilePathsWithPrefixList) < 1 {
		LogPrintln(p, "JacocoPlugin Error in IsExecFileArgOk: No jacoco exec files found")
		return GetNewError("Error in IsExecFileArgOk: No jacoco exec files found")
	}

	return nil
}

func (p *JacocoPlugin) GetExecFilesList() []PathWithPrefix {
	return p.ExecFilePathsWithPrefixList
}

/*
java -jar jacoco.jar \
    report   ./gameoflife-core/target/jacoco.exec   ./gameoflife-web/target/jacoco.exec   \
    --classfiles ./gameoflife-core/target/classes   \
    --sourcefiles ./gameoflife-core/src/main/java   \
    --html ./gameoflife-core/target/site/jacoco_html   \
    --xml ./gameoflife-core/target/site/jacoco.xml


func main() {
	// Define the command and its arguments
	cmd := exec.Command(
		"java", "-jar", "jacoco.jar",
		"report",
		"./gameoflife-core/target/jacoco.exec",
		"./gameoflife-web/target/jacoco.exec",
		"--classfiles", "./gameoflife-core/target/classes",
		"--sourcefiles", "./gameoflife-core/src/main/java",
		"--html", "./gameoflife-core/target/site/jacoco_html",
		"--xml", "./gameoflife-core/target/site/jacoco.xml",
	)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running command: %v\nOutput: %s", err, output)
	}

	// Print the output if the command succeeds
	fmt.Println(string(output))
}

*/

func (p *JacocoPlugin) Run() error {
	LogPrintln(p, "JacocoPlugin Run")

	args := []string{}

	args = append(args, "/usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java"+" ")
	args = append(args, "-jar"+" "+p.JacocoJarPath+" ")
	args = append(args, p.GetReportArgs()+" ")
	args = append(args, p.GetClassFilesPathArgs()+" ")

	if p.SkipCopyOfSrcFiles == false {
		args = append(args, p.GetSourceFilesPathArgs()+" ")
	}

	args = append(args, p.GetHtmlReportArgs()+" ")
	args = append(args, p.GetXmlReportArgs()+" ")

	cmdStr := strings.Join(args, " ")
	parts := strings.Fields(cmdStr)

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		LogPrintln(p, "JacocoPlugin Error in Run: "+err.Error())
		return GetNewError("Error in Run: " + err.Error())
	} else {
		fmt.Println("Command executed successfully.")
	}

	AnalyzeJacocoXml(p.GetJacocoXmlReportFilePath())
	return nil
}

func (p *JacocoPlugin) GetReportArgs() string {
	reportArg := "report"
	for _, execFilePath := range p.ExecFilesFinalCompletePath {
		reportArg = reportArg + " " + execFilePath
	}
	return reportArg
}

func (p *JacocoPlugin) GetClassFilesPathArgs() string {
	classFilePathArg := "--classfiles"
	classFilePathArg = classFilePathArg + " " + p.GetClassesWorkSpaceDir()
	return classFilePathArg
}

func (p *JacocoPlugin) GetSourceFilesPathArgs() string {
	sourceFilePathArg := "--sourcefiles"
	sourceFilePathArg = sourceFilePathArg + " " + p.GetSourcesWorkSpaceDir()
	return sourceFilePathArg
}

func (p *JacocoPlugin) GetHtmlReportArgs() string {
	htmlReportArg := "--html"
	htmlReportArg = htmlReportArg + " " + p.GetOutputReportsWorkSpaceDir() + "/" + "jacoco_html" + " "
	return htmlReportArg
}

func (p *JacocoPlugin) GetXmlReportArgs() string {
	xmlReportArg := "--xml"
	xmlReportArg = xmlReportArg + " " + p.GetJacocoXmlReportFilePath() + " "
	return xmlReportArg
}

func (p *JacocoPlugin) GetJacocoXmlReportFilePath() string {
	return filepath.Join(p.GetOutputReportsWorkSpaceDir(), "jacoco.xml")
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
	BuildRootPathKeyStr          = "BUILD_ROOT_PATH"
	ClassFilesListParamKey       = "ClassFilesList"
	ClassesInfoStoreListParamKey = "ClassesInfoStoreList"
	FinalizedSourcesListParamKey = "FinalizedSourcesList"
	WorkSpaceCompletePathKeyStr  = "WorkSpaceCompletePathKeyStr"
	AllClassesAutoFillGlob       = "**/*.class"
	AllSourcesAutoFillGlob       = "**/*.java"
	DefaultJacocoJarPath         = "/opt/harness/plugins-deps/jacoco/0.8.12/jacoco.jar"
)

//
//
