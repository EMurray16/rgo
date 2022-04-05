/*
Package rgo provides a translation between R and Go using cgo.

Rgo helps translate data from R's internal representation (a SEXP) in C to o objects of standard types (floats, ints,
strings, etc.). It contains C functions that help create, modify, and extract data from a SEXP. In order to take
advantage of these functions, the workhorse type of rgo wraps its internal notion of a C.SEXP.

	type RSEXP C.SEXP

Whenever Rgo wants to read data that comes from R, or format data to send back to R, it uses the
RSEXP type. Because C types cannot be exported, the notion of a C.SEXP in the user's package is different from that
of Rgo's. Therefore, the best way to create an RSEXP is to use the NewRSEXP function.

Sending data from R to Go

Rgo uses R's internal functions to extract data from a SEXP and into a Go typed object. More about
these objects can be found in R's documentation at https://cran.r-project.org/doc/manuals/r-release/R-ints.html#SEXPs.
In short, everything in R is a SEXP, which is a pointer to a SEXPREC, which in turn contains some header information,
attributes, and a pointer to the data itself. A SEXP can point to a SEXPREC of up to a couple dozen types which map
to R's types. Rgo only concerns itself with 5 of them:

    1. REALSXP, akin to a Go slice of float64s and, when containing the dimension attributes, a matrix.
    2. INTSXP, akin to a Go slice of integers
    3. CHARSXP, akin to a a Go string
    4. STRSXP, akin to a Go slice of strings
    5. VECSXP, which is an R list and, when containing the correct attributes, data frame

In C, the type of data a SEXP points to can be found using the ''TYPEOF'' function. It returns an integer, which can
be matched to the relevant types based on the constant enumerations declared in this package. As a convenience, Rgo's
TYPEOF function wraps R's TYPEOF function. Rgo's LENGTH function also wraps R's LENGTH (also called LENGTH) function.

Rgo contains generic function which can be used to extract data from a SEXP as a desired type in Go. There are three
functions, which relate to R's numeric type, character type, and matrix type. The numeric and character functions use
a type parameter as an input so that the Go slice can be made to be the supported type which is most convenient for the
caller, within reason. These functions are:

    1. func AsNumeric[t RNumeric](r RSEXP) ([]t, error)
    1. func AsCharacter[t RCharacter](r RSEXP) ([]t, error)
    4. func AsMatrix(r RSEXP) (Matrix, error)

Each of these functions checks the SEXPTYPE of the underlying SEXP and will return an error if it doesn't match the
function that was called.

Sending data from Go to R

Sending data from Go to R is done by creating an RSEXP (which will always point to a newly created C.SEXP) from one
of the supported Go types:

    1. func NumericToRSEXP[t RNumeric](in []t) *RSEXP
    2. func CharacterToRSEXP[t RCharacter](in []t) *RSEXP
    3. func MatrixToRSEXP(in Matrix) *RSEXP

This SEXPTYPE of the output SEXP from these functions will match the R internal type which makes the most sense.

Because R does not allow functions to have multiple returns, the preferred way to return multiple pieces of data
from a function is a list. Therefore, Rgo contains functions to create three types of lists: a generic list,
a named list, and a data frame.

    1. func MakeList(in ...*RSEXP) *RSEXP
    2. func MakeNamedList(names []string, data ...*RSEXP) (*RSEXP, error)
    3. func MakeDataFrame(rowNames, colNames []string, dataColumns ...*RSEXP) (*RSEXP, error)

The functions to create named lists and data frames enforce data quality so that valid R objects can be created. This
includes enforcing the number of names and objects provided, checking the lengths of all the columns provided in a data
frame, and making sure no nested objects (like lists or data frames themselves) are provided as columns for data frames.

In order to send data back to R, it must be a C.SEXP that matches the user's notion of a SEXP, not Rgo's. Therefore, Rgo
provides a generic ExportRSEXP function which is used to create a user's C.SEXP that has the same data as the RSEXP that
has been made using the rgo library. Callers provide a type parameter, which should always be their C.SEXP:

    mySEXP, err := [C.SEXP](rgoRSEXP)

If their provided type is not a C.SEXP (or a *C.SEXP which is also acceptable) an error is returned.

Building Your Package and Calling Functions in R

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

For a more complete demonstration, see the example below, or the demo package at https://github.com/EMurray16/rgo/demo.

Example

The code below contains a functional example of a Go function that can be called from R:

    package main

    // #include <Rinternals.h>
    // We need to include the shared R headers here
    // One way to find this is via the rgo directory
    // Another way is to find them from your local R installation
    // - Typical Linux: /usr/share/R/include/
    // - Typical MacOS: /Library/Frameworks/R.framework/Headers/
    // If all else fails, you can also find the required header files wherever rgo is located on your computer
    // For example, on my computer all github packages are put in /Go/mod/pkg/github.com/...
    // #cgo CFLAGS: -I/Go/mod/pkg/github.com/EMurray16/rgo/Rheader/
    import "C"
    import(
        "github.com/EMurray16/rgo"
    )

    //export DoubleVector
    func DoubleVector(input C.SEXP) C.SEXP {
        // cast the incoming SEXP as a GoSEXP
        r, err := rgo.NewRSEXP(&input)
        if err != nil {
            fmt.Println(err)
            return nil
        }

        // create a slice from the SEXPs data
        floats, err := rgo.AsNumeric[float64](r)
        if err != nil {
            fmt.Println(err)
            return nil
        }

        // double each element of the slice
        for i, _ := range floats {
            floats[i] *= 2
        }

        // create a SEXP and GoSEXP from the new data
        outputRSEXP := rgo.NumericToSEXP(floats)

		mySEXP, err := rgo.ExportRSEXP[C.SEXP](outputSEXP)
		if err != nil {
            fmt.Println(err)
            return nil
        }

        return mySEXP
    }

Once it is compiled to a shared library, the function can be called using R's .Call() interface:

	input = c(0, 2.71, 3.14)
	output = .Call("DoubleVector", input)
	print(output)

The result would look like this:

	[0, 5.52, 6.28]

*/
package rgo
