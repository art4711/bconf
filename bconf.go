// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

type Bconf map[string]interface{}

// Add a value to the Bconf.
func (bc Bconf) AddValue(k []string, v string) {
	n := bc.get(true, k[:len(k)-1]...).(Bconf)
	lastkey := k[len(k)-1]
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

// Get a bconf node under a varargs key
func (bc Bconf) GetNode(k ...string) Bconf {
	n := bc.get(false, k...)
	if bn, ok := n.(Bconf); ok {
		return bn
	}
	return nil
}

// Return a string from bconf under a key, returns empty string if key not found.
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

// Call the callback function for every node under a Bconf node. The nodes are unsorted.
// When the walk has to be sorted, use ForeachSortedNode instead which is slower.
func (bc Bconf) ForeachNode(f func(name string, bc Bconf)) {
	for n, v := range bc {
		if b, ok := v.(Bconf); ok {
			f(n, b)
		}
	}
}

// Call the callback function for every value (not node) in a Bconf node. The values are unsorted.
// When the walk has to be sorted, use ForeachSortedVal instead which is slower.
func (bc Bconf) ForeachVal(f func(k, v string)) {
	for k, v := range bc {
		if s, ok := v.(string); ok {
			f(k, s)
		}
	}
}

// Merge 'src' into 'dst' overwriting values and branches in 'dst'.
func (dst Bconf) Merge(src Bconf) {
	for k, v := range src {
		dv, exists := dst[k]
		dn, disnode := dv.(Bconf)
		sn, sisnode := v.(Bconf)
		if !exists || !disnode || !sisnode {
			dst[k] = v
		} else {
			dn.Merge(sn)
		}
	}
}
