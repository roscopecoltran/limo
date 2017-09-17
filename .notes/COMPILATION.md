https://stackoverflow.com/questions/24601619/how-to-improve-golang-compilation-speed

```
Try go install -a github.com/mattn/go-sqlite3 which will install the compiled-against-Go-1.3 package into your $GOPATH.

Right now, you likely have an older version installed under $GOPATH/pkg/ and therefore Go is recompiling it for every build.
```