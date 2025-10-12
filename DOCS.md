# Pacakge Docs

Rendering HTML docs of the `hit` package locally is possible with installing the `pkgsite` package and running it
from the working folder.

```fish
set GOPATH $(go env GOPATH)
go install golang.org/x/pkgsite/cmd/pkgsite@latest
$GOPATH/bin/pkgsite .
```