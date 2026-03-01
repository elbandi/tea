// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// PackagesList prints a listing of packages
func PackagesList(packages []*gitea.Package, output string) {
	t := tableWithHeader(
		"ID",
		"Type",
		"Name",
		"Version",
		"Creator",
		"Created",
	)

	for _, pkg := range packages {
		t.addRow(
			fmt.Sprintf("%d", pkg.ID),
			string(pkg.Type),
			pkg.Name,
			pkg.Version,
			pkg.Creator.UserName,
			FormatTime(pkg.CreatedAt, isMachineReadable(output)),
		)
	}

	t.print(output)
}

// PackageDetails prints detailed information about a package
func PackageDetails(pkg *gitea.Package, output string) {
	if isMachineReadable(output) {
		// For machine-readable output, use the table format
		t := tableWithHeader("Field", "Value")
		t.addRow("ID", fmt.Sprintf("%d", pkg.ID))
		t.addRow("Type", string(pkg.Type))
		t.addRow("Name", pkg.Name)
		t.addRow("Version", pkg.Version)
		t.addRow("Creator", pkg.Creator.UserName)
		t.addRow("Repository", pkg.Repository.FullName)
		t.addRow("Created", FormatTime(pkg.CreatedAt, true))
		t.print(output)
	} else {
		// For human-readable output, use a more readable format
		fmt.Printf("Package Details:\n")
		fmt.Printf("  ID:         %d\n", pkg.ID)
		fmt.Printf("  Type:       %s\n", pkg.Type)
		fmt.Printf("  Name:       %s\n", pkg.Name)
		fmt.Printf("  Version:    %s\n", pkg.Version)
		fmt.Printf("  Creator:    %s\n", pkg.Creator.UserName)
		if pkg.Repository != nil {
			fmt.Printf("  Repository: %s\n", pkg.Repository.FullName)
		}
		fmt.Printf("  Created:    %s\n", FormatTime(pkg.CreatedAt, false))
	}
}

// PackageFilesList prints a listing of package files
func PackageFilesList(files []*gitea.PackageFile, output string) {
	t := tableWithHeader(
		"ID",
		"Name",
		"Size",
	)

	for _, file := range files {
		t.addRow(
			fmt.Sprintf("%d", file.ID),
			file.Name,
			formatByteSize(file.Size),
		)
	}

	t.print(output)
}
