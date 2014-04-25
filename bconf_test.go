package bconf_test

import (
	"encoding/json"
	"fmt"
	"github.com/art4711/bconf"
	"strings"
	"testing"
)

var testjson = `{"attr": {"name": {"attrind": "4","attronly": "3","body": "5","id": "0","order": "1","status": "6","suborder": "2"},"order": {"0": "id","1": "order","2": "suborder","3": "attronly","4": "attrind","5": "body","6": "status"}},"conf": {"foo.bar": "17","test.conf": "4711"},"counters": {"attr": {"HEJ": "hej_counter"},"word": {"attr": "attr_counter"}},"opers": {"en": {"AND": "0","NOT": "2","OR": "1","and": "0","not": "2","or": "1"}},"state": {"to_id": "4711"},"stopwords": {"ONE": "","one": ""}}`

var testdata = `
#comment
   
   #whitespace before comment

node.1=foo
node.2=bar
a.a.a.a.a=b
some.node.3.a=a
some.node.3.b=b
some.node.3.c=c
some.node.1.a=a
some.node.1.b=b
some.node.1.c=c
some.node.2.a=a
some.node.2.b=b
some.node.2.c=c`

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

func TestForeachNode(t *testing.T) {
	bc := ld(t).GetNode("some", "node")
	nodes := make(map[string]bool)
	nodes["1"] = true
	nodes["2"] = true
	nodes["3"] = true
	i := 0
	bc.ForeachNode(func(n string, b bconf.Bconf) {
		if _, ok := nodes[n]; !ok {
			t.Errorf("node %s not in map", n)
		}
		i++
	})
	if i != 3 {
		t.Errorf("missing nodes: %d/3", i+1)
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

func TestForeachSortedNode(t *testing.T) {
	bc := ld(t).GetNode("some", "node")
	i := 0
	bc.ForeachSortedNode(func(n string, b bconf.Bconf) {
		i++
		if n != fmt.Sprint(i) {
			t.Errorf("out of order nodes: %v != %v", n, i)
		}
	})

	if i != 3 {
		t.Errorf("missing nodes: %d/3", i+1)
	}
}

func TestForeachSortedVal(t *testing.T) {
	bc := lj(t)
	n := bc.GetNode("attr", "order")
	i := 0
	n.ForeachSortedVal(func(k, v string) {
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

type S struct {
	Str string `bconf:"s"`
	I   int
	F   float64
	U   uint
}

var unbconf = `
this.is.s=string
this.is.i=-4711
this.is.f=42.0
this.is.u=4711
`

func TestUnmarshal(t *testing.T) {
	bcThis := make(bconf.Bconf)
	if err := bcThis.LoadConfData(strings.NewReader(unbconf)); err != nil {
		t.Fatalf("LoadConfData: %v", err)
	}

	bc := bcThis.GetNode("this", "is")
	if bc == nil {
		t.Fatalf("GetNode")
	}

	var s S
	if err := bc.Unmarshal(&s); err != nil {
		t.Fatalf("Unmarshal: %v\n", err)
	}
	if s.Str != "string" {
		t.Fatalf("%v != string\n", s.Str)
	}
	if s.I != -4711 {
		t.Fatalf("%v != -4711\n", s.I)
	}
	if s.F != 42.0 {
		t.Fatalf("%v != 42.0\n", s.F)
	}
	if s.U != 4711 {
		t.Fatalf("%v != 4711\n", s.U)
	}
}

func BenchmarkJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bc := make(bconf.Bconf)
		if err := bc.LoadJson([]byte(testjson)); err != nil {
			b.Fatalf("LoadJson: %v", err)
		}
	}
}
