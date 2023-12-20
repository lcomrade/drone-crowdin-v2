// Copyright (C) 2022-2023 Leonid Maslakov.

// This file is part of drone-crowdin-v2.

// drone-crowdin-v2 is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// drone-crowdin-v2 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with drone-crowdin-v2.
// If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lcomrade/drone-crowdin-v2/internal/crowdin"
)

const (
	author          = "Leonid Maslakov <root@lcomrade.su>"
	downloadSources = "https://github.com/lcomrade/drone-crowdin-v2"

	version = "20231220"
)

func exitOnError(a ...interface{}) {
	fmt.Fprint(os.Stderr, "error: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func getBoolVal(envName string) bool {
	valStr := os.Getenv(envName)
	if valStr != "" && valStr != "false" && valStr != "true" {
		exitOnError("bad " + envName + " parameter value")
	}

	val := false
	if valStr == "true" {
		val = true
	}

	return val
}

func targetUpload(client *crowdin.Client, projectID string) {
	// Get files list from parameters
	filesList := os.Getenv("PLUGIN_UPLOAD_FILES")

	if filesList == "" {
		exitOnError("upload files list cannot be empty")
	}

	targetFiles := make(map[string]string)
	err := json.Unmarshal([]byte(filesList), &targetFiles)
	if err != nil {
		exitOnError("failed read upload files list:", err.Error())
	}

	if len(filesList) == 0 {
		exitOnError("upload files list cannot be empty")
	}

	cloudBadSymbols := string(crowdin.BadSymbols)
	for localPath, cloudName := range targetFiles {
		if localPath == "" {
			exitOnError("local file path cannot be empty")
		}

		if cloudName == "" {
			exitOnError("Crowdin file name cannot be empty")
		}

		if strings.ContainsAny(cloudName, cloudBadSymbols) {
			exitOnError("Crowdin file name cannot contain '"+cloudBadSymbols+"':", cloudName)
		}
	}

	// Add or update files
	for localPath, cloudName := range targetFiles {
		// Check file exist in Crowdin
		fileID, err := client.FindFileId(projectID, cloudName)
		if err != nil {
			exitOnError(err)
		}

		// Add if not exist
		if fileID == "" {
			fmt.Println("- Add:   ", localPath, "->", cloudName)
			err = client.AddFile(projectID, localPath, cloudName)
			if err != nil {
				exitOnError(err)
			}

			// Else update file
		} else {
			fmt.Println("- Update:", localPath, "->", cloudName)
			err = client.UpdateFile(projectID, localPath, cloudName, fileID)
			if err != nil {
				exitOnError(err)
			}
		}
	}
}

func targetDownload(client *crowdin.Client, projectID string) {
	// Get download parameters
	downloadTo := os.Getenv("PLUGIN_DOWNLOAD_TO")
	if downloadTo == "" {
		exitOnError("empty 'download to' parameter")
	}

	skipUntranslatedStrings := getBoolVal("PLUGIN_DOWNLOAD_SKIP_UNTRANSLATED_STRINGS")
	skipUntranslatedFiles := getBoolVal("PLUGIN_DOWNLOAD_SKIP_UNTRANSLATED_FILES")
	exportApprovedOnly := getBoolVal("PLUGIN_DOWNLOAD_EXPORT_APPROVED_ONLY")

	// Download
	extracted, err := client.Download(downloadTo, projectID, skipUntranslatedStrings, skipUntranslatedFiles, exportApprovedOnly)
	if err != nil {
		exitOnError(err)
	}

	for _, part := range extracted {
		fmt.Println("- Extract:", part, "->", filepath.Join(downloadTo, part))
	}
}

func main() {
	fmt.Println("Drone plugin author:", author)
	fmt.Println("Drone plugin source code:", downloadSources)
	fmt.Println("Drone plugin version:", version)

	var err error
	tipProjectID := false

	// Get parameters
	target := os.Getenv("PLUGIN_TARGET")
	key := os.Getenv("PLUGIN_CROWDIN_KEY")
	projectID := os.Getenv("PLUGIN_PROJECT_ID")
	projectName := os.Getenv("PLUGIN_PROJECT_NAME")

	if key == "" {
		exitOnError("empty Crowdin API key")
	}

	if projectID == "" && projectName == "" {
		exitOnError("Crowdin project ID or name not set")
	}

	// Prepare Crowdin API client
	crowdin := crowdin.NewClient(key)

	// Get project ID if need
	if projectID == "" {
		projectID, err = crowdin.FindProjectIdByName(projectName)
		if err != nil {
			exitOnError(err)
		}

		tipProjectID = true
	}

	fmt.Println("Crowdin project ID:", projectID)

	// Run
	switch target {
	case "upload":
		targetUpload(crowdin, projectID)
	case "download":
		targetDownload(crowdin, projectID)
	default:
		exitOnError("unknown target '" + target + "' (possible targets: upload, download)")
	}

	// Tips
	if tipProjectID {
		fmt.Println("TIP: Use the 'project ID' parameter instead of 'project name'. Your Project ID:", projectID)
	}
}
