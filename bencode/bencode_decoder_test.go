package bencode

import (
	"strings"
	"testing"
)

func TestDecodeInt(t *testing.T) {
	d := NewDecoder(strings.NewReader("i42e"))
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := v.(int)
	if !ok {
		t.Fatalf("expected int, got %T", v)
	}
	if n != 42 {
		t.Fatalf("expected 42, got %d", n)
	}
}

func TestDecodeString(t *testing.T) {
	d := NewDecoder(strings.NewReader("4:spam"))
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s != "spam" {
		t.Fatalf("expected spam, got %q", s)
	}
}

func TestDecodeList(t *testing.T) {
	d := NewDecoder(strings.NewReader("l4:spam4:eggsi123ee"))
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := v.([]BValue)
	if !ok {
		t.Fatalf("expected []BValue, got %T", v)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list))
	}
	if list[0].(string) != "spam" || list[1].(string) != "eggs" || list[2].(int) != 123 {
		t.Fatalf("unexpected list contents: %#v", list)
	}
}

func TestDecodeDict(t *testing.T) {
	// dictionary keys are strings
	d := NewDecoder(strings.NewReader("d3:cow3:moo4:spam4:eggse"))
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := v.(map[string]BValue)
	if !ok {
		t.Fatalf("expected map[string]BValue, got %T", v)
	}
	if m["cow"].(string) != "moo" {
		t.Fatalf("expected cow->moo, got %v", m["cow"])
	}
	if m["spam"].(string) != "eggs" {
		t.Fatalf("expected spam->eggs, got %v", m["spam"])
	}
}

func TestDecodeNested(t *testing.T) {
	in := "d4:dictd3:key5:valuee4:listl3:one3:twoi3eee"
	d := NewDecoder(strings.NewReader(in))
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := v.(map[string]BValue)
	// check inner dict
	inner := m["dict"].(map[string]BValue)
	if inner["key"].(string) != "value" {
		t.Fatalf("unexpected inner dict value: %v", inner["key"])
	}
	// check inner list
	lst := m["list"].([]BValue)
	if lst[0].(string) != "one" || lst[1].(string) != "two" || lst[2].(int) != 3 {
		t.Fatalf("unexpected inner list: %#v", lst)
	}
}

func TestDecodeErrors(t *testing.T) {
	cases := []string{
		"i12",      // missing ending 'e'
		"5:abc",    // length 5 but only 3 bytes
		"x5:abcde", // invalid prefix
	}
	for _, c := range cases {
		d := NewDecoder(strings.NewReader(c))
		if _, err := d.Decode(); err == nil {
			t.Fatalf("expected error for input %q", c)
		}
	}
}
