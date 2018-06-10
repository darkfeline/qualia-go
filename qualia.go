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

// Package qualia implements the primary functionality of the qualia
// command.
//
// qualia recognizes special blocks (called qualified blocks) and
// comments or uncomments them. A qualified block looks like the
// following:
//
//   # BEGIN laptop
//   export PATH="$HOME/bin:$PATH"
//   # END laptop
//
// The quality of this block is laptop. If laptop is given as a
// quality, then qualia will make sure the contents of the block are
// uncommented. If laptop isn't given as a quality, then qualia will
// make sure the contents of the block are commented.
package qualia

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// Qualify qualifies the lines read from the Reader with the qualities
// and writes the lines to the Writer.
func Qualify(q []string, r io.Reader, w io.Writer) error {
	qm := make(map[string]bool)
	for _, q := range q {
		qm[q] = true
	}
	qf := qualifier{
		qualities: qm,
	}
	return qf.qualify(r, w)
}

type qualifier struct {
	qualities map[string]bool

	block      []string
	blockAttrs blockAttrs
	inBlock    bool
}

func (q *qualifier) qualify(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		for _, l := range q.feedLine(s.Text()) {
			_, err := fmt.Fprintln(w, l)
			if err != nil {
				return err
			}
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	if q.inBlock {
		for _, l := range q.block {
			_, err := fmt.Fprintln(w, l)
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("unclosed block %s", q.blockAttrs.quality)
	}
	return nil
}

func (q *qualifier) feedLine(l string) []string {
	if !q.inBlock {
		if b, ok := findBlockStart(l); ok {
			q.inBlock = true
			q.blockAttrs = b
			q.block = make([]string, 0, 8)
		}
		return []string{l}
	}
	if !q.blockAttrs.findBlockEnd(l) {
		q.block = append(q.block, l)
		return nil
	}
	c := newCommenter(q.blockAttrs.prefix)
	var n []string
	if q.qualities[q.blockAttrs.quality] {
		n = c.uncomment(q.block)
	} else {
		n = c.comment(q.block)
	}
	q.inBlock = false
	return append(n, l)
}

var blockStart = regexp.MustCompile(`^\s*(\S+)\s*BEGIN\s+(\S+)`)
var blockEnd = regexp.MustCompile(`^\s*(\S+)\s*END\s+(\S+)`)

type blockAttrs struct {
	prefix, quality string
}

func findBlockStart(l string) (b blockAttrs, ok bool) {
	r := blockStart.FindStringSubmatch(l)
	if r == nil {
		return blockAttrs{}, false
	}
	b = blockAttrs{
		prefix:  r[1],
		quality: r[2],
	}
	return b, true
}

func (b blockAttrs) findBlockEnd(l string) bool {
	r := blockEnd.FindStringSubmatch(l)
	if r == nil {
		return false
	}
	return r[1] == b.prefix && r[2] == b.quality
}
