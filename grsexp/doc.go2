/*
Package rsexp provides a translation between R and Go using cgo.

rsexp helps translate data from R's internal representation (a SEXP, hence the package name) in C to
Go objects of standard types (floats, ints, strings, etc.). In cgo, C objects are always unexported, so the
only way to use this package is via unsafe pointers, which can be cast both within the rsexp package and
in other packages as a C.SEXP.

Hence, the workhorse object of the rsexp package is a GoSEXP:
	type GoSEXP struct {
		Point unsafe.Pointer
	}

Whenever the rsexp package wants to read data that comes from R, or format data to send back to R, it uses the
GoSEXP object. While there is no enforcement mechanism, a GoSEXP should always point to a SEXP object in C.
The rsexp package will always use a GoSEXP assuming this is the case. Therefore, the cleanest and safest (also most
convenient) way to create a new GoSEXP object is with the NewGoSEXP function.

Internally, the GoSEXP type has a method which dereferences the unsafe pointer and casts it as a C.SEXP:
	func (g GoSEXP) deref() C.SEXP {
		return *(*C.SEXP)(g.Point)
	}

In other packages, a similar function can be written to do the same thing. In a perfect world, the rsexp package
would export this method for other packages to use, but because C types are always unexported in cgo, rsexp's internal
notion of a C.SEXP will be different from the package that imports it. In this specific case, Go's strict type
safety is a bit of a hindrance.

Sending data from R to Go

The rsexp package uses R's internal functions to extract data from a SEXP and into a Go typed object. More about
these objects can be found in R's documentation at https://cran.r-project.org/doc/manuals/r-release/R-ints.html#SEXPs.
In short, everything in R is a SEXP, which is a pointer to a SEXPREC, which in turn contains some header information
and a pointer to the data itself. A SEXP can point to a SEXPREC of up to a couple dozen types. The rsexp package
only concerns itself with 5 of them:

    1. REALSXP, akin to a Go slice of float64s
    2. INTSXP, akin to a Go slice of integers
    3. CHARSXP, akin to a a Go string
    4. STRSXP, akin to a Go slice of strings
    5. VECSXP, which is an R list and as akin to a rsexp.List. It has no parallel in base Go.

In C, the type of data a SEXP points to can be found using the ``TYPEOF'' function. It returns an integer, which can
be matched to the relevant types based on the constant enumerations declared in this package.

The GoSEXP has several methods which will pull the data out of an underlying SEXP and into a Go slice (or Matrix).
They are:

    1. func (g GoSEXP) AsFloats() ([]float64, error)
    2. func (g GoSEXP) AsInts() ([]int, error)
    3. func (g GoSEXP) AsStrings() ([]string, error)
    4. func (g GoSEXP) AsMatrix(nrow, ncol int) (Matrix, error)

Each of these functions checks the SEXPTYPE of the underlying SEXP and will return an error if it
doesn't match the method that was called.

Matrices are a special case. In R, matrices are represented internally as a single vector containing all of the
underlying data. However, using cgo to get the dimension metadata along with the underlying data is not feasible.
Therefore, the user needs to know the size of the matrix being sent from R to Go *a priori*.

Sending data from Go to R

Sending data from Go to R is done by creating a GoSEXP (which will always point to a newly created C.SEXP) from one
of the supported Go types:

    1. func Float2sexp(in []float64) GoSEXP
    2. func Int2sexp(in []int) GoSEXP
    3. func String2sexp(in []string) GoSEXP
    4. func Matrix2sexp(in Matrix) GoSEXP

This SEXPTYPE of the output SEXP from these functions will match the R internal type which makes the most sense.

Because R does not allow functions to have multiple returns, the preferred way to return multiple pieces of data
from a function is a list. Therefore, the rsexp package contains a wrapper type List, which is just a slice of
GoSEXP objects.

	type List []GoSEXP

One advantage of using a GoSEXP is that the SEXP underlying each element of the List can be of any type without
violating Go's type safety. R lists work the same way. The function ``List2sexp'` will condense all the SEXPs
underlying the List to a single SEXP/GoSEXP that can be sent back to R as a VECSXP, or list.

Once again, matrices are a special case. Because the underlying data is simply a vector which is indistinguishable
from a normal numeric vector, the Matrix2sexp function returns a GoSEXP that points to a list rather than a numeric
vector. The list always has two elements. The first is a pair of integers with the size information, and the second
is the data itself. That way, the matrix is easier to create in R than it would be otherwise:

    # This is some R code to quickly make an R matrix from a matrix, formatted as a list, sent from Go
    rsexp.ParseGoMatrix <- function(rsexpMatrix) {
        # The matrix is the second element of the list, while the first is an integer vector that is c(nrow,ncol)
        outMat = matrix(data=rsexpMatrix[[2]], nrow=rsexpMatrix[[1]][1], ncol=rsexpMatrix[[1]][2])
        return(outMat)
    }

In order to send data back to R, a GoSEXP must be dereferenced and cast as a C.SEXP, similar to rsexp's internal
deref method. For convenience, the following function can be copy/pasted into other packages:

    func derefGoSEXP(g GoSEXP) C.SEXP {
         return *(*C.SEXP)(g.Point)
    }

Building Your Package and Calling Function in R

In order for a package of Go functions to be callable from R, they must take any number of C.SEXP objects as input and
return a single C.SEXP object. They also need to be marked to be exported, by including an export statement immediately
above the function signature. Note that if there is a space between the comment slashes and the export mark, Go will
parse it as a vanilla comment and the function won't be exported.

	//export DoubleVector
	func DoubleVector(C.SEXP) C.SEXP {}

The Go package must then be compiled to a C shared library:

	go build -o <libName>.so -buildmode=c-shared <package>

Finally, the Go functions can be called in R using the .Call function:

    output = .Call("DoubleVector", input)

For a more complete demonstration, see the example below, or the demo package at https://github.com/EMurray16/Rgo/demo.

Example

The code below contains a functional example of a Go function that can be called from R:

    package main

    // #include <Rinternals.h>
    // We need to include the shared R headers here
    // One way to find this is via the Rgo/rsexp directory
    // Another way is to find them from your local R installation
    // - Typical Linux: /usr/share/R/include/
    // - Typical MacOS: /Library/Frameworks/R.framework/Headers/
    // If all else fails, you can also find the required header files wherever Rgo is located on your computer
    // For example, on my computer all github packages are put in /Go/mod/pkg/github.com/...
    // #cgo CFLAGS: -I/Go/mod/pkg/github.com/EMurray16/Rgo/rsexp/Rheader/
    import "C"
    import(
        "github.com/EMurray16/Rgo/rsexp"
    )

    //export DoubleVector
    func DoubleVector(input C.SEXP) C.SEXP {
        // cast the incoming SEXP as a GoSEXP
        gs, err := rsexp.NewGoSEXP(&input)
        if err != nil {
            fmt.Println(err)
            return nil
        }

        // create a slice from the SEXPs data
        floats, err := gs.AsFloats()
        if err != nil {
            fmt.Println(err)
            return nil
        }

        // double each element of the slice
        for i, _ := range floats {
            floats[i] *= 2
        }

        // create a SEXP and GoSEXP from the new data
        outputGoSEXP := rsexp.Float2sexp(floats)

        // return the result, dereferenced and casted as a C.SEXP
        return *(*C.SEXP)(outputGoSEXP.Point)
    }

Once it is compiled to a shared library, the function can be called using R's .Call() interface:

	input = c(0,2.71,3.14)
	output = .Call("DoubleVector", input)
	print(output)

The result would look like this:

	[0, 5.52, 6.28]

*/
package rsexp
