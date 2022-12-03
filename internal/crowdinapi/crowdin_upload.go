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
	"os"
)

// Read more: https://developer.crowdin.com/api/v2/#operation/api.storages.post
type addStorageResp struct {
	Data struct {
		ID       int64  `json:"id"`
		FileName string `json:"fileName"`
	} `json:"data"`
}

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.files.post
type addFileReq struct {
	StorageID   int64  `json:"storageId"`
	Name        string `json:"name"`
	BranchID    int64  `json:"branchId,omitempty"`
	DirectoryID int64  `json:"directoryId,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
}

// Read more: https://developer.crowdin.com/api/v2/#operation/api.projects.files.put
type updateFileReq struct {
	StorageID int64 `json:"storageId"`
}

func (client *Client) uploadToCloudStorage(localPath string, cloudFileName string) (int64, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return 0, errors.New("crowdin api: upload file: " + err.Error())
	}
	defer file.Close()

	resp, err := client.uploadFileExtra("/api/v2/storages", 201, file, map[string]string{"Crowdin-API-FileName": cloudFileName})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data addStorageResp
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, errors.New("crowdin api: POST /api/v2/storages: failed decode JSON: " + err.Error())
	}

	return data.Data.ID, nil
}

func (client *Client) AddFile(projectID string, localPath string, cloudFileName string) error {
	// Add file to Crowdin cloud storage
	cloudFileID, err := client.uploadToCloudStorage(localPath, cloudFileName)
	if err != nil {
		return err
	}

	// Move file from Crowdin cloud storage to Crowdin project
	addReq := addFileReq{
		StorageID: cloudFileID,
		Name:      cloudFileName,
	}

	resp, err := client.sendJSON("POST", "/api/v2/projects/"+projectID+"/files", 201, addReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) UpdateFile(projectID string, localPath string, cloudFileName string, fileID string) error {
	// Add file to Crowdin cloud storage
	cloudFileID, err := client.uploadToCloudStorage(localPath, cloudFileName)
	if err != nil {
		return err
	}

	// Move file from Crowdin cloud storage to Crowdin project
	updateReq := updateFileReq{
		StorageID: cloudFileID,
	}

	resp, err := client.sendJSON("PUT", "/api/v2/projects/"+projectID+"/files/"+fileID, 200, updateReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
