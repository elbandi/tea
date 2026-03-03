// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package generic

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdPackageGenericDeleteFile represents a sub command of Package to delete a package
var CmdPackageGenericDeleteFile = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete a package file",
	Description: "Delete a package file from the package registry",
	ArgsUsage:   "<package-name> <package-version> <file-name>",
	Action:      runPackageGenericDeleteFile,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "confirm",
			Aliases: []string{"y"},
			Usage:   "Confirm deletion (required)",
		},
		&flags.OrgFlag,
	}, flags.AllDefaultFlags...),
}

func runPackageGenericDeleteFile(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)

	if cmd.NArg() < 3 {
		return fmt.Errorf("package name, version, and file are required")
	}

	if !ctx.Bool("confirm") {
		fmt.Println("Are you sure? Please confirm with -y or --confirm.")
		return nil
	}

	packageName := cmd.Args().Get(0)
	packageVersion := cmd.Args().Get(1)
	packageFile := cmd.Args().Get(2)

	client := ctx.Login.Client()

	owner := ctx.Owner
	if ctx.Org != "" {
		owner = ctx.Org
	}
	_, err := client.DeletePackageGenericFile(owner, packageName, packageVersion, packageFile)
	if err != nil {
		return err
	}

	fmt.Printf("File %s from package generic %s@%s deleted successfully\n", packageFile, packageName, packageVersion)
	return nil
}
