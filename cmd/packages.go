// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/packages"

	"github.com/urfave/cli/v3"
)

// CmdPackages represents the package command
var CmdPackages = cli.Command{
	Name:        "packages",
	Aliases:     []string{"package", "pkg"},
	Category:    catEntities,
	Usage:       "Manage packages",
	Description: "Manage packages in the package registry",
	ArgsUsage:   " ", // command does not accept arguments
	Action:      packages.RunPackagesList,
	Commands: []*cli.Command{
		&packages.CmdPackageList,
		&packages.CmdPackageShow,
		&packages.CmdPackageFiles,
		&packages.CmdPackageDelete,
		&packages.CmdPackagePublish,
	},
	Flags: flags.AllDefaultFlags,
}
