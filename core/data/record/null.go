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

package record

import (
	"github.com/google/gapid/core/event"
	"github.com/google/gapid/core/log"
)

// nullShelf is an implementation of Shelf that creates null ledgers.
type nullShelf struct{}

// nullLedger is an implementation of Ledger that just ignores all record append requests.
type nullLedger struct{}

func NewNullShelf(ctx log.Context) (Shelf, error)                         { return &nullShelf{}, nil }
func (nullShelf) Open(log.Context, string, interface{}) (Ledger, error)   { return nullLedger{}, nil }
func (nullShelf) Create(log.Context, string, interface{}) (Ledger, error) { return nullLedger{}, nil }
func (nullLedger) Read(ctx log.Context, h event.Handler) error            { return nil }
func (nullLedger) Watch(ctx log.Context, w event.Handler)                 {}
func (nullLedger) Add(ctx log.Context, record interface{}) error          { return nil }
func (nullLedger) Close(log.Context)                                      {}
func (nullLedger) New(log.Context) interface{}                            { return nil }
