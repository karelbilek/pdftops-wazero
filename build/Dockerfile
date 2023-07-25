FROM emscripten/emsdk:3.1.44

# remove prebuilt cache
WORKDIR /emsdk/upstream/emscripten
RUN rm -Rf cache

# check in https://github.com/emscripten-core/emscripten/pull/18285
COPY patched_emscripten/system/lib/standalone/standalone.c /emsdk/upstream/emscripten/system/lib/standalone/standalone.c
COPY patched_emscripten/system/lib/libc/musl/src/stdio/freopen.c /emsdk/upstream/emscripten/system/lib/libc/musl/src/stdio/freopen.c

# copy xpdf code from mirror
COPY xpdfreader-mirror /xpdf

# get pre-built emscripten ports
WORKDIR /xpdf
RUN mkdir -p build
WORKDIR /xpdf/build
RUN embuilder build zlib
RUN embuilder build libpng
RUN embuilder build libjpeg
RUN embuilder build freetype

RUN apt-get install -y make gcc

# make the cmake target
RUN emcmake cmake .. -DCMAKE_TOOLCHAIN_FILE=$EMSDK/upstream/emscripten/cmake/Modules/Platform/Emscripten.cmake -DCMAKE_CXX_FLAGS="-std=c++14 -O2 -DLOAD_FONTS_FROM_MEM=1" -DCMAKE_EXE_LINKER_FLAGS="-static -sERROR_ON_UNDEFINED_SYMBOLS=0 -s WASM=1 -s ALLOW_MEMORY_GROWTH=1 -s STANDALONE_WASM=1 -s USE_FREETYPE=1 -s USE_ZLIB=1 -s USE_LIBPNG=1 -s USE_LIBJPEG=1"

# actual make
RUN nproc | xargs -I % emmake make -j%