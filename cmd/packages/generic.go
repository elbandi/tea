// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	generic "code.gitea.io/tea/cmd/packages/generic"

	"github.com/urfave/cli/v3"
)

// CmdPackagesGeneric represents the actions runs command
var CmdPackagesGeneric = cli.Command{
	Name:        "generic",
	Usage:       "Manage generics packages",
	Description: "List, view, and manage generic packages in the package registry",
	Commands: []*cli.Command{
		&generic.CmdPackageGenericPublish,
		&generic.CmdPackageGenericDeleteFile,
	},
}
