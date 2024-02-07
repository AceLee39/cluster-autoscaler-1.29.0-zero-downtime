/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package civocloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// InstanceSize represents an available size for instances to launch
type InstanceSize struct {
	Type              string `json:"type,omitempty"`
	Name              string `json:"name,omitempty"`
	NiceName          string `json:"nice_name,omitempty"`
	CPUCores          int    `json:"cpu_cores,omitempty"`
	GPUCount          int    `json:"gpu_count,omitempty"`
	GPUType           string `json:"gpu_type,omitempty"`
	RAMMegabytes      int    `json:"ram_mb,omitempty"`
	DiskGigabytes     int    `json:"disk_gb,omitempty"`
	TransferTerabytes int    `json:"transfer_tb,omitempty"`
	Description       string `json:"description,omitempty"`
	Selectable        bool   `json:"selectable,omitempty"`
}

// ListInstanceSizes returns all available sizes of instances
// TODO: Rename to Size because this return all size (k8s, vm, database, kfaas)
func (c *Client) ListInstanceSizes() ([]InstanceSize, error) {
	resp, err := c.SendGetRequest("/v2/sizes")
	if err != nil {
		return nil, decodeError(err)
	}

	sizes := make([]InstanceSize, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&sizes); err != nil {
		return nil, err
	}

	return sizes, nil
}

// FindInstanceSizes finds a instance size name by either part of the ID or part of the name
func (c *Client) FindInstanceSizes(search string) (*InstanceSize, error) {
	instanceSize, err := c.ListInstanceSizes()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := InstanceSize{}

	for _, value := range instanceSize {
		if value.Name == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) {
			if !exactMatch {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}
