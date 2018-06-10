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
	"reflect"
	"testing"
)

func TestCommenter_comment(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, out []string
	}{
		{[]string{}, []string{}},
		{[]string{"foo", "bar"}, []string{"#foo", "#bar"}},
		// Idempotent
		{[]string{"#foo", "#bar"}, []string{"#foo", "#bar"}},
		// Preserve different indent
		{[]string{"foo", " bar"}, []string{"#foo", "# bar"}},
		// Skip common indent
		{[]string{" foo", "  bar"}, []string{" #foo", " # bar"}},
	}
	cm := newCommenter("#")
	for _, c := range cases {
		got := cm.comment(c.in)
		if !reflect.DeepEqual(got, c.out) {
			t.Errorf("For %#v, expected %#v, got %#v", c.in, c.out, got)
		}
	}
}

func TestCommenter_uncomment(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, out []string
	}{
		{[]string{}, []string{}},
		{[]string{"#foo", "#bar"}, []string{"foo", "bar"}},
		// Idempotent
		{[]string{"foo", "bar"}, []string{"foo", "bar"}},
		// Preserve partial
		{[]string{"foo", "#bar"}, []string{"foo", "#bar"}},
		// Uncomment many
		{[]string{"##foo", "###bar"}, []string{"foo", "#bar"}},
		// With indentation before
		{[]string{" #foo", "  #bar"}, []string{" foo", "  bar"}},
		// With indentation after
		{[]string{"#foo", "# bar"}, []string{"foo", " bar"}},
		// With indentation before and after
		{[]string{"#foo", " # bar"}, []string{"foo", "  bar"}},
	}
	cm := newCommenter("#")
	for _, c := range cases {
		got := cm.uncomment(c.in)
		if !reflect.DeepEqual(got, c.out) {
			t.Errorf("For %#v, expected %#v, got %#v", c.in, c.out, got)
		}
	}
}
