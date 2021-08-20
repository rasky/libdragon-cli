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

### Installation (binary)

Binary releases coming soon (stay tuned). For now, see below installation from
source.


### Quick sheet

 1. Start from an empty git repository (`git init`).
 1. Run `libdragon init` to create the skeleton project.
 1. Run `libdragon make` to build the first ROM.
 1. Hack away and enjoy!


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

