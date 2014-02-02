Release script for nord
=======================

This is an adapted script from websocketd project by Joe Walnes. Thanks!

Perform a fully automated repeatable release of websocketd. Including:

*   Downloading the correct version of Go
*   Cross-compiling Go for all required platforms
*   Inferring the next nord release version
*   Cross-compiling websocketd for all required platforms

The actual Go release will be downloaded and cross-compiled. All compilation
occurs in an isolated environment to ensure there are no conflicts with any
system install Go. The first build will be really slow (cross compiling Go core),
but subsequent builds are incremental.

To build the packages for all platforms:

```bash
make build
```

Now all compiled binaries can be found in out/$RELEASE_VERSION directory

To create zip archives for uploading on Github:

```bash
cd out/$RELEASE_VERSION
for i in *; do zip -r "$i" "$i" -x *.built* -x *.released*; done
```

Now new release with binaries can be created manually on Github

To clean up:

```bash
make clean
```