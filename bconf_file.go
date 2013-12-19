// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"bufio"
	"os"
	"strings"
	"errors"
	"fmt"
	"io"
)

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
