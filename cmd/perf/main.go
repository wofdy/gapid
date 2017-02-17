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

package main

import (
	"flag"

	"github.com/google/gapid/core/app"
	_ "github.com/google/gapid/framework/binary/any"
)

var (
	flagVerifyHashes bool
)

func main() {
	app.ShortHelp = "Perf is a performance regression testing tool"
	app.Version = app.VersionSpec{Major: 0, Minor: 1}
	flag.BoolVar(&flagVerifyHashes, "v", false, "warn if any external files have changed")
	app.Run(app.VerbMain)
}
