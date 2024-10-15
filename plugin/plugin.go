// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
)

type Plugin interface {
	Init() error
	SetBuildRoot(buildRootPath string) error
	DeInit() error
	ValidateAndProcessArgs(args Args) error
	DoPostArgsValidationSetup(args Args) error
	Run() error
	WriteOutputVariables() error
	PersistResults() error
	GetPluginType() string
	IsQuiet() bool
	InspectProcessArgs(argNamesList []string) (map[string]interface{}, error)
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
	EnvPluginInputArgs
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`
}

type CoveragePluginArgs struct {
	PluginToolType        string `envconfig:"PLUGIN_TOOL"`
	PluginFailOnThreshold bool   `envconfig:"PLUGIN_FAIL_ON_THRESHOLD"`
	PluginFailIfNoReports bool   `envconfig:"PLUGIN_FAIL_IF_NO_REPORTS"`
}

type EnvPluginInputArgs struct {
	ExecFilesPathPattern string `envconfig:"PLUGIN_REPORTS_PATH_PATTERN"`

	ClassPatterns          string `envconfig:"PLUGIN_CLASS_DIRECTORIES"`
	ClassInclusionPatterns string `envconfig:"PLUGIN_CLASS_INCLUSION_PATTERN"`
	ClassExclusionPatterns string `envconfig:"PLUGIN_CLASS_EXCLUSION_PATTERN"`

	SourcePattern          string `envconfig:"PLUGIN_SOURCE_DIRECTORIES"`
	SourceInclusionPattern string `envconfig:"PLUGIN_SOURCE_INCLUSION_PATTERN"`
	SourceExclusionPattern string `envconfig:"PLUGIN_SOURCE_EXCLUSION_PATTERN"`
}

func Exec(ctx context.Context, args Args) (Plugin, error) {

	plugin, err := GetNewPlugin(ctx, args)
	if err != nil {
		return plugin, err
	}

	err = plugin.Init()
	if err != nil {
		return plugin, err
	}
	defer func(p Plugin) {
		err := p.DeInit()
		if err != nil {
			LogPrintln(p, "Error in DeInit: "+err.Error())
		}
	}(plugin)

	err = plugin.ValidateAndProcessArgs(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.DoPostArgsValidationSetup(args)
	if err != nil {
		return plugin, err
	}

	err = plugin.Run()
	if err != nil {
		return plugin, err
	}

	err = plugin.PersistResults()
	if err != nil {
		return plugin, err
	}

	err = plugin.WriteOutputVariables()
	if err != nil {
		return plugin, err
	}

	return plugin, nil
}

const (
	JacocoPluginType    = "jacoco"
	JacocoXmlPluginType = "jacoco-xml"
	CorbeturaPluginType = "corbetura"
)
