# Rgo + Generics

This should be fun!

This is a development branch of Rgo that takes advantage of generics. Specifically, this branch is based on [this blogpost](https://blog.golang.org/generics-next-step), which talks about their latest [draft design](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md) and they talk about their `go2go` tool, which can run code with generics in it.

## Running Rgo with generics

The `go2go` tool is on the `dev.go2go` branch of the main Go repository. In their blogpost, the Go team links to [instructions on how to build Go from source](https://golang.org/doc/install/source), which is the only way to get the tool for now. Instead of checking out and building master, check out and build the `dev.go2go` branch.

The next issue is how to import Rgo into code that is usable by the `go2go` tool. In the [`go2go` readme](https://github.com/golang/go/blob/dev.go2go/README.go2go.md), they mention that using the `GO2PATH` environment variable is the best way to import Go code with generics. They recommend setting it accordingly:

```
export GO2PATH=$GOROOT/src/cmd/go2go/testdata/go2path
```

I recommend cloning the Rgo/generics branch to this location. From there, Rgo should be able to be imported like a normal Go standard library. It's not the most elegant solution, but it works well enough for us to start playing with generics, which is the point after all!