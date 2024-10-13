package plugin

/*
	Input Parameter	Type

	PLUGIN_TOOL	String (required)
	PLUGIN_FAIL_ON_THRESHOLD	boolean (optional)
	PLUGIN_FAIL_IF_NO_REPORTS	boolean (optional)

	execPattern	PLUGIN_REPORTS_PATH_PATTERN	String (optional)
	classPattern	PLUGIN_CLASS_DIRECTORIES	String (optional)
	exclusionPattern	PLUGIN_CLASS_EXCLUSION_PATTERN	String(optional)
	inclusionPattern	PLUGIN_CLASS_INCLUSION_PATTERN	String(optional)
	skipCopyOfSrcFiles	PLUGIN_SKIP_SOURCE_COPY	boolean(optional)
	sourcePattern	PLUGIN_SOURCE_DIRECTORIES	String(optional)
	sourceInclusionPattern	PLUGIN_SOURCE_INCLUSION_PATTERN	String(optional)
	sourceExclusionPattern	PLUGIN_SOURCE_EXCLUSION_PATTERN	String(optional)
	minimumClassCoverage	PLUGIN_THRESHOLD_CLASS	float(optional)
	minimumMethodCoverage	PLUGIN_THRESHOLD_METHOD	float(optional)
	minimumLineCoverage	PLUGIN_THRESHOLD_LINE	float(optional)
	minimumInstructionCoverage	PLUGIN_THRESHOLD_INSTRUCTION	float(optional)
	minimumBranchCoverage	PLUGIN_THRESHOLD_BRANCH	float(optional)
	minimumComplexityCoverage	PLUGIN_THRESHOLD_COMPLEXITY	int(optional)

	PLUGIN_THRESHOLD_MODULE	float(optional)
	PLUGIN_THRESHOLD_PACKAGE	float(optional)
	PLUGIN_THRESHOLD_FILE	float(optional)
	PLUGIN_THRESHOLD_COMPLEXITY_DENSITY	float(optional)
	PLUGIN_THRESHOLD_LOC	int(optional)

*/

type JacocoPlugin struct {
	CoveragePluginArgs
}

func (p *JacocoPlugin) Init() error {
	LogPrintln(p, "JacocoPlugin Init")
	return nil
}

func (p *JacocoPlugin) DeInit() error {
	LogPrintln(p, "JacocoPlugin DeInit")
	return nil
}

func (p *JacocoPlugin) BuildAndValidateArgs(args Args) error {
	LogPrintln(p, "JacocoPlugin BuildAndValidateArgs")
	return nil
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

//
//
