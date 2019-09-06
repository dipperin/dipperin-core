# Install Compiler

Dipc is a compiler that compile Dipperin C++ smart contract code to WebAssembly program.

## Build From Source Code

### Required

- GCC 5.4+ or Clang 4.0+
- CMake 3.5+
- Git
- Python

### Ubuntu 

**Required:** 16.04+

- **Install Dependencies**

``` shell
sudo apt install build-essential cmake libz-dev libtinfo-dev
```

- **Get Source Code**

```shell
git clone https://github.com/dipperin/dipc.git
cd dipc
git submodule update --init --recursive
```
- **Build Code**

``` sh
cd dipc
mkdir build && cd build
cmake .. 
make && make install
```

### Windows

**Required:** [MinGW-W64 GCC-8.1.0](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z)

**NOTES:** _MinGW and CMake must be installed in a directory without space._

- **Get Source Code**

```shell
git clone https://github.com/dipperin/dipc.git
cd dipc
git submodule update --init --recursive
```
- **Build Code**

``` sh
cd dipc
mkdir build && cd build
cmake -G "MinGW Makefiles" .. -DCMAKE_INSTALL_PREFIX="C:/dipc.cdt" -DCMAKE_MAKE_PROGRAM=mingw32-make
mingw32-make && mingw32-make install
```

## Use Dipc

### Skeleton Smart Contract Without CMake Support

- Init a project

``` sh
dipc-init -project example -bare
```

- Build contract

``` sh
cd example
dipc-cpp -o example.wasm example.cpp -abigen
```

### Skeleton Smart Contract With CMake Support

- Init CMake project

``` sh
dipc-init -project cmake_example
```

- Build contract
  * Linux
  ```
  cd cmake_example/build
  cmake ..
  ```
  * Windows
  >**Required:**
  >+ [MinGW-W64 GCC-8.1.0](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z)
  >+ CMake 3.5 or higher

  ```sh
  cd cmake_example/build
  cmake .. -G "MinGW Makefiles" -DCMAKE_PREFIX_PATH=<cdt_install_dir>
  ```