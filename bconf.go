// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"encoding/json"
	"sort"
	"bufio"
	"os"
	"strings"
	"errors"
	"fmt"
	"io"
)

type Bconf map[string]interface{}

func (bc *Bconf) LoadJson(js []byte) error {
	err := json.Unmarshal(js, bc)

	if bc != nil && len(*bc) > 0 {
		*bc = normalize(*bc)
	}

	return err
}

// Split function for the scanner. This is what removes comments, whitespace and
// empty lines and returns just the valid bconf lines.
func scanLines(data []byte, atEOF bool) (int, []byte, error) {
	totadvance := 0
	for true {
		advance, token, err := bufio.ScanLines(data, atEOF)
		if advance == 0 || err != nil {
			return advance, token, err
		}
		totadvance += advance
		// Eat leading whitespace
		i := 0
		for ; i < len(token); i++ {
			if token[i] != ' ' && token[i] != '\t' {
				break
			}
		}
		token = token[i:]

		// Empty lines and comments.
		if len(token) == 0 || token[0] == '#' {
			data = data[advance:]
			continue
		}
		return totadvance, token, nil
	}
	return -1, nil, nil
}

func (bc *Bconf) LoadConfFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	return bc.LoadConfData(f)
}

func (bc *Bconf) LoadConfData(f io.Reader) error {
	scanner := bufio.NewScanner(f)
	scanner.Split(scanLines)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "include") {
			incpath := strings.Split(t, " \t")
			if len(incpath) < 2 {
				return errors.New(fmt.Sprintf("include without path: %v", t))
			}
			bc.LoadConfFile(incpath[len(incpath) - 1])
		}
		skv := strings.SplitN(t, "=", 2)
		if len(skv) != 2 {
			return errors.New(fmt.Sprintf("malformed bconf line: %v", t))
		}
		k, v := skv[0], skv[1]
		bc.AddValue(strings.Split(k, "."), v)
		
	}
	if err := scanner.Err(); err != nil {
		return err
	}	

	return nil
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
		var nv interface{}
		switch v.(type) {
		case map[string]interface{}:
			nv = normalize(Bconf(v.(map[string]interface{})))
		default:
			nv = v
		}
		nb[k] = nv
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

	switch n.(type) {
	case Bconf:
		return n.(Bconf).get(alloc, k[1:]...)
	}

	return nil
}

func (bc Bconf) GetNode(k ...string) Bconf {
	n := bc.get(false, k...)
	switch n.(type) {
	case Bconf:
		return n.(Bconf)
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
	switch n.(type) {
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

func (bc Bconf) ForeachVal(f func(k, v string)) {
	for k, v := range bc {
		switch v.(type) {
		case string:
			f(k, v.(string))
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
		switch s.v.(type) {
		case string:
			f(s.k, s.v.(string))
		}
	}
}
