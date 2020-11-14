package rsexp

import (
	"errors"
	"unsafe"
)

// These constants are enumerations of the SEXPTYPEs that are part of R's internals. There are about 2 dozen in all,
// the rsexp package only supports 5 of them.
const (
	CHARSXP = 9
	INTSXP  = 13
	REALSXP = 14

	// A STRSXP is actually a vector of strings, where each element points to a CHARSXP.
	STRSXP = 16

	// It's not obvious from the name, but a VECSXP is a list. Each element of a
	// VECSXP is a SEXP and can be of any type.
	VECSXP = 19
)

// TypeMismatch is most often returned from an AsX method when the caller tries to extract the incorrect type from
// a SEXP.
var TypeMismatch = errors.New("input SEXP type does not match desired output type")

// NotASEXP is returned by NewGoSEXP when it cannot coerce the input object into a *C.SEXP.
var NotASEXP = errors.New("non-SEXP object provided to a function that needs a SEXP")

// Any Matrix function or method which can return an error will first check the input Matrix or matrices for
// validity and return an ImpossibleMatrix error if there is an inconsistency between the matrix data and
// dimensions.
var ImpossibleMatrix = errors.New("matrix size and underlying data length are not compatible")

// All matrix operations check inputs for validity and will return errors where applicable.
var (
	SizeMismatch     = errors.New("operation is not possible with given input dimensions")
	InvalidIndex     = errors.New("given index is impossible (ie < 0)")
	IndexOutOfBounds = errors.New("index is out of bounds (ie too large)")
)

// GoSEXP wraps an unsafe pointer, which should always point towards a C.SEXP object.
// Because cgo doesn't allow for the exporting of C types, a GoSEXP is used as a translation object to pass
// a C.SEXP into the rsexp package. The preferred way to create a new GoSEXP object is with the function NewGoSEXP.
//
// A GoSEXP can be dereferenced in any package and asserted as a C.SEXP. Internally,
// the rsexp package uses an unexported method to do this:
//
//     func (g GoSEXP) deref() C.SEXP {
//         return *(*C.SEXP)(g.Point)
//     }
//
// Other packages can interact with a GoSEXP in the exact same way and get the same results, but must have their own
// dereference implementation which returns their own package's definition of a C.SEXP.
type GoSEXP struct {
	Point unsafe.Pointer
}

// List is a Go correlate to R's lists. The List type is a vector of GoSEXPs, which can be of any type. Just as
// in R, a List is the preferred way to return multiple objects from a function. However, a List is not the preferred
// way to to provide multiple inputs - both R and Go support any number of function arguments.
type List []GoSEXP

// NewList is a convenience function for creating a List from any number GoSEXPs
func NewList(s ...GoSEXP) List {
	return s
}

// Matrix is a representation of a matrix in Go that mirrors how matrices are represented in R. The Matrix
// contains a vector of all the data, and a header of two integers that contain the dimensions of the matrix.
// The Data vector is organized so that column indices are together, but row indices are not. In other words, the
// data can be thought of as a concatenation of several vectors, each of which contains the data for one column.
//
// For example, the following Matrix:
//     Matrix{Nrow: 3, Ncol: 2, Data: []float64{1.1,2.2,3.3,4.4,5.5,6.6}}
// will look like this:
//    [1.1, 4.4
//     2.2, 5.5
//     3.3, 6.6]
//
// Matrix data is accessed using 0-based indexing, which is natural in Go but differs from R. For example, the 0th
// row in the example matrix is [1.1, 4.4], while the "1st" row is [2.2, 5.5].
type Matrix struct {
	// The Matrix header - two integers which specify its dimension
	Nrow, Ncol int

	// The data in a matrix is represented as a single slice of data
	Data []float64
}

// isSizeValid checks to make sure a matrix's dimensions and data length match.
func (m *Matrix) isSizeValid() bool {
	return m.Nrow*m.Ncol == len(m.Data)
}
