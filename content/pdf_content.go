// Copyright (c) 2019 Siemens AG
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

package content

import (
	"bytes"
	"github.com/forensicanalysis/fslib/fsio"

	"github.com/ledongthuc/pdf"
)

// PDFContent returns the text data from a pdf file.
func PDFContent(r fsio.ReadSeekerAt) (string, error) {
	size, err := fsio.GetSize(r)
	if err != nil {
		return "", err
	}

	file, err := pdf.NewReader(r, size)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := file.GetPlainText()
	if err != nil {
		return "", err
	}
	if _, err = buf.ReadFrom(b); err != nil {
		return "", err
	}
	return buf.String(), nil
}