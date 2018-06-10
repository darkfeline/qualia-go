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

// Command qualia conditionally comments and uncomments blocks in
// files, for example configuration files (dotfiles).  This can be
// used to keep dotfiles for different machines in a single version
// control repository and check out the right copy on each machine.
//
// qualia recognizes special blocks (called qualified blocks) and
// comments or uncomments them.  A qualified block looks like the
// following:
//
//   # BEGIN laptop
//   export PATH="$HOME/bin:$PATH"
//   # END laptop
//
// The quality of this block is laptop.  If laptop is given as a
// quality, then qualia will make sure the contents of the block are
// uncommented.  If laptop isn't given as a quality, then qualia will
// make sure the contents of the block are commented.
//
// It is possible to pass multiple qualities or no qualities:
//
//   qualia audio games
//   qualia
//
// qualia is idempotent, so you can run it multiple times; only the
// last time takes effect:
//
//   qualia <infile | qualia laptop | qualia desktop | qualia laptop
package main

import (
	"fmt"
	"os"

	"go.felesatra.moe/qualia"
)

func main() {
	qs := os.Args[1:]
	if err := qualia.Qualify(qs, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
