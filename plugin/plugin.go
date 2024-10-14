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
	Run() error
	WriteOutputVariables() error
	PersistResults() error
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

	err = plugin.ValidateAndProcessArgs(args)
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
