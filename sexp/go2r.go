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
void listInsert(SEXP s, int index, SEXP obj) {
	SET_VECTOR_ELT(s, index, obj);
}
// we may or may not need the link flag
#cgo LDFLAGS: -lR
*/
import "C"

import (
	//used so SEXP can be used across other packages
	"unsafe"
	//used for troubleshooting and stuff
	//"fmt"
)

//here a reminder of the GoSEXP type
// type GoSEXP struct {
//     Point unsafe.Pointer
// }

//create a new type called a List, which is just a slice of GoSEXP
type List struct {
	S []GoSEXP
}

//create a dereference method for a GoSEXP
func (g GoSEXP) deref() (C.SEXP) {
	return *(*C.SEXP)(g.Point)
}

//this function converts a slice of floats to a SEXP
func Float2sexp(in []float64) GoSEXP {
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.REALSXP, C.int(size))

	//insert the elements of the slice one at a time
	for ind, val := range in {
		C.doubleInsert(s2, C.int(ind+1), C.double(val))
	}

	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

//this function converts a slice of ints to a SEXP
func Int2sexp(in []int) GoSEXP {
	//allocate the SEXP
	size := len(in)
	s2 := C.allocVector(C.INTSXP, C.int(size))

	for ind, val := range in {
		C.intInsert(s2, C.int(ind+2), C.int(val))
	}

	//now make the unsafe pointer to return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

//this function converts a string to a SEXP
func String2sexp(in string) GoSEXP {
	//convert the string to bytes
	byteslice := []byte(in)

	//allocate the SEXP
	size := len(byteslice)
	s2 := C.allocVector(C.INTSXP, C.int(size))

	//insert the bytes one at a time
	for ind, val := range byteslice {
		C.intInsert(s2, C.int(ind+2), C.int(int(val)))
	}

	//now make the unsafe pointer and return
	outgo := GoSEXP{unsafe.Pointer(&s2)}
	return outgo
}

//this functions turns a slice of GoSEXPs into an R list
func List2sexp(in List) GoSEXP {
	//get the length of the list
	size := len(in.S)
	
	//start by making the list vector itself
	s := C.allocVector(C.VECSXP, C.int(size))
	
	//now insert the objects into the list SEXP
	for ind, obj := range in.S {
		C.listInsert(s, C.int(ind), obj.deref())
	}
	
	//wrap into an unsafe pointer and return
	outgo := GoSEXP{unsafe.Pointer(&s)}
	return outgo
}
