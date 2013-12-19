// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

type Bconf map[string]interface{}

// Add a value to the Bconf.
func (bc Bconf) AddValue(k []string, v string) {
	n := bc.get(true, k[:len(k) - 1]...).(Bconf)
	lastkey := k[len(k) - 1]
	n[lastkey] = v
}

func (bc Bconf) get(alloc bool, k ...string) interface{} {
	if len(k) == 0 {
		return bc
	}

	n, exists := bc[k[0]]

	if !exists {
		if !alloc {
			return nil
		}
		bc[k[0]] = make(Bconf)
		n = bc[k[0]]
	}

	if len(k) == 1 {
		return n
	}

	if bn, ok := n.(Bconf); ok {
		return bn.get(alloc, k[1:]...)
	}
	return nil
}

func (bc Bconf) GetNode(k ...string) Bconf {
	n := bc.get(false, k...)
	if bn, ok := n.(Bconf); ok {
		return bn
	}
	return nil
}

func (bc Bconf) GetString(k ...string) string {
	if len(k) == 0 {
		return ""
	}
	n := bc[k[0]]
	if n == nil {
		return ""
	}
	if bn, ok := n.(Bconf); ok {
		return bn.GetString(k[1:]...)
	}
	if len(k) == 1 {
		return n.(string)
	}
	return ""
}

func (bc Bconf) ForeachVal(f func(k, v string)) {
	for k, v := range bc {
		if s, ok := v.(string); ok {
			f(k, s)
		}
	}
}

