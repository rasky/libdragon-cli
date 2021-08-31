## libdragon-cli: a tool for managing libdragon projects

Main features:

 * Create a libdragon ROM skeleton, with a fully-featured Makefile-based
   build system
 * Vendor the libdragon library in the repository, and compile it as part
   of the build, so it's easy to do local modifications when required.
   Vendoring can be done with either git submodule or git subtree.
 * Download and use a complete toolchain as a Docker repository, so that
   it doesn't pollute your PATH. You can have different toolchain versions
   for different projects, if required so.
 * Update libdragon and toolchain.
 * Available on Windows, Mac and Linux.

### Installation (binary)

Download the [latest binary release](https://github.com/rasky/libdragon-cli/releases/latest) from GitHub.
`libdragon-cli` ships as a single binary available on Windows, Mac and Linux.
Simply decompress the binary and put it in a directory in your PATH.

Otherwise, see below for instructions to install from source.


### Quick sheet

 1. Start from an empty git repository (`git init`).
 1. Run `libdragon init` to create the skeleton project, with a vendored copy
    of lidbragon, and download an updated Docker toolchain image.
 1. Run `libdragon make` to build the first ROM. This will auto-start the docker
    image and run `make` within it.
 1. Hack away and enjoy!

### Other commands:

 * `libdragon disasm`: show disassembly of the current project, You can pass 
   a symbol as argument to request disassembly of a single function 
   (eg: `libdragon disasm dfs_read`).
 * `libdragon exec`: run a command within the Docker container. This can be
   useful to manually execute libdragon tools. For instance: 
   `libdragon exec makedfs <arguments>`
 * `libdragon start` and `libdragon stop` help explicitly managing the
   docker instance associated to the current git repository. In general,
   `libdragon` will create one container per repository.
 * `libdragon update` will update both the vendored copy of libdragon
   (lastest version on Github) and the lastest toolchain (from Docker Hub).
   You can update only either of the two with specific options (see the help).
 * `libdragon init` can vendor libdragon with `git subtree` (default) or
   with `git submodule`. If you prefer the latter, use `libdragon init --submodule`.


### FAQ

  * Can I use this tool on an existing libdragon project?

Sure. Use `libdragon update toolchain` to download the toolchain, and then
`libdragon make` for your development cycle.

  * Can I use this tool for a project that doesn't use git?

Sure, but you will not be able to manage libdragon vendoring wihtout git.
Use `libdragon update toolchain` to download the toolchain, and then
`libdragon make` for your development cycle. You won't be able to manage
libdragon vendoring through the tool, though.

  * How can I use a different toolchain?

libdragon-cli uses the official libdragon Docker toolchain. If you want to
switch to your own toolchain, use `libdragon update toolchain --image user/image:tag`
to specify a Docker image in Docker Hub format.


### Building libdragon-cli from source

You need Go 1.16 or later to build libdragon-cli. Once you have Go installed
in your system, simply run:

	$ git clone https://github.com/rasky/libdragon-cli
	$ cd libdragon-cli
	$ go build -o libdragon

This will build a `libdragon` binary. Move it to a directory in your `$PATH`
and you're done, there are no additional dependencies.

