# Repro kit

Exact steps to reproduce a green build:

```
go version
# expect go1.21+

go mod tidy

go test ./...

for f in examples/*.mars; do echo "==> $f"; go run ./cmd/mars run "$f" || echo "FAILED: $f"; done | cat
```

