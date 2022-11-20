# Depdency Visualizer

A simple tool to transform Go module dependencies into mermaid graphs.

*Warning:* This does not really work.

Build with `go build`

```sh
./dependencyviz -ignorePrefix="github.com/hrisp/" -ignore="/mocks/,/cmd/,testutils" /Users/chris/get-doer
```
