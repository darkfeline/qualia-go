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
	"bytes"
	"strings"
	"testing"
)

func TestQualify(t *testing.T) {
	t.Parallel()
	cases := []struct {
		qs      []string
		in, out string
	}{
		// No quality idempotence
		{
			[]string{},
			"# BEGIN spam\n#spam\n# END spam\n",
			"# BEGIN spam\n#spam\n# END spam\n",
		},
		// With quality idempotence
		{
			[]string{"spam"},
			"# BEGIN spam\nspam\n# END spam\n",
			"# BEGIN spam\nspam\n# END spam\n",
		},
		// With quality
		{
			[]string{"spam"},
			"# BEGIN spam\n#spam\n# END spam\n",
			"# BEGIN spam\nspam\n# END spam\n",
		},
		// No quality
		{
			[]string{},
			"# BEGIN spam\nspam\n# END spam\n",
			"# BEGIN spam\n#spam\n# END spam\n",
		},
		// Ignore unclosed
		{
			[]string{"spam", "eggs"},
			"# BEGIN spam\n#spam\n# BEGIN eggs\n#eggs\n# END eggs\n",
			"# BEGIN spam\n#spam\n# BEGIN eggs\n#eggs\n# END eggs\n",
		},
		// Ignore unqualified
		{[]string{}, "spam\n", "spam\n"},
		// No whitespace before keyword
		{
			[]string{"spam"},
			"#BEGIN spam\n#spam\n#END spam\n",
			"#BEGIN spam\nspam\n#END spam\n"},
	}
	for _, c := range cases {
		r := strings.NewReader(c.in)
		b := new(bytes.Buffer)
		err := Qualify(c.qs, r, b)
		got := b.String()
		if got != c.out {
			t.Errorf("For %#v %#v, expected %#v, got %#v with error %#v", c.qs, c.in, c.out, got, err)
		}
	}

}
