// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"net/http"
	"net/url"
	"fmt"
)

func (bc *Bconf) LoadHTTP(bconfurl, host, appl string) error {
	v := url.Values{}
	if host != "" {
		v.Set("host", host)
	}
	if appl != "" {
		v.Set("appl", appl)
	}
	res, err := http.PostForm(bconfurl, v)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("bconf.LoadHTTP - response: %v", res.Status)
	}
	bc.LoadJSONReader(res.Body)
	res.Body.Close()
	return nil
}