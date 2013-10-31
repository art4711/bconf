package bconf_test

import (
	"github.com/art4711/bconf"
	"testing"
	"fmt"
)

var testdata = `{"attr": {"name": {"attrind": "4","attronly": "3","body": "5","id": "0","order": "1","status": "6","suborder": "2"},"order": {"0": "id","1": "order","2": "suborder","3": "attronly","4": "attrind","5": "body","6": "status"}},"conf": {"foo.bar": "17","test.conf": "4711"},"counters": {"attr": {"HEJ": "hej_counter"},"word": {"attr": "attr_counter"}},"opers": {"en": {"AND": "0","NOT": "2","OR": "1","and": "0","not": "2","or": "1"}},"state": {"to_id": "4711"},"stopwords": {"ONE": "","one": ""}}`

func lj(t *testing.T) bconf.Bconf {
	bc := make(bconf.Bconf)
	if err := bc.LoadJson([]byte(testdata)); err != nil {
		t.Fatalf("LoadJson: %v", err)
	}
	return bc
}

func TestJsonLoad(t *testing.T) {
	lj(t)
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
	n.ForeachVal(func(k,v string) {
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
	n.ForeachSorted(func(k,v string) {
		if fmt.Sprint(i) != k {
			t.Errorf("out of order keys: %v != %v", i, k)
		}
		i++
	})
	if i != 7 {
		t.Errorf("too few keys: %v", i)
	}
}
