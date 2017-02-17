// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import "github.com/google/gapid/core/log"

// Get returns a versioned file using `git show` at the specified CL.
func (g Git) Get(ctx log.Context, path string, at SHA) ([]byte, error) {
	str, _, err := g.run(ctx, "show", at.String()+":"+path)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
