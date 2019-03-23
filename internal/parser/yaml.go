// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ParseYAML takes a YAML file path and tries to parse it as a set
// of peridotctl instructions. It returns the parsed PeridotReq
// object, or an error if unable to load or if YAML contents are
// invalid.
func ParseYAML(filePath string) (*PeridotReq, error) {
	// read in the YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// unmarshal the YAML into the request object
	req := PeridotReq{}
	err = yaml.Unmarshal([]byte(data), &req)
	if err != nil {
		return nil, err
	}

	// now, inspect the request object and its subparts to confirm
	// they are valid
	err = ValidateReq(&req)
	if err != nil {
		return nil, err
	}

	// request is loaded and valid! return it
	return &req, nil
}
