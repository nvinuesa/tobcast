# Total order broadcast (tobcast)

Please visit https://underscorenico.github.io/blog/2019/11/12/total-order-broadcast/ for implementation details.

### Build
The project uses makefile, to build simply run: 

```bash
make all
```

and to test it simply run:

```bash
make test
```
or 

```bash
make test-race
```

### Run

Run the binary:

```bash
.bin/github.com/underscorenico/tobcast
```

Note: keep in mind that every time you run the binary, the configuration is read, so you need to update the 
tcp listen port every time you launch a new instance.

### Dependency

If you want to use tobcast as a Go dependency simply: 

```bash
go get github.com/underscorenico/tobcast/tobcast
```

and import it in your project:

```go
import "github.com/underscorenico/tobcast/pkg/tobcast"
```