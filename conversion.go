package rgo

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
	"strconv"
	"strings"
	"unsafe"
)

// RSEXP is the workhorse type of the rsexp package. It is an identical to R's SEXP implementation in C. It is used
// for any operation that deals with sending, receiving, or modifying data that moves from R to Go or from Go back to
// R.
//
// The RSEXP type exists because Cgo does not allow packages to export C types. Therefore, the C.SEXP that is
// defined in the rsexp package is different from a C.SEXP in the caller's main package. The RSEXP type acts as a
// go-between, allowing the use of functions in this package to extract, modify, or create SEXP objects without
// needing to write additional C code.
type RSEXP C.SEXP

// TYPEOF is Rgo's convenience wrapper of R's internal TYPEOF function. It returns the type of an RSEXP, matching the
// enumerations used by R itself.
func TYPEOF(r RSEXP) RSEXPTYPE {
	return RSEXPTYPE(int(C.TYPEOF(r)))
}

// LENGTH is Rgo's convenience wrapper of R's internal LENGTH/XLENGTH function. It returns the length of an RSEXP.
func LENGTH(r RSEXP) int {
	return int(C.XLENGTH(r))
}

// NewRSEXP creates an RSEXP object from the function input. It attempts to create an RSEXP object from any input, but
// only succeeds if the provided type is a C.SEXP or a *C.SEXP.
//
// It would be ideal for this function to have a more limited set of input types (like only those that can be coerced
// to a C.SEXP), but checking the "coercibility" of an input (without knowing the universe of input types a priori) is
// impossible at compile time in Go.
//
// Generally speaking, a failed type cast in Go results in a panic. NewRSEXP returns useful errors to the fullest extent
// possible, but guaranteeing no runtime failures or panics is impossible.
//
// NewRSEXP uses a combination of reflection and the unsafe package to coerce the input into Rgo's C.SEXP type.
// First, it creates an unsafe pointer to the underlying data. Then, it uses reflection to verify that the input
// type is either a C.SEXP or a *C.SEXP. Then, it performs the coercion. If this fails at any point, it returns a
// NotASEXP error. If the input is a SEXP, but not one of the types supported by the Rgo package, then it returns an
// UnsupportedType error.
func NewRSEXP(in any) (r RSEXP, err error) {
	// This function requires in to be a C.SEXP or a *C.SEXP, but there is no compile-time enforcement
	// mechanism as the function argument is an empty interface. Therefore, we will proceed as though we
	// have a SEXP, and recover any panics which are likely caused by trying to cast something as a SEXP
	// that isn't. We assume this is the cause of any panic.
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%w. more detail: %s", NotASEXP, rec)
		}
	}()

	// the first thing to check is the input type. It should be a _Ctype_SEXP or *_Ctype_SEXP

	// first, we want to create a type-unsafe pointer to the underlying data provided by the caller
	underlyingData := reflect.ValueOf(in)
	underlyingDataCpointer := underlyingData.Pointer()
	pointerToUnderlyingData := unsafe.Pointer(underlyingDataCpointer)

	// we can use reflect.TypeOf to get the type of the input. If it is a *C.SEXP, we need to dereference the pointer
	// before coercion. If it is a C.SEXP we can coerce right away. If it is neither of these, return a NotASEXP error.
	// if it contains *C.SEXP but not C.SEXP, we need to dereference the pointer before coercion
	typeString := reflect.TypeOf(in).String()
	// fmt.Println(typeString)
	if strings.Contains(typeString, "*") && strings.Contains(typeString, "_Ctype_SEXP") {
		// this means we need to dereference the pointer
		rPointer := *(*C.SEXP)(pointerToUnderlyingData)
		r = RSEXP(rPointer)
	} else if strings.Contains(typeString, "_Ctype_SEXP") {
		r = RSEXP((C.SEXP)(pointerToUnderlyingData))
	} else {
		r = nil
		return r, NotASEXP
	}

	typeEnum := TYPEOF(r)
	// Even if we have a C.SEXP, we still have no guarantee that the SEXP is of a type supported type
	if !(typeEnum == REALSXP || typeEnum == INTSXP || typeEnum == STRSXP || typeEnum == CHARSXP) {
		// fmt.Println(typeEnum)
		return r, UnsupportedType
	}

	// if there's a panic, err will be TypeMismatch here
	return r, err
}

// AsNumeric extracts data from the input RSEXP and returns it as a slice of the given type parameter. The data is the
// same data that is contained in the RSEXP, but a new copy that can be modified independently. If the underlying data
// cannot be coerced into numeric data, the TypeMismatch error is returned.
func AsNumeric[t RNumeric](r RSEXP) (out []t, err error) {
	//start by finding the length of the SEXP and making a slice
	Slen := LENGTH(r)
	out = make([]t, Slen)

	rsexpType := TYPEOF(r)
	if rsexpType != INTSXP && rsexpType != REALSXP {
		return nil, TypeMismatch
	}

	// in order to determine the type (so that we call the right extraction function) we
	// create an interface of type t we can perform a type switch on
	for i := 0; i < Slen; i++ {
		switch rsexpType {
		case INTSXP:
			out[i] = t(C.intExtract(r, C.int(i)))
		case REALSXP:
			out[i] = t(C.doubleExtract(r, C.int(i)))
		}
	}

	return out, nil
}

// AsMatrix returns a matrix based on the input RSEXP. All matrices must contain doubles/float64s with a dimension
// attribute. The data returned by this function is a copy of the data contained in the RSEXP that can be modified
// independently. If the data in the RSEXP cannot be coerced into a matrix, the TypeMismatch error is returned.
func AsMatrix(r RSEXP) (out Matrix, err error) {
	rsexpType := TYPEOF(r)
	// technically we could accept integer matrices but for now we won't
	if rsexpType != REALSXP {
		return out, TypeMismatch
	}
	dataVec, err := AsNumeric[float64](r)
	if err != nil {
		return out, err
	}

	out.Data = dataVec
	out.Nrow = int(C.nrows(r))
	out.Ncol = int(C.ncols(r))

	return out, err
}

// AsCharacter extracts the data from the input RSEXP and returns it as a slice of the given type parameter. The resulting
// slice that contains the same data as the contents of the RSEXP, but a new copy that can be modified independently.
// If the underlying data connot be coerced into string data, the TypeMismatch error is returned.
func AsCharacter[t RCharacter](r RSEXP) (out []t, err error) {
	//start by finding the length of the SEXP and making a slice
	Slen := LENGTH(r)
	out = make([]t, Slen)

	rsexpType := TYPEOF(r)
	if rsexpType != STRSXP {
		return nil, TypeMismatch
	}

	for i := 0; i < Slen; i++ {
		// first, pull out the CHARSXP of the string
		charsxp := C.STRING_ELT(r, C.long(i))

		// we want to build a string using each index of the character vector
		// to do this we'll use a byte slice as the go-between (get it?)
		nChar := LENGTH(RSEXP(charsxp))
		goBytes := make([]byte, 0, nChar)
		for charInd := 0; charInd < nChar; charInd++ {
			// pull out the specific character and add it to the slice
			indChar := C.charExtract(charsxp, C.int(charInd))
			charAsByte := C.GoBytes(unsafe.Pointer(&indChar), 1)
			// using append here is safer, in case we unexpectedly get more than 1 byte from charExtract
			goBytes = append(goBytes, charAsByte...)
		}
		// convert the assembled slice to a string
		out[i] = t(goBytes)
	}

	return out, nil
}

// NumericToRSEXP converts a slice of numeric data into a C.SEXP, represented by the returned RSEXP data. The R
// representation will have the same data as the input slice and be the REALSXP type (aka a double in R). Because
// the intent of this function is to prepare data to be sent back to R, which largely treats doubles and integers the
// same, this function cannot return an RSEXP of type INTSXP.
func NumericToRSEXP[t RNumeric](in []t) *RSEXP {
	size := len(in)
	s := C.allocVector(C.REALSXP, C.long(size))

	for i, num := range in {
		C.doubleInsert(s, C.int(i), C.double(float64(num)))
	}

	out := RSEXP(s)
	return &out
}

// MatrixToRSEXP converts a Matrix a C.SEXP, represented by the returned RSEXP data. The R representation
// will have the same data and dimensions as the input Matrix and be of the REALSXP type (aka a double in R).
func MatrixToRSEXP(in Matrix) *RSEXP {
	s := NumericToRSEXP(in.Data)

	dimSEXP := NumericToRSEXP([]int{in.Nrow, in.Ncol})

	C.setAttrib(*s, C.R_DimSymbol, *dimSEXP)

	return s
}

// CharacterToRSEXP converts a slice of strings (or byte slices) into a C.SEXP, represented by the returned RSEXP
// data. The R representation will have the same data as the input slice and be the STRSXP type (aka the character type
// in R).
func CharacterToRSEXP[t RCharacter](in []t) *RSEXP {
	size := len(in)
	s := C.allocVector(C.STRSXP, C.long(size))

	for i, str := range in {
		C.charInsert(s, C.int(i), C.CString(string(str)))
	}

	out := RSEXP(s)
	return &out
}

// MakeList creates an R list from the provided inputs and returns its representing RSEXP object. Unlike MakeDataFrame
// and MakeNamedList, there are no restrictions on the data that is provided.
func MakeList(in ...*RSEXP) *RSEXP {
	//start by making the list vector itself
	s := C.allocVector(C.VECSXP, C.long(len(in)))

	//now insert the objects into the list SEXP
	for ind, obj := range in {
		C.listInsert(s, C.int(ind), *obj)
	}

	// wrap the result and return
	out := RSEXP(s)
	return &out
}

// MakeNamedList creates a named list based on the provided input names and data and returns its representing RSEXP
// object. Users must provide the names and each element of the list. If number of names and elements provided do not
// match, a LengthMismatch error is returned. Elements of a named list may be of any valid SEXP type.
func MakeNamedList(names []string, data ...*RSEXP) (*RSEXP, error) {
	// check the length of the names versus input data
	if len(names) != len(data) {
		return nil, LengthMismatch
	}

	//start by making the list vector itself
	s := C.allocVector(C.VECSXP, C.long(len(data)))

	//now insert the objects into the list SEXP
	for ind, obj := range data {
		C.listInsert(s, C.int(ind), *obj)
	}

	// set the names attribute
	nameSEXP := CharacterToRSEXP(names)
	C.setAttrib(s, C.R_NamesSymbol, *nameSEXP)

	// wrap the result and return
	out := RSEXP(s)
	return &out, nil
}

// MakeDataFrame creates an R data frame from the provided inputs and returns its representing RSEXP object.
//
// It creates a data frame based on the provided data columns (as RSEXP objects) and column names. Users may directly
// provide row names or an empty slice. If an empty slice is provided, then row names are automatically generated
// as the row indexes, starting at 1 instead of 0 to be consistent with R.
//
// MakeDataFrame is strict about making sure the provided inputs can create a valid R data frame. The number of column
// names and data columns provided must match. Likewise, the lengths of all the provided data columns must match
// and, if row names are provided, match the number of provided row names. If any of these conditions are false, then
// a LengthMismatch error will be returned with more detail about which condition failed.
//
// MakeDataFrame also checks that the types of all the provided data columns are valid types according to Rgo and
// that they can be used to create a column in a data frame. If these conditions are not met, an UnsupportedType error
// will be returned. Right now, this list includes integer vectors, real vectors, and string vectors only.
// Examples of invalid types include lists, data frames, or other nested SEXP objects.
func MakeDataFrame(rowNames, colNames []string, dataColumns ...*RSEXP) (*RSEXP, error) {
	// first, check to make sure the number of column names and number of columns match
	if len(colNames) != len(dataColumns) {
		return nil, fmt.Errorf("%w: problem is number of column names vs. columns", LengthMismatch)
	}

	// next, check to make sure the length of all the provided dataColumns columns match
	nRowsProvided := LENGTH(*dataColumns[0])
	for _, dataCol := range dataColumns {
		if LENGTH(*dataCol) != nRowsProvided {
			return nil, fmt.Errorf("%w: problem is length of provided dataColumns columns", LengthMismatch)
		}
		colType := TYPEOF(*dataCol)
		if colType != INTSXP && colType != REALSXP && colType != STRSXP {
			return nil, UnsupportedType
		}
	}

	// do we need to create column names?
	var autoRownames bool
	if len(rowNames) == 0 {
		autoRownames = true
	}

	// if not, then we should check the length
	if !autoRownames {
		if len(rowNames) != int(nRowsProvided) {
			return nil, fmt.Errorf("%w: problem is number of row names vs. rows", LengthMismatch)
		}
	}

	// if we need to make row names, make them now
	if autoRownames {
		rowNames = make([]string, nRowsProvided)
		for i := 0; i < int(nRowsProvided); i++ {
			rowNames[i] = strconv.Itoa(i)
		}
	}

	// under the hood, a dataColumns frame is just a list with attributes. Make the list first
	r := C.allocVector(C.VECSXP, C.long(len(dataColumns)))
	for ind, obj := range dataColumns {
		C.listInsert(r, C.int(ind), *obj)
	}

	// in order to set the attributes, we need the row and column names to be SEXPs
	rowSEXP := CharacterToRSEXP(rowNames)
	colSEXP := CharacterToRSEXP(colNames)

	// now we set the attributes
	C.setAttrib(r, C.R_ClassSymbol, C.ScalarString(C.mkChar(C.CString("data.frame"))))
	C.setAttrib(r, C.R_RowNamesSymbol, *rowSEXP)
	C.setAttrib(r, C.R_NamesSymbol, *colSEXP)

	out := RSEXP(r)
	return &out, nil
}

// ExportRSEXP converts an input RSEXP object into the caller's provided C.SEXP type. It is used as a final function to
// prepare data to be sent back to R. The input type is any, because the rsexp package cannot anticipate the strict
// C.SEXP type used by the caller. Like NewRSEXP, ExportRSEXP checks the type parameter provided using reflection and returns
// a NotASEXP error if it is not a C.SEXP.
//
// The intent of ExportRSEXP is that it is always called with the user providing their C.SEXP as the type parameter, like so:
//
//     mySEXP, err := [C.SEXP](rgoRSEXP)
func ExportRSEXP[t any](r *RSEXP) (out t, err error) {
	// this function is like NewRSEXP, but creates the users C.SEXP to send back to R
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%w. more detail: %s", NotASEXP, rec)
		}
	}()

	// check the input type to make sure it's a C.SEXP from the user
	typeString := reflect.TypeOf(out).String()
	// fmt.Println(typeString)
	// We know the provided type should be a C.SEXP, so if the type string doesn't contain "_Ctype_SEXP" then we know
	// to return an error
	if !strings.Contains(typeString, "_Ctype_SEXP") {
		return out, NotASEXP
	}
	// if it contains a pointer or a SEXPREC (what's pointed to by a SEXP), then return an error
	if strings.Contains(typeString, "*") || strings.Contains(typeString, "SEXPREC") {
		return out, NotASEXP
	}

	// first, we want to create a type-unsafe pointer to the underlying data in the RSEXP
	pointerConversion := unsafe.Pointer(r)
	outPoint := (*t)(pointerConversion)

	return *outPoint, err
}
