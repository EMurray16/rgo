//This tests the rgo package
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
	"fmt"
	"github.com/EMurray16/rgo"
)

//export TestFloat
func TestFloat(s C.SEXP) C.SEXP {
	point, err := rgo.NewRSEXP(s)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	TestSlicef, err := rgo.AsNumeric[float64](point)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	//now multiply everything in the slice by two
	var sum float64 = 0
	for i := 0; i < len(TestSlicef); i++ {
		TestSlicef[i] = TestSlicef[i] * 2
		sum += TestSlicef[i]
	}
	TestSlicef = append(TestSlicef, sum)
	// fmt.Println(TestSlicef)

	//now make a new SEXP
	s2 := rgo.NumericToRSEXP(TestSlicef)
	fmt.Println(s2)
	// var out C.SEXP
	out, err := rgo.ExportRSEXP[C.SEXP](s2)
	fmt.Println(out, err)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}
	//defereference the pointer
	return out
}

//export TestInt
func TestInt(s C.SEXP) C.SEXP {
	//make an unsafe pointer to the SEXP
	point, err := rgo.NewRSEXP(s)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	TestSlice, err := rgo.AsNumeric[int](point)
	fmt.Println(TestSlice)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	//now multiply everything in the slice by two
	var Sum int = 0
	for i := 0; i < len(TestSlice); i++ {
		TestSlice[i] = TestSlice[i] * 2
		Sum += TestSlice[i]
	}
	TestSlice = append(TestSlice, Sum)

	//make the new SEXP
	s2 := rgo.NumericToRSEXP(TestSlice)
	out, err := rgo.ExportRSEXP[C.SEXP](s2)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}
	//defereference the pointer
	return out
}

//export TestString
func TestString(s C.SEXP) C.SEXP {
	//make an unsafe pointer to the SEXP
	point, err := rgo.NewRSEXP(s)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	//now get the string and write a file
	TestString, err := rgo.AsCharacter[string](point)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
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
	s2 := rgo.CharacterToRSEXP([][]byte{[]byte(outstring), []byte(outstring + "again")})
	out, err := rgo.ExportRSEXP[C.SEXP](s2)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}
	//defereference the pointer
	return out
}

//export TestMatrix
func TestMatrix(s C.SEXP) C.SEXP {
	vecPoint, err := rgo.NewRSEXP(s)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	TestMat, err := rgo.AsMatrix(vecPoint)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}
	fmt.Println(TestMat)

	//now append the columns to the matrix
	addedcol := make([]float64, TestMat.Nrow)
	for i := 0; i < TestMat.Nrow; i++ {
		row, err := TestMat.GetRow(i)
		if err != nil {
			// return the error as a string
			outString := rgo.CharacterToRSEXP([]string{err.Error()})
			out, _ := rgo.ExportRSEXP[C.SEXP](outString)
			return out
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
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	//now append the columns
	addedrow := make([]float64, TestMat.Ncol)
	for i := 0; i < TestMat.Ncol; i++ {
		col, err := TestMat.GetCol(i)
		if err != nil {
			// return the error as a string
			outString := rgo.CharacterToRSEXP([]string{err.Error()})
			out, _ := rgo.ExportRSEXP[C.SEXP](outString)
			return out
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
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	//convert the matrix back to a slice
	s2 := rgo.MatrixToRSEXP(TestMat)
	out, err := rgo.ExportRSEXP[C.SEXP](s2)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}
	return out
}

//export MakeDataFrame
func MakeDataFrame() C.SEXP {
	column1 := []string{"0", "1", "e", "pi"}
	column2 := []float64{0, 1, 2.71, 3.14}

	col1SEXP := rgo.CharacterToRSEXP(column1)
	col2SEXP := rgo.NumericToRSEXP(column2)

	colNames := []string{"Constant", "Value"}
	rowNames := []string{}

	df, err := rgo.MakeDataFrame(rowNames, colNames, col1SEXP, col2SEXP)
	fmt.Println("made DataFrame:", df)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	out, _ := rgo.ExportRSEXP[C.SEXP](df)

	return out
}

//export MakeNamedList
func MakeNamedList() C.SEXP {
	column1 := []string{"0", "1", "e", "pi"}
	column2 := []float64{0, 1, 2.71, 3.14}

	col1SEXP := rgo.CharacterToRSEXP(column1)
	col2SEXP := rgo.NumericToRSEXP(column2)

	colNames := []string{"Constant", "Value"}

	df, err := rgo.MakeNamedList(colNames, col1SEXP, col2SEXP)
	fmt.Println("made named list:", df)
	if err != nil {
		// return the error as a string
		outString := rgo.CharacterToRSEXP([]string{err.Error()})
		out, _ := rgo.ExportRSEXP[C.SEXP](outString)
		return out
	}

	out, _ := rgo.ExportRSEXP[C.SEXP](df)

	return out
}

func main() {}
