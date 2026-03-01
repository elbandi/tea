// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	terraform "code.gitea.io/tea/cmd/packages/terraform"

	"github.com/urfave/cli/v3"
)

// CmdPackagesTerraformProvider represents the actions runs command
var CmdPackagesTerraformProvider = cli.Command{
	Name:        "terraform",
	Usage:       "Manage terraform packages",
	Description: "List, view, and manage generic packages in the package registry",
	Commands: []*cli.Command{
		&terraform.CmdPackageTerraformProviderPublish,
	},
}
