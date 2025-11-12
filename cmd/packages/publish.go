// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import (
	stdctx "context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdPackagePublish represents a sub command of Package to publish/upload a package
var CmdPackagePublish = cli.Command{
	Name:        "publish",
	Aliases:     []string{"upload", "push"},
	Usage:       "Publish a package file",
	Description: "Publish a package file to the package registry",
	ArgsUsage:   " ",
	Action:      runPackagePublish,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "package-type",
			Aliases:  []string{"t"},
			Usage:    "Package type (e.g., generic, npm, maven, container)",
			Required: false,
			Value:    "generic",
		},
		&cli.StringFlag{
			Name:     "package-name",
			Aliases:  []string{"n"},
			Usage:    "Package name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "package-version",
			Aliases:  []string{"v"},
			Usage:    "Package version",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Usage:    "Path to file(s) to upload. Can be specified multiple times",
			Required: true,
		},
	}, flags.AllDefaultFlags...),
}

func runPackagePublish(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	packageType := cmd.String("package-type")
	packageName := cmd.String("package-name")
	packageVersion := cmd.String("package-version")
	files := cmd.StringSlice("file")

	if len(files) == 0 {
		return fmt.Errorf("at least one file is required")
	}

	// Validate package name and version
	if strings.TrimSpace(packageName) == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if strings.TrimSpace(packageVersion) == "" {
		return fmt.Errorf("package version cannot be empty")
	}

	// Upload each file
	for _, filePath := range files {
		if err := uploadPackageFile(ctx, packageType, packageName, packageVersion, filePath); err != nil {
			return fmt.Errorf("failed to upload %s: %w", filePath, err)
		}
		fmt.Printf("Uploaded: %s\n", filepath.Base(filePath))
	}

	// Print success message with package URL
	fmt.Printf("\nPackage published successfully!\n")
	fmt.Printf("Package: %s/%s@%s\n", packageType, packageName, packageVersion)
	fmt.Printf("URL: %s/%s/-/packages/%s/%s/%s\n",
		ctx.Login.URL, ctx.Owner, packageType, packageName, packageVersion)

	return nil
}

func uploadPackageFile(ctx *context.TeaContext, packageType, packageName, packageVersion, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileName := filepath.Base(filePath)

	// Construct the API URL
	// Generic package upload endpoint: PUT /api/packages/{owner}/generic/{package_name}/{package_version}/{file_name}
	apiURL := fmt.Sprintf("%s/api/packages/%s/%s/%s/%s/%s",
		strings.TrimSuffix(ctx.Login.URL, "/"),
		ctx.Owner,
		packageType,
		packageName,
		packageVersion,
		fileName,
	)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPut, apiURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+ctx.Login.Token)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = fileInfo.Size()

	// Create HTTP client (with insecure option if needed)
	client := &http.Client{}
	if ctx.Login.Insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
