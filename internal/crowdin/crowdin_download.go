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

package crowdin

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.translations.builds.post
type buildProjectReq struct {
	BranchID                int64    `json:"branchId,omitempty"`
	TargetLanguageIds       []string `json:"targetLanguageIds,omitempty"`
	SkipUntranslatedStrings bool     `json:"skipUntranslatedStrings"`
	SkipUntranslatedFiles   bool     `json:"skipUntranslatedFiles"`
	ExportApprovedOnly      bool     `json:"exportApprovedOnly"`
}

type buildProjectResp struct {
	Data struct {
		ID       int    `json:"id"`
		Status   string `json:"status"`
		Progress int    `json:"progress"`
	} `json:"data"`
}

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.translations.builds.download.download
type downloadProjectResp struct {
	Data struct {
		URL string `json:"url"`
		//ExpireIn time.Time `json:"expireIn"`
	} `json:"data"`
}

func (client *Client) Download(destDir string, projectID string, skipUntranslatedStrings bool, skipUntranslatedFiles bool, exportApprovedOnly bool) ([]string, error) {
	// Start build Crowdin project
	buildReq := buildProjectReq{
		SkipUntranslatedStrings: skipUntranslatedStrings,
		SkipUntranslatedFiles:   skipUntranslatedFiles,
		ExportApprovedOnly:      exportApprovedOnly,
	}

	respBuild, err := client.sendJSON("POST", "/api/v2/projects/"+projectID+"/translations/builds", 201, buildReq)
	if err != nil {
		return nil, err
	}
	defer respBuild.Body.Close()

	var dataBuild buildProjectResp
	err = json.NewDecoder(respBuild.Body).Decode(&dataBuild)
	if err != nil {
		return nil, errors.New("crowdin api: POST /api/v2/projects/" + projectID + "/translations/builds: failed decode JSON: " + err.Error())
	}

	buildID := strconv.Itoa(dataBuild.Data.ID)

	// Download build when it finished
	for i := 0; i < 6; i++ {
		// Wait
		time.Sleep(5 * time.Second)

		// Check build status
		var respDl *http.Response
		respDl, err = client.get("/api/v2/projects/"+projectID+"/translations/builds/"+buildID+"/download", 200)
		if err != nil {
			continue
		}
		defer respDl.Body.Close()

		var dataDl downloadProjectResp
		err = json.NewDecoder(respDl.Body).Decode(&dataDl)
		if err != nil {
			return nil, errors.New("crowdin api: GET /api/v2/projects/" + projectID + "/translations/builds/" + buildID + "/download: failed decode JSON: " + err.Error())
		}

		// Request zip archive from server
		tmpFile, err := client.dlToTmpFile(dataDl.Data.URL)
		if err != nil {
			return nil, err
		}
		defer os.Remove(tmpFile)

		// Extract zip archive
		extracted, err := unzip(tmpFile, destDir)
		if err != nil {
			return nil, errors.New("crowdin api: failed extract archive: " + err.Error())
		}
		return extracted, nil
	}

	return nil, err
}
