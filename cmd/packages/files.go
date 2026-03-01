// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v3"
)

// CmdPackageFiles represents a sub command of Package to list package files
var CmdPackageFiles = cli.Command{
	Name:        "files",
	Usage:       "List files in a package",
	Description: "List files in a specific package version",
	ArgsUsage:   "<package-type> <package-name> <package-version>",
	Action:      RunPackageFiles,
	Flags: append([]cli.Flag{
		&flags.OrgFlag,
	}, flags.AllDefaultFlags...),
}

// RunPackageFiles lists files in a package
func RunPackageFiles(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)

	if cmd.NArg() < 3 {
		return fmt.Errorf("package type, name, and version are required")
	}

	packageType := cmd.Args().Get(0)
	packageName := cmd.Args().Get(1)
	packageVersion := cmd.Args().Get(2)

	client := ctx.Login.Client()

	owner := ctx.Owner
	if ctx.Org != "" {
		owner = ctx.Org
	}

	files, _, err := client.ListPackageFiles(owner, packageType, packageName, packageVersion)
	if err != nil {
		return err
	}

	print.PackageFilesList(files, ctx.Output)
	return nil
}
