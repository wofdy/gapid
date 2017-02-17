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

#include "connection_header.h"

#include "core/cc/log.h"
#include "core/cc/stream_reader.h"

namespace gapii {

ConnectionHeader::ConnectionHeader()
    : mVersion(0)
    , mObserveFrameFrequency(0)
    , mObserveDrawFrequency(0)
    , mFlags(0) {}

bool ConnectionHeader::read(core::StreamReader* reader) {
    if (!reader->read(mMagic)) {
        return false;
    }
    if (mMagic[0] != 's'
     || mMagic[1] != 'p'
     || mMagic[2] != 'y'
     || mMagic[3] != '0') {
        GAPID_WARNING("ConnectionHeader magic was not as expected. Got %c%c%c%c",
            mMagic[0], mMagic[1], mMagic[2], mMagic[3]);
        return false;
     }

    // TODO: Endian-swap data if GAPII is running on a big-endian architecture.

    if (!reader->read(mVersion)) {
        return false;
    }

    const int kMinSupportedVersion = 2;
    const int kMaxSupportedVersion = 3;

    if (mVersion < kMinSupportedVersion || mVersion > kMaxSupportedVersion) {
        GAPID_WARNING("Unsupported ConnectionHeader version %d. Only understand [%d to %d].",
                mVersion, kMinSupportedVersion, kMaxSupportedVersion);
        return false;
    }
    if (mVersion >= 2) {
        if (!reader->read(mObserveFrameFrequency) ||
            !reader->read(mObserveDrawFrequency)) {
            return false;
        }
    }
    if (mVersion >= 3) {
        if (!reader->read(mFlags)) {
            return false;
        }
    }

    // Insert new version handling here. Don't forget to bump kMaxSupportedVersion!
    return true;
}

} // namespace gapii
