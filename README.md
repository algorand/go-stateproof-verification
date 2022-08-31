
State Proof Verification
====================


A Go implementation of functionality required to verify Algorand state proofs.
The functions exported in stateproof.verifier provide the verification interface.

# Install

```bash
go get github.com/algorand/go-stateproof-verification
```
Alternatively the same can be achieved if you use import in a package:

```bash
import "github.com/algorand/go-stateproof-verification"
```
and run go get without parameters.

# Testing

```go
go test ./test
```