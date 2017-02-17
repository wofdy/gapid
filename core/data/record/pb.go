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
	"encoding/binary"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/google/gapid/core/context/jot"
	"github.com/google/gapid/core/event"
	"github.com/google/gapid/core/fault/cause"
	"github.com/google/gapid/core/log"
)

// pbbHandler is an implementation of fileType that stores it's records in binary proto format.
type pbbFileType struct{}

type pbbHandler struct {
	f    *os.File
	null proto.Message
}

type pbbReader struct {
	buf  []byte
	f    io.Reader
	null proto.Message
}

func (pbbFileType) Ext() string { return ".pb" }

func (pbbFileType) Open(ctx log.Context, f *os.File, null interface{}) (LedgerInstance, error) {
	m, ok := null.(proto.Message)
	if !ok {
		return nil, cause.Explain(ctx, nil, "Cannot create proto ledger with non proto type")
	}
	return &pbbHandler{f: f, null: m}, nil
}

func (h *pbbHandler) Write(ctx log.Context, record interface{}) error {
	buf, err := proto.Marshal(record.(proto.Message))
	if err != nil {
		return err
	}
	size := int32(len(buf))
	if err := binary.Write(h.f, binary.LittleEndian, &size); err != nil {
		return err
	}
	_, err = h.f.Write(buf)
	return err
}

func (h *pbbHandler) Reader(ctx log.Context) event.Source {
	return &pbbReader{f: &readAt{f: h.f}, null: h.null}
}

func (h *pbbHandler) Close(ctx log.Context) {
	h.f.Close()
}

func (h *pbbHandler) New(ctx log.Context) interface{} {
	return proto.Clone(h.null)
}

func (r *pbbReader) Next(ctx log.Context) interface{} {
	size := int32(0)
	if err := binary.Read(r.f, binary.LittleEndian, &size); err != nil {
		if err != io.EOF {
			jot.Fail(ctx, err, "Invalid proto record header in ledger")
		}
		return nil
	}
	if cap(r.buf) < int(size) {
		r.buf = make([]byte, size*2) // TODO: very naive growth algorithm
	}
	r.buf = r.buf[0:size]
	io.ReadFull(r.f, r.buf)
	message := proto.Clone(r.null)
	err := proto.Unmarshal(r.buf, message)
	if err != nil {
		jot.Fail(ctx, err, "Invalid proto in ledger")
		return nil
	}
	return message
}

func (h *pbbReader) Close(ctx log.Context) {}
