package rsexp

/*
#define USE_RINTERNALS
#include <Rinternals.h>
void doubleInsert(SEXP s, int index, double v) {
	double *rp = REAL(s);
	*(rp+index) = v;
}
void intInsert(SEXP s, int index, int v) {
	int *ip = INTEGER(s);
	*(ip+index) = v;
}
void listInsert(SEXP s, int index, SEXP obj) {
	SET_VECTOR_ELT(s, index, obj);
}
double doubleExtract(SEXP s, int index) {
	double output;
	//now get the pointer to the beginning of the vector in the SEXP
	double *rp = REAL(s);
	output = *(rp+index);
	return output;
}
int intExtract(SEXP s, int index) {
	int output;
	//get the pointer to the beginning of the vectors in the SEXP
	int *ip = INTEGER(s);
	output = *(ip+index);
	return output;
}
char charExtract(SEXP s, int index) {
	char output;
	const char *cp = CHAR(s);
	output = *(cp+index);
	return output;
}
void charInsert(SEXP s, int index, char* c) {
	SET_STRING_ELT(s, index, mkChar(c));
}

// we use {SRCDIR} to make sure we can always find the R header files regardless of where this file is located
#cgo CFLAGS: -I${SRCDIR}/Rheader
// we need to link the R dynamic library for the actual implementation though. By default, we can link the most common
// paths to the shared libraries for each operating system.
// Default Mac location
#cgo LDFLAGS: -L/Library/Frameworks/R.framework/Libraries
// default linux location
#cgo LDFLAGS: -L/usr/lib
#cgo LDFLAGS: -lR
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// deref is a convenience method for dereferencing a GoSEXP and casting it as a C.SEXP. We define it here instead of
// in types.go, a more natural place for it, so that we only have to import C in one file. If we import C in multiple
// files, we have to duplicate the CFLAGS and LDFLAGS for each one.
func (g GoSEXP) deref() C.SEXP {
	return *(*C.SEXP)(g.Point)
}

// NewGoSEXP creates a new GoSEXP from the input object. Because C types are not able to be exported by cgo, the input
// is the dreaded empty interface. Despite the empty interface, the input to NewGoSEXP must always be a pointer to a
// C.SEXP, like so:
//		// assume s is a C.SEXP
//		gs, err := NewGoSEXP(&s)
//
// To try and enforce as much type safety as possible, the NewGoSEXP will return an error if the input is not a *SEXP.
// It will also return a TypeMismatch error if the SEXP is not of a type that the rsexp package supports, like a list
// or a closure.
//
// For a demonstration of how to use NewGoSEXP, see the example provided in the documentation or the demo of this
// package that can be found in the same repository on Github.
func NewGoSEXP(in interface{}) (g GoSEXP, err error) {
	// This function requires in to be a C.SEXP, but there is no compile-time enforcement mechanism
	// Therefore, we will proceed as though we have a SEXP, and recover any panics
	// We assume any panics are a result of a type mismatch, where the provided interface is not a C.SEXP
	defer func() {
		if r := recover(); r != nil {
			//TODO: It would be nice to capture when a C.SEXP is provided instead of a *C.SEXP and return a more helpful error
			err = fmt.Errorf("%w: %s", NotASEXP, r)
		}
	}()

	// The underlying data of the input is the same as it would be for a rsexp.C.SEXP
	// We get the underlying data with reflection, and then make an unsafe pointer
	underlyingData := reflect.ValueOf(in)
	underlyingDataCpointer := underlyingData.Pointer()
	pointerToUnderlyingData := unsafe.Pointer(underlyingDataCpointer)
	sexpSEXP := *(*C.SEXP)(pointerToUnderlyingData)

	// now put the correctly type SEXP in the output
	g.Point = unsafe.Pointer(&sexpSEXP)

	typeEnum := C.TYPEOF(sexpSEXP)
	if !(typeEnum == REALSXP || typeEnum == INTSXP || typeEnum == STRSXP || typeEnum == CHARSXP) {
		return g, TypeMismatch
	}

	// if there's a panic, err will be TypeMismatch here
	return g, err
}

// AsFloats reads data from a SEXP into a slice of float64s. This function is only compatible with SEXPs which are of
// SEXPTYPE 14 - REALSXP. Attempts to read SEXPs of other types using this function will result in a TypeMismatch error.
func (g GoSEXP) AsFloats() ([]float64, error) {
	//cast the GoSEXP as a C.SEXP
	cs := g.deref()

	// ensure the type of the SEXP is actually integers
	if C.TYPEOF(cs) != REALSXP {
		return nil, TypeMismatch
	}

	//start by finding the length of the SEXP and making a slice
	Slen := int(C.XLENGTH(cs))
	OutSlice := make([]float64, Slen)

	//for each element of the slice, pull out the SEXP part
	for i := 0; i < Slen; i++ {
		OutSlice[i] = float64(C.doubleExtract(cs, C.int(i)))
	}

	//now return
	return OutSlice, nil
}

// Float2sexp creates a SEXP, of type REALSXP, from data contained in a slice of floats. The output of this function is
// a GoSEXP, which can be dereferenced and asserted as a C.SEXP in an external package and returned to R. In R, the
// result is a numeric (aka double) vector.
func Float2sexp(in []float64) GoSEXP {
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.REALSXP, C.long(size))

	//insert the elements of the slice one at a time
	for ind, val := range in {
		C.doubleInsert(s2, C.int(ind), C.double(val))
	}

	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

// AsInts reads data from a SEXP into a vector of ints. This function is only compatible with SEXPs which are of
// SEXPTYPE 13 - INTSXP. Attempts to read SEXPs of other types using this function will result in a TypeMismatch error.
func (g GoSEXP) AsInts() ([]int, error) {
	//cast the GoSEXP as a C.SEXP
	cs := g.deref()

	// ensure the type of the SEXP is actually integers
	if C.TYPEOF(cs) != INTSXP {
		return nil, TypeMismatch
	}

	//start by finding the length of the SEXP and making a slice
	Slen := int(C.XLENGTH(cs))
	OutSlice := make([]int, Slen)

	//for each element of the slice, pull out the SEXP part
	for i := 0; i < Slen; i++ {
		OutSlice[i] = int(C.intExtract(cs, C.int(i)))
	}

	//now return
	return OutSlice, nil
}

// Int2sexp creates a SEXP, of type INTSXP, from data contained in a slice of integers. The output of this function is
// a GoSEXP, which can be dereferenced and asserted as a C.SEXP in an external package and returned to R. In R, the
// result is an integer vector.
func Int2sexp(in []int) GoSEXP {
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.INTSXP, C.long(size))

	for ind, val := range in {
		C.intInsert(s2, C.int(ind), C.int(val))
	}

	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

// AsStrings reads data from a SEXP into a vector of strings. This function is only compatible with SEXPs which are of
// SEXPTYPE 16 - STRSXP. Attempts to read SEXPs of other types using this function will result in a TypeMismatch error.
func (g GoSEXP) AsStrings() ([]string, error) {
	// cast the GoSEXP as a C.SEXP
	cs := g.deref()

	// ensure the type of the SEXP
	if C.TYPEOF(cs) != STRSXP {
		return nil, TypeMismatch
	}

	// see how many strings we need to return
	Slen := int(C.XLENGTH(cs))
	OutSlice := make([]string, Slen)

	for stringInd := 0; stringInd < Slen; stringInd++ {
		// first, pull out the CHARSXP of the string
		charsxp := C.STRING_ELT(cs, C.long(stringInd))

		// now we want to build a string using each index of the character vector
		nChar := int(C.XLENGTH(charsxp))
		s := ""
		for charInd := 0; charInd < nChar; charInd++ {
			// pull out the specific character and add it to the string
			indChar := C.charExtract(charsxp, C.int(charInd))
			s += C.GoString(&indChar)
		}
		OutSlice[stringInd] = s
	}

	//now convert this to a string and return
	return OutSlice, nil
}

// String2sexp creates a SEXP of type STRSXP from data contained in a slice of strings. The output of this function
// is a GoSEXP, which can be dereferenced and asserted as a C.SEXP in an external package and returned to R. In R, the
// result is a string vector.
func String2sexp(in []string) GoSEXP {
	// allocate the STRSXP
	size := len(in)
	s2 := C.allocVector(C.STRSXP, C.long(size))

	for ind, str := range in {
		C.charInsert(s2, C.int(ind), C.CString(str))
	}

	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

// List2sexp creates a SEXP of type VECSXP from data contained in a List, or slice of GoSEXPs. The input of this
// function is a list of GoSEXP objects, which should already point to SEXPs of the correct types. The output of this
// function is a GoSEXP, which can be dereferenced and asserted as a C.SEXP in an external package and returned to R.
// In R the result is a list.
func List2sexp(in List) GoSEXP {
	// we could make this a method on the List type instead of a function, but then it wouldn't match the
	// interface that already exists for ints, floats, and strings

	//get the length of the list
	size := len(in)

	//start by making the list vector itself
	s := C.allocVector(C.VECSXP, C.long(size))

	//now insert the objects into the list SEXP
	for ind, obj := range in {
		C.listInsert(s, C.int(ind), obj.deref())
	}

	//wrap into an unsafe pointer and return
	outgo := GoSEXP{unsafe.Pointer(&s)}
	return outgo
}

// AsMatrix reads data from a SEXP into a Matrix type. This function is only compatible with SEXPs which are of
// SEXPTYPE 14 - REALSXP. This simply wraps a call to AsFloats and prepends the input matrix size. Because there
// is no interface to impute the size of a matrix from the SEXP in and of itself, the size of a matrix must be known
// a priori and provided as an input. If the input dimensions don't match the length of the vector in the SEXP, a
// SizeMismatch error will be returned.
func (g GoSEXP) AsMatrix(nrow, ncol int) (Matrix, error) {
	// try to get the vector as a slice of floats
	floatVec, err := g.AsFloats()
	if err != nil {
		return Matrix{}, nil
	}

	// make sure the length matches our expectations
	if len(floatVec) != nrow*ncol {
		return Matrix{}, SizeMismatch
	}

	// create the output matrix object
	m := Matrix{Nrow: nrow, Ncol: ncol, Data: floatVec}
	return m, nil
}

// Matrix2sexp creates a SEXP of type VECSXP (a list) from the data contained in a matrix. The resulting SEXP
// will always have two elements. The first is a length 2 vector of integers, containing the number of rows and
// number of columns of the matrix in that order. The second is a SEXP of type REALSXP, created by converting
// the slice of matrix data into a numeric vector.
func Matrix2sexp(in Matrix) GoSEXP {
	// Just as with lists, this could be a method on the matrix type instead of a standalone function. While
	// it may be clearer and a bit more elegant, it wouldn't match the interface users will know from ints,
	// floats, and strings. Therefore, we leave this as a function.

	// here, we send back a list with 2 objects: the dimensions and the data
	dimensions := Int2sexp([]int{in.Nrow, in.Ncol})
	data := Float2sexp(in.Data)

	// now we want to create a LIST
	list := NewList(dimensions, data)
	return List2sexp(list)
}
