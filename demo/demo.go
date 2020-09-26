//This tests the conversions/rsexp package
package main

/*
#include <Rinternals.h>
// We need to include the shared R headers here
// One way to find this is via the Rgo/sexp directory
// 		For example, my computer all github packages are put in /Go/mod/pkg/github.com/...
// Another way is to find them from your local R installation
// 		- Typical Linux: /usr/share/R/include/
// 		- Typical MacOS: /Library/Frameworks/R.framework/Headers/
// 		- Typical Windows (I think): C:/Program Files/R/R4.0.0/include/
#cgo CFLAGS: -I/Library/Frameworks/R.framework/Headers/ -I/usr/share/R/include/
*/
import "C"

import (
	"rsexp"
)

func deref(g rsexp.GoSEXP) C.SEXP {
	return *(*C.SEXP)(g.Point)
}

//export TestFloat
func TestFloat(s C.SEXP) C.SEXP {
	point, err := rsexp.NewGoSEXP(&s)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	TestSlicef, err := point.AsFloats()
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now multiply everything in the slice by two
	var sum float64 = 0
	for i := 0; i < len(TestSlicef); i++ {
		TestSlicef[i] = TestSlicef[i] * 2
		sum += TestSlicef[i]
	}
	TestSlicef = append(TestSlicef, sum)

	//now make a new SEXP
	s2 := rsexp.Float2sexp(TestSlicef)
	//defereference the pointer
	return deref(s2)
}

//export TestInt
func TestInt(s C.SEXP) C.SEXP {
	//make an unsafe pointer to the SEXP
	point, err := rsexp.NewGoSEXP(&s)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	TestSlice, err := point.AsInts()
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now multiply everything in the slice by two
	var Sum int = 0
	for i := 0; i < len(TestSlice); i++ {
		TestSlice[i] = TestSlice[i] * 2
		Sum += TestSlice[i]
	}
	TestSlice = append(TestSlice, Sum)

	//now make a new SEXP
	s2 := rsexp.Int2sexp(TestSlice)
	//defereference the pointer
	return deref(s2)
}

//export TestString
func TestString(s C.SEXP) C.SEXP {
	//make an unsafe pointer to the SEXP
	point, err := rsexp.NewGoSEXP(&s)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now get the string and write a file
	TestString, err := point.AsStrings()
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now send a message back to R
	outstring := "Gopher says: Thanks for saying '"
	for i, s := range TestString {
		outstring += s
		if i < len(TestString)-1 {
			outstring += " and "
		}
	}
	outstring += "'!"

	//make the new SEXP
	s2 := rsexp.String2sexp([]string{outstring, outstring + "again"})
	//dereference the pointer
	return deref(s2)
}

//export TestMatrix
func TestMatrix(s, dim C.SEXP) C.SEXP {
	vecPoint, err := rsexp.NewGoSEXP(&s)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	sizePoint, err := rsexp.NewGoSEXP(&dim)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	// parse the indices
	inds, err := sizePoint.AsInts()
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now get the matrix
	TestMat, err := vecPoint.AsMatrix(inds[0], inds[1])
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now append the columns to the matrix
	addedcol := make([]float64, TestMat.Nrow)
	for i := 0; i < TestMat.Nrow; i++ {
		row, err := TestMat.GetRow(i)
		if err != nil {
			// return the error as a string
			outString := rsexp.String2sexp([]string{err.Error()})
			return deref(outString)
		}
		var sum float64
		for _, f := range row {
			sum += f
		}
		addedcol[i] = sum
	}
	err = TestMat.AppendCol(addedcol)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//now append the columns
	addedrow := make([]float64, TestMat.Ncol)
	for i := 0; i < TestMat.Ncol; i++ {
		col, err := TestMat.GetCol(i)
		if err != nil {
			// return the error as a string
			outString := rsexp.String2sexp([]string{err.Error()})
			return deref(outString)
		}
		var sum float64
		for _, f := range col {
			sum += f
		}
		addedrow[i] = sum
	}
	err = TestMat.AppendRow(addedrow)
	if err != nil {
		// return the error as a string
		outString := rsexp.String2sexp([]string{err.Error()})
		return deref(outString)
	}

	//convert the matrix back to a slice
	s2 := rsexp.Matrix2sexp(TestMat)
	return deref(s2)
}

func main() {}
