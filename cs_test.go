package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestInputFiles(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"sample.csv", false},
		{"not-exist.csv", true},
		{"error.csv", true},
		{"https://raw.githubusercontent.com/tacklehop/csvsearch/main/sample.csv", false},
		{"https://raw.githubusercontent.com/tacklehop/csvsearch/main/non-exist.csv", true},
		{"https://raw.githubusercontent.com/tacklehop/csvsearch/main/error.csv", true},
		{"https://not-exist.com", true},
	}

	for _, c := range cases {
		err := searchCsv(c.in, "")
		if c.wantErr {
			if err == nil {
				t.Errorf("cs %v results in %v", c.in, err)
			}
		} else {
			if err != nil {
				t.Errorf("cs %v results in %v", c.in, err)
			}
		}
	}
}

func TestSearchWords(t *testing.T) {
	testfile := "sample.csv"
	cases := []struct {
		in   string
		want string
	}{
		{"Key1", "[Key1 Word1  ]"},
		{"Word3", "[Key2 Word1 Word2 Word3][Key3 Word3 Word1]"},
	}

	for _, c := range cases {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		stdout := os.Stdout
		os.Stdout = w
		err = searchCsv(testfile, c.in)
		os.Stdout = stdout
		w.Close()

		if err != nil {
			t.Fatalf("%v does not exist\n", testfile)
		}

		var buf bytes.Buffer
		io.Copy(&buf, r)
		s := buf.String()
		s = strings.Replace(s, "\ufeff", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		if s != c.want {
			t.Fatalf("Search result: %#v, want: %#v", s, c.want)
		}
	}

}
