//this is for converting R objects to Go objects
package sexp

/*
#define USE_RINTERNALS
#include "Rheader/Rinternals.h"
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
//we may or may not need the linkage
#cgo  LDFLAGS: -lR
*/
import "C"

import (
	//used so SEXP can be used across other packages
	"unsafe"
)

//create the GoSEXP type to make the C.SEXP from this package accessible to others
type GoSEXP struct {
	Point unsafe.Pointer
}

//this function converts a SEXP to a slice of integers
func AsInts(s GoSEXP) []int {
	//cast the GoSEXP as a C.SEXP
	cs := *(*C.SEXP)(s.Point)

	//start by finding the length of the SEXP and making a slice
	Slen := int(C.XLENGTH(cs))
	OutSlice := make([]int, Slen)

	//for each element of the slice, pull out the SEXP part
	for i := 0; i < Slen; i++ {
		OutSlice[i] = int(C.intExtract(cs, C.int(i+2)))
	}

	//now return
	return OutSlice
}

//this function converts a SEXP to a slice of floats
func AsFloats(s GoSEXP) []float64 {
	//cast the GoSEXP as a C.SEXP
	cs := *(*C.SEXP)(s.Point)

	//start by finding the length of the SEXP and making a slice
	Slen := int(C.XLENGTH(cs))
	//Stype := C.TYPEOF(cs)
	OutSlice := make([]float64, Slen)

	//for each element of the slice, pull out the SEXP part
	for i := 0; i < Slen; i++ {
		OutSlice[i] = float64(C.doubleExtract(cs, C.int(i+1)))
	}

	//now return
	return OutSlice
}

//this function converts a SEXP to a string
func AsString(s GoSEXP) string {
	//start by pulling the Ints (which in this case we know are bytes)
	IntSlice := AsInts(s)

	//now make the slice of bytes
	ByteSlice := make([]byte, len(IntSlice))
	for i := 0; i < len(IntSlice); i++ {
		ByteSlice[i] = byte(IntSlice[i])
	}

	//now convert this to a string and return
	OutString := string(ByteSlice)
	return OutString
}
