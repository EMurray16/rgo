# Rgo

Connecting R and Go

Rgo allows one to connect R to Go, exporting Go functions as callable C functions. This makes it easier to use the good parts of Go (performance, reading online data, and non-embarrassing parallelization) and the good parts of R (plotting and nice data analysis libraries) into one workflow. 

This repository also contains an example of using `rgo` to export a library of simple Go functions that can be compiled and then called from R.

# Using rgo

rgo helps translate data from R's internal object representation in C (a `SEXP`) to Go objects of standard types (floats, ints, strings, etc.) and back again. This allows for easy, RAM-only passage of data from R to Go and/or Go to R by calling Go functions in R. This can be desirable for a number of reasons, which I [covered in a blogpost when I first started the project](https://overthinkdciscores.com/2018/11/20/introducing-sexp-connecting-r-and-go/).

The intended use of the `rgo` package to provide an interface for writing Go functions that can be called from R without having to worry too much about R's C internals. This is especially valuable for programmers like me who are fluent in R and Go but don't have the C chops to mess around with R's internals directly. 

## Requirements

Rgo requires a working installation of R (at least version 4.0.0) and Go (at least version 1.18). Rgo uses cgo to call R's internal C functions, which means the Go installation must have cgo enabled and there must be a C compiler. 

While rsexp contains its own header files which define the C functions called in the rsexp package, the location of the R shared libraries must also be included at compile time. This means the R libraries must be either in the default linker path, or be in one of the following directories that rsexp links automatically:

- Linux: `/usr/lib/`
- MacOS: `/Library/Frameworks/R.framework/Libraries`

Windows is neither well supported or tested in this package. Moreover, rgo does not look for a default Windows path to the R shared libraries.

If R's shared libraries are not in the default linker path or in the default locations which are included, the best solution is to either [use environment variables as specified in this SO post](https://stackoverflow.com/questions/28710276/override-an-external-packages-cgo-compiler-and-linker-flags) or to modify the contents of conversion.go file to link to the appropriate path. If you are using the former, note that it will require using the `replace` directive in your go.mod file.

In addition to the requirements for getting rgo to compile, there are additional requirements to use the package. Because [cgo does not allow for exported C types](https://golang.org/cmd/cgo/#hdr-Go_references_to_C) (see quoted text), the package which imports `rgo` must also include a link to R's internal definitions. Therefore, the file which uses the `C.SEXP` type must include a link to R's header files.

> Cgo translates C types into equivalent unexported Go types. Because the translations are unexported, a Go package should not expose C types in its exported API: a C type used in one Go package is different from the same C type used in another.

In general, this will result in a code snippet like the following:

```go
/*
#define USE_RINTERNALS // this is optional
#include <Rinternals.h>
// We need to include the shared R headers here
// One way to find this is via the rgo/rsexp directory
// Another way is to find them from your local R installation
// - Typical Linux: /usr/share/R/include/
// - Typical MacOS: /Library/Frameworks/R.framework/Headers/
// If all else fails, you can also find the required header files wherever rgo is located on your computer
// For example, on my computer all github packages are put in /Go/mod/pkg/github.com/...
#cgo CFLAGS: -I/Library/Frameworks/R.framework/Headers/
*/
import "C"
```

In order to have access to R's internal functions that are used in rgo like `TYPEOF` or `XLENGTH`, it's necessary to include a `#define USE_RINTERNALS`. Like in the rgo package itself, some functionality will also depend on linking the R shared libraries.

In order to avoid having to worry about the user's location of these header files, rgo keeps a copy in the repository. This means the headers are also available wherever `go get` saves the package files, typically in something like `/Go/mod/pkg/`. 

## Building Go functions that can be called from R

For a working example of how to use the rsexp package, refer to the `demo` package in the Rgo repository.

The Go functions are made in the `main` package. They *must* use the `C.SEXP` type as both the input and output (remember that R doesn't allow multiple returns, which means the output must be only one `C.SEXP`), and have an export comment before their signature that looks like this:

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
go build -o <package>.so -buildmode=c-shared <package>
```

## Supported Types

Rgo uses Go's generics to support as many types as possible. However, it only uses types that are easily translated and understood in both R's and Go's type systems. Therefore, it contains two types of type constraints which match R's notion of the `numeric` and `character` types.

```go
type RCharacter interface {
    ~string | ~[]byte
}
type RNumeric interface {
    ~float64 | ~float32
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

All functions to create, modify, or extract data from an `RSEXP` use these constraints as type parameters.

The only supported types that can be used as inputs to a Go function are R's numeric and character vectors and matrices. Because Go functions can have any number of inputs, supporting a more complex type like lists or data frames doesn't seem necessary. 

However, the types that can be created in Go and sent to R are much more diverse, because R functions can only have a single output. Rgo contains functions to create matrices, lists, named lists, and data frames, as well as more common vectors.

## Writing Go code using Rgo

Because the `C.SEXP` type used in rgo is different from the `C.SEXP` type that is used elsewhere, the only way to pass the type around is to use a combination of unsafe pointers, reflection, and type parameters.  The workhorse object of the `rgo` package is an `RSEXP`:

```go
type RSEXP C.SEXP
```

While Go considers a `C.SEXP` in one package to be different from the `C.SEXP` from another, the underlying data is the same. Therefore, the data can be converted within the rgo package using a combination of reflection and unsafe pointers. The two most important functions in rgo are therefore those that convert between a users `C.SEXP` and rgo's. They are called `NewRSEXP` and `ExportRSEXP`, respectively.

```go
func NewRSEXP(in any) (RSEXP, error)
func ExportRSEXP[t any](*RSEXP) (t, error)
```

Both functions use the `any` type parameter as its input, which means techincally anything can be provided by the caller.  The rgo package doesn't have knowledge of a user's `C.SEXP` type a priori, so there's no way around allowing users to provide anything. Both functions, however, perform basic type checking and will return an error if the provided variable's underlying data isn't either a `C.SEXP` or a `*C.SEXP`. 

Rgo contains functions to create, modify, and extract data from R's internal `SEXP` type. These functions wrap R's internal C functions, and use the constraints of supported types to help users avoid panics, runtime errors, or segmentatin faults.

## Writing R Code Using Rgo

Because the Go code is compiled to be executable in C, all we need to do in R is load the shared library using `dyn.load` and then call it using R's `.Call` interface:

```R
dyn.load("MYLIB.so")
outputR = .Call("MYFUNC", inputR)
```

It is important to be careful to only load the library once per R session, as loading it multiple times can result in instability. Likewise, loading a library in R, changing it in Go and then recompiling, and then loading it again in the same R session will most likely crash R.

## The Matrix Type

Because lots of R code focuses on matrices, data frames, and `data.table`s, rsexp contains an implementation of the matrix type which mirrors the R `matrix` implementation. This allows for easier matrix operations in Go and provides a Go type which will return an identical matrix back to R.

Specifically, a matrix is specified as a single vector with metadata describing the dimensions. While matrices in R can be of any numeric type, in Rgo they are always the float64 type, matching the way R will handle operations of just about any matrix.

```go
type Matrix struct {
	Nrow, Ncol int
	Data       []float64
}
```

Consistent with R's implmentation, the `Data` vector is a single concatenation of all the data, with each column serving as a vector itself. For example:

```go
// this matrix
Matrix{Nrow: 2, Ncol: 2, Data: []float64{1.1,2.2,3.3,4.4}}
// looks like this:
// [1.1 3.3
//  2.2 4.4]
```

In addition to providing the `Matrix` type, the rsexp package provides many functions and methods to get and set subsets of data within a matrix and do simple linear algebra operations.

In order to ensure matrix data quality, all matrix operation functions which can return an error first check the input matrix for internal consistency (such as the length of the data vector matching the `Nrow` and `Ncol` metadata). 

The `Matrix` struct is exported in order to allow users to be as flexible as possible in using it, but that comes with responsibility. Sloppy handling of matrices will likely result in compiler issues and/or panics at runtime. Sticking to the methods and functions provided in the package is much safer, although somewhat restricting.

# How Rgo Works

The workhorse file of rgo is conversion.go. It defines the C code used as the *go-between* (get it?) between R and Go and defines the functions that convert an `RSEXP` to a useful Go type and a useful Go type back to an `RSEXP`. 

## Extracting data from R

Rgo is based on the C interface for R's internals. More about R's internals can be found [here](https://cran.r-project.org/doc/manuals/r-release/R-ints.html), and Hadley Wickham's book [R's C Interface](http://adv-r.had.co.nz/C-interface.html) is also a good resource on the topic.

Everything in R is a `SEXP`, which is always a pointer to a `SEXPREC`, which in turn contains some header information and a pointer to the data itself. A `SEXP` can point to a `SEXPREC` of up to a couple dozen types. The rsexp package only concerns itself with 5 of them:

1. `REALSXP`, akin to a Go slice of `float64`s
2. `INTSXP`, akin to a Go slice of `int`s
3. `CHARSXP`, akin to a a Go `string`
4. `STRSXP`, akin to a Go slice of `string`s
5. `VECSXP`, which is an R list and contains no parallel in Go

In C, the type of data a `SEXP` points to can be found using the `TYPEOF` function. It returns an integer, which can be matched to the relevant types based on the rsexp's constants. When using one of the functions to convert an `RSEXP` to a Go object, they first check to make sure the type of the `SEXP` matches the type list allowed by the function. If the type doesn't match, they return an error.

The `SEXPREC` type, which points to the underlying data of the R object itself, does not explicitly point to a vector. Instead, it points to the *beginning* of the vector and the rest can be found using pointer arithmetic. Go doesn't support pointer arithmetic. Therefore, Rgo uses C  functions to extract a single element of a vector based on the original pointer location and an index. The functions that extract data from a `SEXP` determine the length of the underlying data (via `XLENGTH`), make a slice to hold it, and fill the slice one index at a time using these extractor functions.

## Sending data to R

Rgo only accepts relatively simple types from R but sends a much more varied list of types back, including named lists and data frames. This is because Go functions can support any number of inputs, but R functions can only return one, requiring more complex types to be useful.

Creating R output objects in Go is, intuitively, the same process as extracting data in reverse. First, Rgo determines the underlying R type to create, allocates the size of the vector, and fills it in element by element. Each Rgo generic type has a function to create corresponding `SEXP` objects. 

There are Go functions to facilitate the creation of more complex types as well, [which is all done by setting class attributes](https://stackoverflow.com/a/37070440) of the underlying `SEXP`.

# CRAN

Right now, there is no Rgo CRAN package, nor is there a way, using Rgo, to ship a CRAN package which runs Go under the hood rather than C or R itself. 

In principle, it should be doable to create a CRAN package which simply wraps `.Call` to run C functions built using Go. The tricky part with this is compilation. I don't think CRAN can check for a Go installation before downloading and building a package. It should be possible to bundle the package using pre-compiled versions for each operating system, but I don't make my own R packages so I'm not sure how easy that is to implement.

If someone who knows more about CRAN than I do would like to contribute, please do.