## SPIRV-Tools Generated Headers

This directory contains the generated headers and include files used by SPIRV-Tools.

The headers can be re-generated by following these steps:

 1. `mkdir <tmp folder> && cd <tmp folder>`
 1. `git clone https://github.com/KhronosGroup/SPIRV-Tools`
 2. `cd SPIRV-Tools`
 3. `git checkout <sha from WORKSPACE>`
 4. `git clone https://github.com/KhronosGroup/SPIRV-Headers.git external/spirv-headers`
 5. `git clone https://github.com/google/googletest.git external/googletest`
 6. `mkdir build && cd build`
 7. `cmake -G 'Unix Makefiles' ..`
 8. `make -j <N>`
 9. `cp *.inc *.h <gapid>/tools/build/third_party/SPIRV-Tools-Generated`
