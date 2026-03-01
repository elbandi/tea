// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdPackageDelete represents a sub command of Package to delete a package
var CmdPackageDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete a package",
	Description: "Delete a package from the package registry",
	ArgsUsage:   "<package-type> <package-name> <package-version>",
	Action:      runPackageDelete,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "confirm",
			Aliases: []string{"y"},
			Usage:   "Confirm deletion (required)",
		},
		&flags.OrgFlag,
	}, flags.AllDefaultFlags...),
}

func runPackageDelete(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)

	if cmd.NArg() < 3 {
		return fmt.Errorf("package type, name, and version are required")
	}

	if !ctx.Bool("confirm") {
		fmt.Println("Are you sure? Please confirm with -y or --confirm.")
		return nil
	}

	packageType := cmd.Args().Get(0)
	packageName := cmd.Args().Get(1)
	packageVersion := cmd.Args().Get(2)

	client := ctx.Login.Client()

	owner := ctx.Owner
	if ctx.Org != "" {
		owner = ctx.Org
	}
	_, err := client.DeletePackage(owner, packageType, packageName, packageVersion)
	if err != nil {
		return err
	}

	fmt.Printf("Package %s/%s@%s deleted successfully\n", packageType, packageName, packageVersion)
	return nil
}
