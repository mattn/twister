// Copyright 2010 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package web

import (
	"testing"
	"reflect"
	"bufio"
	"bytes"
)

var quoteHeaderValueTests = []struct {
	s            string
	quote        string
	quoteOrToken string
}{
	{s: "a", quote: "\"a\"", quoteOrToken: "a"},
	{s: "x\"y", quote: "\"x\\\"y\"", quoteOrToken: "\"x\\\"y\""},
	{s: "x\\y", quote: "\"x\\\\y\"", quoteOrToken: "\"x\\\\y\""},
}

func TestQuoteHeaderValue(t *testing.T) {
	for _, tt := range quoteHeaderValueTests {
		if quote := QuoteHeaderValue(tt.s); quote != tt.quote {
			t.Errorf("QuoteHeaderValue(%q) = %q, want %q", tt.s, quote, tt.quote)
		}
		if quoteOrToken := QuoteHeaderValueOrToken(tt.s); quoteOrToken != tt.quoteOrToken {
			t.Errorf("QuoteHeaderValueOrToken(%q) = %q, want %q", tt.s, quoteOrToken, tt.quoteOrToken)
		}
	}
}

var UnquoteHeaderValueTests = []struct {
	s       string
	unquote string
}{
	{s: "a", unquote: "a"},
	{s: "a b", unquote: "a b"},
	{s: "\"a\"", unquote: "a"},
	{s: "\"a \\\\ b \\\" c\"", unquote: "a \\ b \" c"},
}

func TestUnquoteHeaderValue(t *testing.T) {
	for _, tt := range UnquoteHeaderValueTests {
		if unquote := UnquoteHeaderValue(tt.s); unquote != tt.unquote {
			t.Errorf("UnquoteHeaderValue(%q) = %q, want %q", tt.s, unquote, tt.unquote)
		}
	}
}

var getHeaderListTests = []struct {
	s string
	l []string
}{
	{s: "a", l: []string{"a"}},
	{s: "a, b , c ", l: []string{"a", "b", "c"}},
	{s: "a,b,c", l: []string{"a", "b", "c"}},
	{s: " a b, c d ", l: []string{"a b", "c d"}},
	{s: "\"a, b, c\", d ", l: []string{"\"a, b, c\"", "d"}},
	{s: "\" \"", l: []string{"\" \""}},
}

func TestGetHeaderList(t *testing.T) {
	for _, tt := range getHeaderListTests {
		header := NewHeaderMap("foo", tt.s)
		if l := header.GetList("foo"); !reflect.DeepEqual(tt.l, l) {
			t.Errorf("GetList for %q = %q, want %q", tt.s, l, tt.l)
		}
	}
}

var parseHTTPHeaderTests = []struct {
	name   string
	header HeaderMap
	s      string
}{
	{"multihdr", NewHeaderMap(
		HeaderContentType, "text/html",
		HeaderCookie, "hello=world",
		HeaderCookie, "foo=bar"),
`Content-Type: text/html
CoOkie: hello=world
Cookie: foo=bar

`},
	{"continuation", NewHeaderMap(
		HeaderContentType, "text/html",
		HeaderCookie, "hello=world, foo=bar"),
`Cookie: hello=world,
 foo=bar
Content-Type: text/html

`},
}

func TestParseHttpHeader(t *testing.T) {
	for _, tt := range parseHTTPHeaderTests {
		b := bufio.NewReader(bytes.NewBufferString(tt.s))
		header := HeaderMap{}
		err := header.ParseHttpHeader(b)
		if err != nil {
			t.Errorf("ParseHttpHeader error for %s = %v", tt.name, err)
		}
		if !reflect.DeepEqual(tt.header, header) {
			t.Errorf("ParseHttpHeader for %s = %q, want %q", tt.name, header, tt.header)
		}
	}
}

func TestParseHttpHeaderBytes(t *testing.T) {
	for _, tt := range parseHTTPHeaderTests {
		header := HeaderMap{}
		n, err := header.ParseHttpHeaderBytes([]byte(tt.s))
		if err != nil {
			t.Errorf("ParseHttpHeaderBytes error for %s = %v", tt.name, err)
		}
		if !reflect.DeepEqual(tt.header, header) {
			t.Errorf("ParseHttpHeaderBytes for %s = %q, want %q", tt.name, header, tt.header)
		}
		if n != len(tt.s) {
			t.Errorf("ParseHEaderBytes returned n = %d, want %d", n, len(tt.s))
		}
	}
}