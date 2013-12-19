package bconf_test

import (
	"fmt"
	"github.com/art4711/bconf"
	"strings"
	"testing"
	"encoding/json"
)

var testjson = `{"attr": {"name": {"attrind": "4","attronly": "3","body": "5","id": "0","order": "1","status": "6","suborder": "2"},"order": {"0": "id","1": "order","2": "suborder","3": "attronly","4": "attrind","5": "body","6": "status"}},"conf": {"foo.bar": "17","test.conf": "4711"},"counters": {"attr": {"HEJ": "hej_counter"},"word": {"attr": "attr_counter"}},"opers": {"en": {"AND": "0","NOT": "2","OR": "1","and": "0","not": "2","or": "1"}},"state": {"to_id": "4711"},"stopwords": {"ONE": "","one": ""}}`

var testdata = `
#comment
   
   #whitespace before comment

node.1=foo
node.2=bar
a.a.a.a.a=b`

func lj(t *testing.T) bconf.Bconf {
	bc := make(bconf.Bconf)
	if err := bc.LoadJson([]byte(testjson)); err != nil {
		t.Fatalf("LoadJson: %v", err)
	}
	return bc
}

func ld(t *testing.T) bconf.Bconf {
	bc := make(bconf.Bconf)
	if err := bc.LoadConfData(strings.NewReader(testdata)); err != nil {
		t.Fatalf("LoadConfData: %v", err)
	}
	return bc
}

func TestLoadJson(t *testing.T) {
	lj(t)
}

func TestLoadData(t *testing.T) {
	ld(t)
}

func TestGet(t *testing.T) {
	bc := lj(t)
	if s := bc.GetString("attr", "name", "attrind"); s != "4" {
		t.Errorf("attr.name.attrind != 4 (%v)", s)
	}
}

func TestForeachVal(t *testing.T) {
	bc := lj(t)
	n := bc.GetNode("conf")
	foo := make(map[string]string)
	foo["foo.bar"] = "17"
	foo["test.conf"] = "4711"
	n.ForeachVal(func(k, v string) {
		x := foo[k]
		if x != v {
			t.Errorf("wrong/missing/repeated k: %v v: %v x: %v", k, v, x)
		}
		foo[k] = "!"
	})
}

func TestForeachSorted(t *testing.T) {
	bc := lj(t)
	n := bc.GetNode("attr", "order")
	i := 0
	n.ForeachSorted(func(k, v string) {
		if fmt.Sprint(i) != k {
			t.Errorf("out of order keys: %v != %v", i, k)
		}
		i++
	})
	if i != 7 {
		t.Errorf("too few keys: %v", i)
	}
}

func TestData(t *testing.T) {
	bc := ld(t)
	if s := bc.GetString("a", "a", "a", "a", "a"); s != "b" {
		t.Errorf("a.a.a.a.a != b (%v)", s)
	}
}

var tm1 = `
a.a=1
a.b=2
a.c.a=3
a.d=4
`

var tm2 = `
a.a=11
a.b.c=22
a.c=33
`

var mexpect = `{"a":{"a":"11","b":{"c":"22"},"c":"33","d":"4"}}`

func TestMerge(t *testing.T) {
	bc := make(bconf.Bconf)
	if err := bc.LoadConfData(strings.NewReader(tm1)); err != nil {
		t.Fatalf("LoadConfData: %v", err)
	}
	b2 := make(bconf.Bconf)
	if err := b2.LoadConfData(strings.NewReader(tm2)); err != nil {
		t.Fatalf("LoadConfData: %v", err)
	}
	bc.Merge(b2)
	j, err := json.Marshal(bc)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	js := string(j)
	if js != mexpect {
		t.Fatalf("expected mismatch: %v != %v", js, mexpect)
	}
}