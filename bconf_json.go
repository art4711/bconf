// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"encoding/json"
	"io"
)

// Populate a Bconf with data from a byte array that contains json. Returns json parsing errors.
func (bc *Bconf) LoadJson(js []byte) error {
	err := json.Unmarshal(js, bc)

	if bc != nil && len(*bc) > 0 {
		*bc = normalize(*bc)
	}

	return err
}

func (bc *Bconf) LoadJSONReader(r io.Reader) error {
	d := json.NewDecoder(r)
	err := d.Decode(bc)

	if bc != nil && len(*bc) > 0 {
		*bc = normalize(*bc)
	}

	return err
}

/*
 * Normalize what we read from json by changing the data types from
 * map[string]interface{} into Bconf
 */
func normalize(bc Bconf) Bconf {
	for k, v := range bc {
		if iv, ok := v.(map[string]interface{}); ok {
			bc[k] = normalize(iv)
		} else {
			bc[k] = v
		}
	}
	return bc
}
