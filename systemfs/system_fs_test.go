// Copyright (c) 2019-2020 Siemens AG
// Copyright (c) 2019-2021 Jonas Plum
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package systemfs

import (
	"io"
	"io/fs"
	"os"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/forensicanalysis/fslib"
)

func TestFS(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	fsys, err := New()
	if err != nil {
		t.Fatal(err)
	}

	wd, err = fslib.ToFSPath(wd)
	if err != nil {
		t.Fatal(err)
	}
	fsys, err = fs.Sub(fsys, wd)
	if err != nil {
		t.Fatal(err)
	}

	if err := fstest.TestFS(fsys, "systemfs.go"); err != nil {
		t.Fatal(err)
	}
}

func Test_LocalNTFS(t *testing.T) {
	if runtime.GOOS == "windows" {
		_, err := os.OpenFile(`\\.\C:`, os.O_RDONLY, fs.FileMode(0666))
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		path   string
		header string
	}{
		{"C/$MFT", "FILE"},
		{"C/Windows/System32/config/SOFTWARE", "regf"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			if runtime.GOOS == "windows" {
				fsys, err := New()
				if err != nil {
					t.Fatalf("Error %s", err)
				}

				file, err := fsys.Open(test.path)
				if err != nil {
					t.Fatalf("Error %s", err)
				}

				info, err := file.Stat()
				if err != nil {
					t.Fatalf("Error %s", err)
				}
				if info.Size() == 0 {
					t.Errorf("file is 0 byte")
				}

				header := make([]byte, len(test.header))
				n, err := file.Read(header)
				if err != nil && err != io.EOF {
					t.Errorf("read error %s", err)
				}
				if n != len(test.header) {
					t.Errorf("Wrong read count got: %d, want: %d", n, len(test.header))
				}
				if string(header) != test.header {
					t.Errorf("Wrong header got: %s, want: %s", header, test.header)
				}
			}
		})
	}
}
