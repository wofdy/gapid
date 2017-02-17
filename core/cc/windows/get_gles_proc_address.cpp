/*
 * Copyright (C) 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "../dl_loader.h"
#include "../get_gles_proc_address.h"
#include "../log.h"
#include "../target.h" // STDCALL

#include <string>
#include <sstream>

#include <windows.h>
#include <wingdi.h>

namespace {

std::string systemOpengl32Path() {
    char sysdir[MAX_PATH];
    GetSystemDirectoryA(sysdir, MAX_PATH-1);

    std::stringstream path;
    path << sysdir << "\\opengl32.dll";
    return path.str();
}

void* getGlesProcAddress(const char* name, bool bypassLocal) {
    using namespace core;
    typedef void* (*GPAPROC)(const char *name);

    static DlLoader opengl(bypassLocal ? systemOpengl32Path().c_str() : "opengl32.dll");
    if (GPAPROC gpa = reinterpret_cast<GPAPROC>(opengl.lookup("wglGetProcAddress"))) {
        if (void* proc = gpa(name)) {
            GAPID_INFO("GetGlesProcAddress(%s, %d) -> 0x%x (via opengl32 wglGetProcAddress)", name, bypassLocal, proc);
            return proc;
        }
    }
    if (void* proc = opengl.lookup(name)) {
        GAPID_INFO("GetGlesProcAddress(%s, %d) -> 0x%x (from opengl32 symbol)", name, bypassLocal, proc);
        return proc;
    }

    GAPID_INFO("GetGlesProcAddress(%s, %d) -> not found", name, bypassLocal);
    return nullptr;
}

}  // anonymous namespace

namespace core {

GetGlesProcAddressFunc* GetGlesProcAddress = getGlesProcAddress;

}  // namespace core
