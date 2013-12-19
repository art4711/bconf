// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"encoding/json"
)

func (bc *Bconf) LoadJson(js []byte) error {
	err := json.Unmarshal(js, bc)

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
	nb := make(Bconf)
	for k, v := range bc {
		if iv, ok := v.(map[string]interface{}); ok {
			nb[k] = normalize(iv)
		} else {
			nb[k] = v;
		}
	}
	return nb
}

