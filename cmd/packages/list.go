// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	stdctx "context"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdPackageList represents a sub command of Package to list packages
var CmdPackageList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List packages",
	Description: "List packages in the package registry",
	ArgsUsage:   " ", // command does not accept arguments
	Action:      RunPackagesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "Filter by package type (e.g., generic, npm, maven, container)",
		},
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Search query to filter packages by name",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunPackagesList lists packages
func RunPackagesList(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	client := ctx.Login.Client()

	packages, _, err := client.ListPackages(ctx.Owner, gitea.ListPackagesOptions{
		ListOptions: flags.GetListOptions(),
	})
	if err != nil {
		return err
	}

	print.PackagesList(packages, ctx.Output)
	return nil
}
