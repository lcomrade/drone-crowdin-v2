// Copyright (C) 2022 Leonid Maslakov.

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

package crowdinapi

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.getMany
type projectsListResp struct {
	Data []struct {
		Data struct {
			ID                int64    `json:"id"`
			SourceLanguageID  string   `json:"sourceLanguageId"`
			TargetLanguageIds []string `json:"targetLanguageIds"`
			Name              string   `json:"name"`
			//CName string `json:"cname"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.files.getMany
type listFilesResp struct {
	Data []struct {
		Data struct {
			ID          int64  `json:"id"`
			ProjectID   int64  `json:"projectId"`
			BranchID    int64  `json:"branchId"`
			DirectoryID int64  `json:"directoryId"`
			Name        string `json:"name"`
			Title       string `json:"title"`
			Type        string `json:"type"`
			Path        string `json:"path"`
			Status      string `json:"status"`
		} `json:"data"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

func (client *Client) FindProjectIdByName(name string) (string, error) {
	for offset := 0; ; offset++ {
		resp, err := client.get("/api/v2/projects?limit="+strconv.Itoa(paginationLimit)+"&offset="+strconv.Itoa(offset), 200)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var data projectsListResp
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", errors.New("crowdin api: GET /api/v2/projects: failed decode JSON: " + err.Error())
		}

		if len(data.Data) == 0 {
			break
		}

		for _, part := range data.Data {
			if part.Data.Name == name {
				return strconv.FormatInt(part.Data.ID, 10), nil
			}
		}
	}

	return "", errors.New("crowdin api: project name not found: " + name)
}

func (client *Client) FindFileId(projectID string, fileName string) (string, error) {
	for offset := 0; ; offset++ {
		resp, err := client.get("/api/v2/projects/"+projectID+"/files?limit="+strconv.Itoa(paginationLimit)+"&offset="+strconv.Itoa(offset), 200)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var data listFilesResp
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", errors.New("crowdin api: GET /api/v2/projects/" + projectID + "/files: failed decode JSON: " + err.Error())
		}

		if len(data.Data) == 0 {
			break
		}

		for _, part := range data.Data {
			if part.Data.Name == fileName {
				return strconv.FormatInt(part.Data.ID, 10), nil
			}
		}
	}

	return "", nil
}
