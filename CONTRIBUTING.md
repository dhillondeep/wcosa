# Development Setup

Clone the repository into your `GOPATH` as `$(GOPATH)/src/wio` with `--recursive` or `git submodule update --init`
to retrieve some of the CMake toolchain modules.

Run `make get` to retrieve development dependencies and make sure to install your OS-specific go library.
For linux users `wmake linux-setup` should handle most things. You will need to symlink `bin/toolchain` to
the `toolchain` folder in order to use `wio build` and `wio run`.

The repo includes some scripts like `wmake` that shorthand some common development commands. And running
`source wenv` will allow the local `wio` executable to take precedence over an installed version. It also
adds `wmake` to your path for easier development

### Installing Go

Wio development should be done on Go version 1.10.2. On linux this can be installed with

```bash
curl -o golang.tar https://storage.googleapis.com/golang/go1.10.2.linux-amd64.tar.gz
tar -xzvf golang.tar
```
