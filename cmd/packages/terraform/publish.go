// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package terraform

import (
	"bufio"
	stdctx "context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdPackageTerraformProviderPublish represents a sub command of Package to publish/upload a terraform provider
var CmdPackageTerraformProviderPublish = cli.Command{
	Name:        "publish-provider",
	Aliases:     []string{"upload", "push"},
	Usage:       "Publish a terraform provider",
	Description: "Publish a terraform provider file to the package registry",
	ArgsUsage:   " ",
	Action:      runPackageTerraformProviderPublish,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "provider-name",
			Aliases:  []string{"n"},
			Usage:    "Provider name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "provider-version",
			Aliases:  []string{"v"},
			Usage:    "Provider version",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "path",
			Aliases:  []string{"p"},
			Usage:    "Path to file(s) to upload",
			Required: true,
		},
		&flags.OrgFlag,
	}, flags.AllDefaultFlags...),
}

func runPackageTerraformProviderPublish(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)

	packageName := cmd.String("provider-name")
	packageVersion := cmd.String("provider-version")

	// Validate package name and version
	if strings.TrimSpace(packageName) == "" {
		return fmt.Errorf("provider name cannot be empty")
	}
	if strings.TrimSpace(packageVersion) == "" {
		return fmt.Errorf("provider version cannot be empty")
	}

	sourcePath := cmd.String("path")
	sha256sumsFile := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", packageName, packageVersion)

	result := make([]string, 0)
	result = append(result, sha256sumsFile, sha256sumsFile+".sig")

	// Read the sha256sums file
	file, err := os.Open(path.Join(sourcePath, sha256sumsFile))
	if err != nil {
		return fmt.Errorf("failed to open sha256sums file: %w", err)
	}
	r := regexp.MustCompile(`\s+`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := r.Split(line, -1)

		if len(parts) != 2 {
			break
		}
		result = append(result, parts[1])
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan sha256sums file: %w", err)
	}

	owner := ctx.Owner
	if ctx.Org != "" {
		owner = ctx.Org
	}
	// Upload files
	err = uploadPackageFile(ctx, owner, packageName, packageVersion, sourcePath, result)
	if err != nil {
		return fmt.Errorf("failed to upload terraform provider: %w", err)
	}

	// Print success message with package URL
	fmt.Printf("\nTerraform provider published successfully!\n")
	fmt.Printf("Provider name: %s@%s\n", packageName, packageVersion)
	fmt.Printf("URL: %s/%s/-/packages/%s/%s/%s\n",
		ctx.Login.URL, ctx.Owner, "terraform", packageName, packageVersion)

	return nil
}

func uploadPackageFile(ctx *context.TeaContext, owner, packageName, packageVersion, sourcePath string, filePaths []string) error {
	readers := make(map[string]io.Reader, len(filePaths))
	files := make([]*os.File, 0, len(filePaths))

	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	// Add all files to form
	for _, filePath := range filePaths {
		f, err := os.Open(path.Join(sourcePath, filePath))
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		files = append(files, f)
		readers[filePath] = f
	}

	client := ctx.Login.Client()

	_, err := client.PublishPackageTerraformProvider(owner, packageName, packageVersion, readers)
	if err != nil {
		return err
	}
	fmt.Printf("Terraform provider %s@%s published successfully\n", packageName, packageVersion)

	return nil
}
