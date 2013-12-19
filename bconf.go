// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"encoding/json"
	"sort"
)

type Bconf map[string]interface{}

func (bc *Bconf) LoadJson(js []byte) error {
	err := json.Unmarshal(js, bc)

	if bc != nil && len(*bc) > 0 {
		*bc = normalize(*bc)
	}

	return err
}


// Add a value to the Bconf.
func (bc Bconf) AddValue(k []string, v string) {
	n := bc.get(true, k[:len(k) - 1]...).(Bconf)
	lastkey := k[len(k) - 1]
	n[lastkey] = v
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

type sortednode []struct {
	k string
	v interface{}
}

func (sn sortednode) Len() int {
	return len(sn)
}

func (sn sortednode) Less(i, j int) bool {
	return sn[i].k < sn[j].k
}

func (sn sortednode) Swap(i, j int) {
	sn[i], sn[j] = sn[j], sn[i]
}

func (bc Bconf) tosortednode() sortednode {
	sn := make(sortednode, len(bc))
	i := 0
	for sn[i].k, sn[i].v = range bc {
		i++
	}
	sort.Sort(sn)
	return sn
}

func (bc Bconf) ForeachSorted(f func(k, v string)) {
	sn := bc.tosortednode()
	for _, s := range sn {
		if sv, ok := s.v.(string); ok {
			f(s.k, sv)
		}
	}
}
