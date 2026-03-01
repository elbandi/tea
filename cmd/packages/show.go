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

// CmdPackageShow represents a sub command of Package to show package details
var CmdPackageShow = cli.Command{
	Name:        "show",
	Usage:       "Show package details",
	Description: "Show details of a specific package",
	ArgsUsage:   "<package-type> <package-name> <package-version>",
	Action:      RunPackageShow,
	Flags:       flags.AllDefaultFlags,
}

// RunPackageShow shows package details
func RunPackageShow(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if cmd.NArg() < 3 {
		return fmt.Errorf("package type, name, and version are required")
	}

	packageType := cmd.Args().Get(0)
	packageName := cmd.Args().Get(1)
	packageVersion := cmd.Args().Get(2)

	client := ctx.Login.Client()

	pkg, _, err := client.GetPackage(ctx.Owner, packageType, packageName, packageVersion)
	if err != nil {
		return err
	}

	print.PackageDetails(pkg, ctx.Output)
	return nil
}
