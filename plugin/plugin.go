// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
)

type Plugin interface {
	Init() error
	DeInit() error
	BuildAndValidateArgs(args Args) error
	Run() error
	WriteOutputVariables() error
	PersistResults() error

	/* Attribute Methods */
	GetPluginType() string
	IsQuiet() bool
}

func GetNewPlugin(ctx context.Context, args Args) (Plugin, error) {

	pluginToolType := args.PluginToolType

	switch pluginToolType {
	case JacocoPluginType:
		return &JacocoPlugin{}, nil

	default:
		return nil, GetNewError("Unknown plugin type: " + pluginToolType)
	}
}

type Args struct {
	Pipeline
	CoveragePluginArgs

	//// Level defines the plugin log level.
	//Level string `envconfig:"PLUGIN_LOG_LEVEL"`
	//
	//// TODO replace or remove
	//Param1 string `envconfig:"PLUGIN_PARAM1"`
	//Param2 string `envconfig:"PLUGIN_PARAM2"`
}

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

type CoveragePluginArgs struct {
	Level                 string `envconfig:"PLUGIN_LOG_LEVEL"`
	PluginToolType        string `envconfig:"PLUGIN_TOOL"`
	PluginFailOnThreshold bool   `envconfig:"PLUGIN_FAIL_ON_THRESHOLD"`
	PluginFailIfNoReports bool   `envconfig:"PLUGIN_FAIL_IF_NO_REPORTS"`
}

func Exec(ctx context.Context, args Args) error {

	plugin, err := GetNewPlugin(ctx, args)
	if err != nil {
		return err
	}

	err = plugin.Init()
	if err != nil {
		return err
	}
	defer func(p Plugin) {
		err := p.DeInit()
		if err != nil {
			LogPrintln(p, "Error in DeInit: "+err.Error())
		}
	}(plugin)

	err = plugin.BuildAndValidateArgs(args)
	if err != nil {
		return err
	}

	err = plugin.Run()
	if err != nil {
		return err
	}

	err = plugin.PersistResults()
	if err != nil {
		return err
	}

	err = plugin.WriteOutputVariables()
	if err != nil {
		return err
	}

	return nil
}

const (
	JacocoPluginType    = "jacoco"
	JacocoXmlPluginType = "jacoco-xml"
	CorbeturaPluginType = "corbetura"
)
