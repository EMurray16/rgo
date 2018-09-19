//this contains functions and stuff to convert from Go to R objects
package sexp


/*
#define USE_RINTERNALS
#include "Rheader/Rinternals.h"
void doubleInsert(SEXP s, int index, double v) {
	double *rp = REAL(s);
	*(rp+index) = v;
}
void intInsert(SEXP s, int index, int v) {
	int *ip = INTEGER(s);
	*(ip+index) = v;
}
// we may or may not need the link flag
#cgo LDFLAGS: -lR
*/
import "C"

import(
	//used so SEXP can be used across other packages
	"unsafe"
)

//here a reminder of the GoSEXP type
// type GoSEXP struct {
//     Point unsafe.Pointer
// }

//this function converts a slice of floats to a SEXP
func Float2sexp(in []float64) (GoSEXP) {
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.REALSXP, C.int(size))
	
	//insert the elements of the slice one at a time
	for ind,val := range in {
		C.doubleInsert(s2, C.int(ind), C.double(val))
	}
	
	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

//this function converts a slice of ints to a SEXP
func Int2sexp(in []int) (GoSEXP) {	
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.INTSXP, C.int(size))
		
	for ind,val := range in {
		C.intInsert(s2, C.int(ind), C.int(val))
	}
	
	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

//this function converts a string to a SEXP
func String2sexp(in string) (GoSEXP) {
	//convert the string to bytes
	byteslice := []byte(in)
	
	//allocate the SEXP
	size := len(byteslice)
	s2 := C.allocVector(C.INTSXP, C.int(size))
	
	//insert the bytes one at a time
	for ind,val := range byteslice {
		C.intInsert(s2, C.int(ind+1), C.int(int(val)))
	}
	
	//now make the unsafe pointer and return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}
