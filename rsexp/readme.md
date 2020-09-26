# Rgo/rsexp

rsexp helps translate data from R's internal object representation in C (a `SEXP`, hence the package name) to Go objects of standard types (floats, ints, strings, etc.) and back again. This allows for easy, RAM-only passage of data from R to Go and/or Go to R by calling Go functions in R. This can be desirable for a number of reasons, which I [covered in a blogpost when I first started the project](https://overthinkdciscores.com/2018/11/20/introducing-sexp-connecting-r-and-go/).

The intended use of the `rsexp` package to provide an easier interface for writing Go functions that can be called from R than having to worry about R's C internals. This is especially valuable for programmers who, are fluent in R and Go but don't have the C chops to mess around with R's internals. 

## Requirements

rsexp requires a working installation of R (at least version 4.0.0) and Go (at least version 1.14). rsexp uses cgo to call R's internal C functions, which means the Go installation must have cgo enabled and there must be a C compiler. 

While rsexp contains its own header files which define the C functions called in the rsexp package, the location of the R shared libraries must also be included at compile time. This means the R libraries must be either in the default linker path, or be in one of the following directories that rsexp links automatically:

- Linux: `/usr/lib/`
- MacOS: `/Library/Frameworks/R.framework/Libraries`

Note that Windows is not well supported by this package and is not tested. Moreover, rsexp does not look for a default Windows path 

If R's shared libraries are not in the default linker path or in the default locations which are included, the best solution is to either [use environment variables as specified in this SO post](https://stackoverflow.com/questions/28710276/override-an-external-packages-cgo-compiler-and-linker-flags) or to modify the contents of the conversion.go in this file to link to the appropriate path.

In addition to the requirements for getting rsexp to compile, there are additional requirements to use the package. Because [cgo does not allow for exported C types](https://golang.org/cmd/cgo/#hdr-Go_references_to_C) (see quoted text), the package which imports `rsexp` must also include a link to R's internal definitions. Therefore, the file which uses the `C.SEXP` type must include a link to R's header files.

> Cgo translates C types into equivalent unexported Go types. Because the translations are unexported, a Go package should not expose C types in its exported API: a C type used in one Go package is different from the same C type used in another.

In general, this will result in a code snippet like the following:

```go
/*
#define USE_RINTERNALS // this is optional
#include <Rinternals.h>
// We need to include the shared R headers here
// One way to find this is via the Rgo/rsexp directory
// Another way is to find them from your local R installation
// - Typical Linux: /usr/share/R/include/
// - Typical MacOS: /Library/Frameworks/R.framework/Headers/
// If all else fails, you can also find the required header files wherever Rgo is located on your computer
// For example, on my computer all github packages are put in /Go/mod/pkg/github.com/...
#cgo CFLAGS: -I/Library/Frameworks/R.framework/Headers/
*/
import "C"
```

In order to have access to R's internal functions that are used in rsexp like `TYPEOF` or `XLENGTH`, it's necessary to include a `#define USE_RINTERNALS`. Like in the rsexp package itself, some functionality will also depend on linking the R shared libraries.

In order to avoid having to worry about the user's location of these header files, rsexp keeps a copy in the repository. This means the headers are also available wherever `go get` saves the package files, typically in something like `/Go/mod/pkg/`. 

## Building Go functions that can be called from R

For a working example of how to use the rsexp package, refer to the `example` package in the Rgo repository.

The Go functions are made in the `main` package. They *must* use the `C.SEXP` type as both the input and output, and have an export comment before their signature that looks like this:

```go
//export MYFUNC
func MYFUNC(C.SEXP) C.SEXP {
```

These functions should be sure to parse the `C.SEXP` to Go objects and then format its desired output back to a `C.SEXP` for R, but otherwise can be written like normal Go. 

In order for a package to be compiled into a C library correctly, the main function needs to be defined, but it shouldn't do anything:

```go
func main() {}
```

Finally, the C library can be compiled like so:

```
go build -o <libName>.so -buildmode=c-shared <package>
```

## Writing Go code using Rgo/rsexp

Because the `C.SEXP` type used in the rsexp package is different from the `C.SEXP` type that is used elsewhere, the only way to pass the type around is to use an unsafe pointer.  Therefore, the workhorse object of the `rsexp` package is a `GoSEXP`:

```go
type GoSEXP struct {
	Point unsafe.Pointer // this should always point to a C.SEXP
}
```

While Go considers a `C.SEXP` in one package to be different from the `C.SEXP` from another, the underlying data is the same. Each package can assert that a`GoSEXP` is pointing to it's own interpretation of a `C.SEXP` and both will be correct. The `GoSEXP` has an internal method for dereferencing a `GoSEXP`:

```go
func (g GoSEXP) deref() C.SEXP {
	return *(*C.SEXP)(g.Point)
}
```

This isn't exported because the returned `C.SEXP` will not match the `C.SEXP` definition of the package that calls it. However, the code above can be effectively copy-pasted to other packages (using a function instead of a method) and it will work just fine.

Whenever the rsexp package wants to read data that comes from R, or format data to send back to R, it uses the `GoSEXP` type. While there is no enforcement mechanism in place, **a GoSEXP should always point to a `SEXP` object in C**. The rsexp package will always use a `GoSEXP` assuming this is the case. 

To make translation between packages easier, rsexp *does* export a function `NewGoSEXP` which will create a GoSEXP from a C.SEXP of any package using reflection. However, because Go doesn't have generics, the input in an empty interface. If the input is not a SEXP, or is a SEXP of the wrong type, `NewGoSEXP` will return an error, rather than panic. `NewGoSEXP` is the safest way to create a GoSEXP.

rsexp contains exported functions which convert data from a `SEXP` to Go types and functions which convert Go types to a `SEXP`. These functions use a `GoSEXP` as inputs and outputs, to allow the data to be passed from package to package. A function using these functions correctly will look something like this:

```go
//export MYFUNC
func MYFUNC(inputC C.SEXP) C.SEXP {
  // to start, we turn the C.SEXP into a GoSEXP
  inputGo := rsexp.GoSEXP{unsafe.Pointer(&input)}
  usefulInputData, err := rsexp.AsFloats(inputGoSEXP)
  
  // now we can work in Go like normal
  
  outputGO, err = rsexp.Float2sexp(usefulOutputData)
  outputC := *(*C.SEXP)(outputGo.Point)
  return outputC
}
```

## Writing R Code Using Rgo/rsexp

Because the Go code is compiled to be executable in C, all we need to do in R is load the shared library using `dyn.load` and then call it using R's `.Call` interface:

```R
dyn.load("MYLIB.so")
outputR = .Call("MYFUNC", inputR)
```

It is important to be careful to only load the library once per R session, as loading it multiple times can result in instability. Likewise, loading a library in R, changing it in Go and then recompiling, and then loading it again in the same R session will most likely crash R.

## The Rgo/rsexp Matrix Type

Because lots of R code focuses on matrices, data frames, and `data.table`s, rsexp contains an implementation of the matrix type which mirrors the R `matrix` implementation. This allows for easier matrix operations in Go and provides a Go type which will return an identical matrix back to R.

Specifically, a matrix is specified as a single vector with all elements, and some metadata for the dimensions:

```go
type Matrix struct {
	Nrow, Ncol int
	Data       []float64
}
```

The `Data` vector is ordered to have columns with adjacent indices, with the result that each row index is offset by the number of columns. Essentially, the data vector is a concatenation of several vectors, one for each column. For example:

```go
// this matrix
Matrix{Nrow: 2, Ncol: 2, Data: []float64{1.1,2.2,3.3,4.4}}
// looks like this:
// [1.1 3.3
//  2.2 4.4]
```

In addition to providing the `Matrix` type, the rsexp package provides many functions and methods to get and set subsets of data within a matrix and do simple linear algebra operations.

In order to ensure matrix data quality, all matrix operation functions which can return an error first check the input matrix for internal consistency (such as the length of the data vector matching the `Nrow` and `Ncol` metadata). 

The `Matrix` struct is exported in order to allow the user more power in its use, but that comes with responsibility. Sloppy handling of matrices will likely result in compiler issues and/or panics at runtime. Sticking to the methods and functions provided in the package is much safer, although somewhat restricting.

# How Rgo/rsexp Works

The workhorse file of rsexp is conversion.go. It defines the C code used as the *go-between* (get it?) between R and Go and defines the functions that convert a `GoSEXP` to a useful Go type and a useful Go type back to a `GoSEXP`. 

## Extracting data from R

The rsexp package is based on the C interface for R's internals. More about R's internals can be found [here](https://cran.r-project.org/doc/manuals/r-release/R-ints.html), and Hadley Wickham's book [R's C Interface](http://adv-r.had.co.nz/C-interface.html) is also a good resource on the topic.

Everything in R is a `SEXP`, which is always a pointer to a `SEXPREC`, which in turn contains some header information and a pointer to the data itself. A `SEXP` can point to a `SEXPREC` of up to a couple dozen types. The rsexp package only concerns itself with 5 of them:

1. `REALSXP`, akin to a Go slice of `float64`s
2. `INTSXP`, akin to a Go slice of `int`s
3. `CHARSXP`, akin to a a Go `string`
4. `STRSXP`, akin to a Go slice of `string`s
5. `VECSXP`, which is an R list and contains no parallel in Go

In C, the type of data a `SEXP` points to can be found using the `TYPEOF` function. It returns an integer, which can be matched to the relevant types based on the rsexp's constants. When using one of the functions to convert a `SEXP` to a Go object (one of `AsFloats`, `AsInts`, or `AsStrings`), they first check to make sure the type of the `SEXP` matches the type implied by the function name. If the type doesn't match, they return an error.

The `SEXPREC` type, which points to the underlying data of the R object itself, does not explicitly point to a vector. Instead, it points to the *beginning* of the vector and the rest can be found using pointer arithmetic. Go, however, doesn't support pointer arithmetic. Therefore, `rsexp` uses C  functions to extract a single element of a vector based on the original pointer location and an index. The functions `AsFloats`, `AsInts`, and `AsStrings` determine the length of the underlying data (via `XLENGTH`), make a slice to hold it, and fill the slice one index at a time using these extractor functions.

rsexp only handles incoming data from R as vectors of three types: `int`, `float64` (doubles in R), and `string`. Functionally, rsexp can also read a matrix from R, but this is done the same way as reading a vector of floats. This should allow enough flexibility to pass any necessary data from R to Go. 

It's important to note that while rsexp *does* have support for sending an R list (a `VECSXP`) back to R, it will not read an R list as input. Rather, the desired way to send multiple pieces of data from R to Go is to use multiple function arguments.

## Formatting data for R

Because everything in R is a `SEXP`, Go needs to create a `C.SEXP` as output. Just as there are three extractor functions, there are three SEXP creation functions for the Go types of `[]int`, `[]float64`, and `[]string`. 

These functions follow the same process as their mirrors that read data from R. They determine the length of the data, allocate the memory using R's internals (`allocVector`), and use a C function to fill in each element one at a time, using pointer arithmetic. 

In addition to the three basic data types, rsexp can also send lists back to R. In Go, a list is represented as a slice of `GoSEXP`s: 

```go
type List []GoSEXP
```

When a `List` is converted into a `SEXP`, it results in a `VECSXP` with one element corresponding to each element of the `Data` slice. Just like R lists, these elements can be of any type (that rsexp supports).

## Matrices

Matrices in rsexp are a special case. Because they are represented as simple vectors in R, they are read in Go the same way. Unfortunately, the matrix size information is stripped when it's passed to Go, so the number of rows and columns of the matrix must be known *a priori* when converting a matrix from a `SEXP` to a `Matrix`. 

There are two ways to send a `Matrix` from Go to R. First, the slice of data can be formatting using `Float2sexp`. However, the rsexp package also exports a `Matrix2sexp` function which creates a list of two elements. The first element is an integer vector containing the matrix size information as [nrow, ncol]. The second element is the underlying vector of data. To provide for easier translation, rsexp also contains an R file with a function to make translation easier.

## CRAN

Right now, there is no Rgo CRAN package, nor is there a way, using Rgo, to ship a CRAN package which runs Go under the hood rather than C or R itself. In principle, it should be doable to create a CRAN package which simply wraps `.Call` interfaces to code written in Go. Compilation could be an issue however - without having pre-compiled versions ready to go from CRAN (like for Windows), a package user may need to have Go installed to use a package that runs Rgo under the hood.

If someone who knows more about CRAN than I do would like to contribute, please do.