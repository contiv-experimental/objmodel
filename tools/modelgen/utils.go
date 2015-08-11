/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*********************** Helper funcs *************************/
var (
	newlines  = regexp.MustCompile(`(?m:\s*$)`)
	acronyms  = regexp.MustCompile(`(Url|Http|Id|Io|Uuid|Api|Uri|Ssl|Cname|Oauth|Otp)$`)
	camelcase = regexp.MustCompile(`(?m)[-.$/:_{}\s]`)
)

func initialCap(ident string) string {
	if ident == "" {
		panic("blank identifier")
	}
	return depunct(ident, true)
}

func initialLow(ident string) string {
	if ident == "" {
		panic("blank identifier")
	}
	return depunct(ident, false)
}

func depunct(ident string, initialCap bool) string {
	matches := camelcase.Split(ident, -1)
	for i, m := range matches {
		if initialCap || i > 0 {
			m = capFirst(m)
		}
		matches[i] = acronyms.ReplaceAllStringFunc(m, func(c string) string {
			if len(c) > 4 {
				return strings.ToUpper(c[:2]) + c[2:]
			}
			return strings.ToUpper(c)
		})
	}
	return strings.Join(matches, "")
}

func capFirst(ident string) string {
	r, n := utf8.DecodeRuneInString(ident)
	return string(unicode.ToUpper(r)) + ident[n:]
}
