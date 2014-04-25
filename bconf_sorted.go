// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"sort"
)

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

// Call the callback function for every node under a Bconf node. The nodes are sorted on their name.
// When the walk doesn't need to be sorted, used ForeachNode instead which is faster.
func (bc Bconf) ForeachSortedNode(f func(name string, bc Bconf)) {
	sn := bc.tosortednode()
	for _, v := range sn {
		f(v.k, v.v.(Bconf))
	}
}

// Call the callback function for every value (not node) in a Bconf node. The values are sorted on their key.
// When the walk doesn't need to be sorted, use ForeachVal instead which is faster.
func (bc Bconf) ForeachSortedVal(f func(k, v string)) {
	sn := bc.tosortednode()
	for _, s := range sn {
		if sv, ok := s.v.(string); ok {
			f(s.k, sv)
		}
	}
}
