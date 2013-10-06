package bconf

import (
	"encoding/json"
	"sort"
)

type Bconf map[string]interface{}

type sortednode []struct {
	k string
	v interface{}
}

func (bc *Bconf)LoadJson(js []byte) error {
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
	for k, v := range(bc) {
		var nv interface{}
		switch (v.(type)) {
		case map[string]interface{}:
			nv = normalize(Bconf(v.(map[string]interface{})))
		default:
			nv = v
		}
		nb[k] = nv
	}
	return nb
}

func (bc Bconf) get(k ...string) interface{} {
	if len(k) == 0 {
		return bc
	}

	n := bc[k[0]]

	if len(k) == 1 {
		return n
	}

	switch (n.(type)) {
	case Bconf:
		return n.(Bconf).get(k[1:]...)
	}

	return nil
}

func (bc Bconf)GetNode(k ...string) Bconf {
	n := bc.get(k...)
	return Bconf(n.(Bconf))
}

func (bc Bconf)GetString(k ...string) string {
	if len(k) == 0 {
		return ""
	}
	n := bc[k[0]]
	switch (n.(type)) {
	case string:
		if len(k) == 1 {
			return n.(string)
		}
		return ""
	default:
		return n.(Bconf).GetString(k[1:]...)
	}
	return ""
}

func (bc Bconf)ForeachVal(f func(k,v string)) {
	for k, v := range(bc) {
		switch(v.(type)) {
		case string:
			f(k, v.(string))
		}
	}
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
	for sn[i].k, sn[i].v = range(bc) {
		i++
	}
	sort.Sort(sn)
	return sn
}

func (bc Bconf)ForeachSorted(f func(k, v string)) {
	sn := bc.tosortednode()
	for _, s := range(sn) {
		switch(s.v.(type)) {
		case string:
			f(s.k, s.v.(string))
		}		
	}
}
