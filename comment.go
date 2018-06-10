// Copyright (C) 2018 Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qualia

import (
	"fmt"
	"regexp"
	"strings"
)

type commenter struct {
	s   string
	pat *regexp.Regexp
}

func newCommenter(c string) commenter {
	return commenter{
		s:   c,
		pat: regexp.MustCompile(fmt.Sprintf(`^(\s*)%s`, regexp.QuoteMeta(c))),
	}
}

// isCommented returns true if all lines are commented.
func (c commenter) isCommented(lines []string) bool {
	for _, l := range lines {
		if !c.pat.MatchString(l) {
			return false
		}
	}
	return true
}

func (c commenter) comment(lines []string) []string {
	if len(lines) == 0 || c.isCommented(lines) {
		return lines
	}
	in := commonIndent(lines)
	sb := strings.Builder{}
	n := make([]string, len(lines))
	for i, l := range lines {
		sb.Reset()
		sb.WriteString(in)
		sb.WriteString(c.s)
		sb.WriteString(l[len(in):])
		n[i] = sb.String()
	}
	return n
}

func (c commenter) uncomment(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	n := make([]string, len(lines))
	copy(n, lines)
	for c.isCommented(n) {
		for i, l := range n {
			n[i] = c.pat.ReplaceAllString(l, "$1")
		}
	}
	return n
}
