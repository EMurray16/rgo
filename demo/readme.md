# Rgo/demo

This folder contains a shiny app which serves to demonstrate how to use the Rgo/rsexp package. It contains:

- app.r : the shiny app
- demo.go : the Go file with simple functions to be called by R

The demo.go file can be built by compiling a C shared library:

```
go build -o demo.so -buildmode=c-shared demo.go
```

Note that it is important that cgo can find the R header files. While demo.go will automatically link the most common locations, it may require some editing to compile. For more information, see the Rgo/rsexp documentation.

You can see this shiny app in action [at my website](https://overthinkDCIscores.com/interactive/RgoDemo).

