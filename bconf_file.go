// Copyright 2013 Artur Grabowski. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
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

// Populate a Bconf with data from the bconf-formatted file provided here as a name.
func (bc *Bconf) LoadConfFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	return bc.LoadConfData(f)
}

// Populate a Bconf with data from the bconf-formatted io.Reader.
func (bc *Bconf) LoadConfData(f io.Reader) error {
	scanner := bufio.NewScanner(f)
	scanner.Split(scanLines)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "include") {
			incpath := strings.TrimLeftFunc(t[7:], unicode.IsSpace)
			if incpath == "" {
				return fmt.Errorf("include without path: '%v'", t)
			}
			if err := bc.LoadConfFile(incpath); err != nil {
				return err
			}
		} else {
			skv := strings.SplitN(t, "=", 2)
			if len(skv) != 2 {
				return fmt.Errorf("malformed bconf line: %v", t)
			}
			k, v := skv[0], skv[1]
			bc.AddValue(strings.Split(k, "."), v)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
