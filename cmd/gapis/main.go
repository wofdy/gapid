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
	"path/filepath"
	"time"

	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/app/auth"
	"github.com/google/gapid/core/context/jot"
	"github.com/google/gapid/core/event/task"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/android/adb"
	"github.com/google/gapid/core/os/device"
	"github.com/google/gapid/core/os/device/bind"
	"github.com/google/gapid/gapis/database"
	"github.com/google/gapid/gapis/replay"
	"github.com/google/gapid/gapis/server"
	"github.com/google/gapid/gapis/service"
	"github.com/google/gapid/gapis/stringtable"
)

var (
	rpc             = flag.String("rpc", "localhost:0", "TCP host:port of the server's RPC listener")
	stringsPath     = flag.String("strings", "strings", "Directory containing string table packages")
	persist         = flag.Bool("persist", false, "Server will keep running even when no connections remain")
	gapisAuthToken  = flag.String("gapis-auth-token", "", "The connection authorization token for gapis")
	gapirAuthToken  = flag.String("gapir-auth-token", "", "The connection authorization token for gapir")
	gapirArgStr     = flag.String("gapir-args", "", `"<The arguments to be passed to gapir>"`)
	scanAndroidDevs = flag.Bool("monitor-android-devices", true, "Server will scan for locally connected Android devices")
	addLocalDevice  = flag.Bool("add-local-device", true, "Server will create a new local replay device")
)

func main() {
	app.ShortHelp = "Gapis is the graphics API server"
	app.Name = "GAPIS" // Has to be this for version parsing compatability
	app.Version = version
	app.Run(run)
}

// features is the reported list of features supported by the server.
// This feature list can be used by the client to determine what new RPCs can be
// called.
var features = []string{}

func run(ctx log.Context) error {
	m, r := replay.New(ctx), bind.NewRegistry()
	ctx = replay.PutManager(ctx, m)
	ctx = bind.PutRegistry(ctx, r)
	ctx = database.Put(ctx, database.NewInMemory(ctx))

	deviceScanDone, onDeviceScanDone := task.NewSignal()
	if *scanAndroidDevs {
		go monitorAndroidDevices(ctx, r, onDeviceScanDone)
	} else {
		onDeviceScanDone(ctx)
	}

	if *addLocalDevice {
		r.AddDevice(ctx, bind.Host(ctx))
	}

	return server.Listen(ctx, *rpc, server.Config{
		Info: &service.ServerInfo{
			Name:         device.Host(ctx).Name,
			VersionMajor: uint32(version.Major),
			VersionMinor: uint32(version.Minor),
			Features:     features,
		},
		StringTables:   loadStrings(ctx),
		AuthToken:      auth.Token(*gapisAuthToken),
		DeviceScanDone: deviceScanDone,
	})
}

func monitorAndroidDevices(ctx log.Context, r *bind.Registry, onDeviceScanDone task.Task) {
	// Populate the registry with all the existing devices.
	func() {
		defer onDeviceScanDone(ctx) // Signal that we have a primed registry.

		if devs, err := adb.Devices(ctx); err == nil {
			for _, d := range devs {
				r.AddDevice(ctx, d)
			}
		}
	}()

	if err := adb.Monitor(ctx, r, time.Second*3); err != nil {
		jot.Warning(ctx).Cause(err).Print("Could not scan for local Android devices")
	}
}

func loadStrings(ctx log.Context) []*stringtable.StringTable {
	files, err := filepath.Glob(filepath.Join(*stringsPath, "*.stb"))
	if err != nil {
		jot.Fail(ctx, err, "Couldn't scan for stringtables")
		return nil
	}

	out := make([]*stringtable.StringTable, 0, len(files))

	for _, path := range files {
		ctx := ctx.S("path", path)
		st, err := stringtable.Load(path)
		if err != nil {
			jot.Fail(ctx, err, "Couldn't load stringtable file")
			continue
		}
		out = append(out, st)
	}

	return out
}
