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
		name	string
		in      string
		wantErr bool
	}{
		{"file: normal", "sample.csv", false},
		{"file: not exist", "not-exist.csv", true},
		{"file: non csv format", "error.csv", true},
		{"http: normal", "https://raw.githubusercontent.com/tacklehop/csvsearch/main/sample.csv", false},
		{"http: uri not exist", "https://raw.githubusercontent.com/tacklehop/csvsearch/main/non-exist.csv", true},
		{"http: non csv format", "https://raw.githubusercontent.com/tacklehop/csvsearch/main/error.csv", true},
		{"http: domain not exist", "https://not-exist.com", true},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// disabled because test results are shown in parallel.
			// t.Parallel()
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
		})
	}
}

func TestSearchWords(t *testing.T) {
	testfile := "sample.csv"
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"case 1", "Key1", "[Key1 Word1  ]"},
		{"case 2", "Word3", "[Key2 Word1 Word2 Word3][Key3 Word3 Word1]"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
		})
	}

}
