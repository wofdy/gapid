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

package app

import (
	"github.com/google/gapid/core/fault/severity"
	"github.com/google/gapid/core/text/note"
)

type (
	AppFlags struct {
		Version bool `help:"_display the application version"`
		Log     LogFlags
		Profile ProfileFlags
	}
	LogFlags struct {
		Level  severity.Level `help:"_The severity to enable logs at"`
		Style  note.Style     `help:"_The style to use when printing the log"`
		Stacks bool           `help:"_If true, stack traces are logged for all errors"`
		File   string         `help:"_The file to store the logs in"`
	}
	ProfileFlags struct {
		CPU string `help:"_write cpu profile to file"`
	}
)
