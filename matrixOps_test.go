package rgo

import (
	"testing"
)

/* for all these tests we start with the following matrix:
{1.1 4.4
2.2 5.5
3.3 6.6} */
var startingMatrix = Matrix{Nrow: 3, Ncol: 2, Data: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}}

/* This is the transpose of the matrix:
{1.1 2.2 3.3
4.4 5.5 6.6}
*/
var startingTranspose = Matrix{Nrow: 2, Ncol: 3, Data: []float64{1.1, 4.4, 2.2, 5.5, 3.3, 6.6}}
var invalidMatrix = Matrix{Nrow: 3, Ncol: 1, Data: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}}

func TestAreMatricesEqual(t *testing.T) {
	foilMatrix1 := Matrix{Nrow: 3, Ncol: 2, Data: []float64{0, 2.2, 3.3, 4.4, 5.5, 6.6}}
	foilMatrix2 := Matrix{Nrow: 2, Ncol: 2, Data: []float64{0, 2.2, 3.3, 4.4}}
	goodFoil := startingMatrix

	// this tests for same dimensions but different data
	if AreMatricesEqual(foilMatrix1, startingMatrix) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, foilMatrix1)
	}
	// this tests for differing columns
	if AreMatricesEqual(startingMatrix, startingTranspose) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, startingTranspose)
	}
	// this tests for different rows
	if AreMatricesEqual(startingMatrix, foilMatrix2) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, foilMatrix2)
	}

	// finally test for a valid equality
	if !AreMatricesEqual(startingMatrix, goodFoil) {
		t.Errorf("Got unexpected unequal matrices: %v and %v", startingMatrix, goodFoil)
	}
}

func TestAreMatricesEqualTol(t *testing.T) {
	foilMatrix1 := Matrix{Nrow: 3, Ncol: 2, Data: []float64{0, 2.2, 3.3, 4.4, 5.5, 6.6}}
	foilMatrix2 := Matrix{Nrow: 2, Ncol: 2, Data: []float64{1, 2.2, 3.3, 4.4}}
	goodFoil := Matrix{Nrow: 3, Ncol: 2, Data: []float64{1.05, 2.21, 3.30003, 4.4, 5.5, 6.6}}

	// this tests for same dimensions but different data
	if AreMatricesEqualTol(foilMatrix1, startingMatrix, 0.1) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, foilMatrix1)
	}
	// this tests for differing columns
	if AreMatricesEqualTol(startingMatrix, startingTranspose, 0.1) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, startingTranspose)
	}
	// this tests for different rows
	if AreMatricesEqualTol(startingMatrix, foilMatrix2, 0.1) {
		t.Errorf("Got unexpected equal matrices: %v and %v", startingMatrix, foilMatrix2)
	}

	// finally test for a valid equality
	if !AreMatricesEqualTol(startingMatrix, goodFoil, 0.1) {
		t.Errorf("Got unexpected unequal matrices: %v and %v", startingMatrix, goodFoil)
	}
}

func TestMatrix_AppendCol(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	finalMat := Matrix{Nrow: 3, Ncol: 3, Data: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 0, 7.7, 8.8}}
	goodColToAppend := []float64{0, 7.7, 8.8}
	badColToAppend1 := []float64{0, 7.7}           // this is too small
	badColToAppend2 := []float64{0, 7.7, 8.8, 9.9} // this is too big

	// test to make sure we pick up and impossible matrix
	err := invalidMatrix.AppendCol(goodColToAppend)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// test the column that is too small
	err = testMat.AppendCol(badColToAppend1)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test the column that is too big
	err = testMat.AppendCol(badColToAppend2)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test where it works
	err = testMat.AppendCol(goodColToAppend)
	if err != nil {
		t.Error("Expected column append to work but got an error instaed")
	}
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("Expected equal matrices, but they aren't: %v and %v", testMat, finalMat)
	}
}
func TestMatrix_AppendRow(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	finalMat := Matrix{Nrow: 4, Ncol: 2, Data: []float64{1.1, 2.2, 3.3, 7.7, 4.4, 5.5, 6.6, 8.8}}
	goodRowToAppend := []float64{7.7, 8.8}
	badRowToAppend1 := []float64{0}           // this is too small
	badRowToAppend2 := []float64{0, 7.7, 8.8} // this is too big

	// test to make sure we pick up and impossible matrix
	err := invalidMatrix.AppendRow(goodRowToAppend)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// test the column that is too small
	err = testMat.AppendRow(badRowToAppend1)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test the column that is too big
	err = testMat.AppendRow(badRowToAppend2)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test where it works
	err = testMat.AppendRow(goodRowToAppend)
	if err != nil {
		t.Error("Expected column append to work but got an error instaed")
	}
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("Expected equal matrices, but they aren't: %v and %v", testMat, finalMat)
	}

}

func TestMatrix_GetCol(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)

	// test to make sure we pick up and impossible matrix
	_, err := invalidMatrix.GetCol(0)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// try to get a negative column (invalid index)
	dummy, err := testMat.GetCol(-1)
	if err != InvalidIndex {
		t.Errorf("Expected to get Invalid index, but got this result instead: %v and %v", dummy, err)
	}
	// try to get a column that too big (index out of bounds)
	dummy, err = testMat.GetCol(3)
	if err != IndexOutOfBounds {
		t.Errorf("Expected to get Index Out of Bounds, but got this result instead: %v and %v", dummy, err)
	}

	// now get an expected column
	expectedData := []float64{1.1, 2.2, 3.3}
	observedData, err := testMat.GetCol(0)
	if err != nil {
		t.Errorf("expected valid data retreival but got error %v", err)
	}
	if len(expectedData) != len(observedData) {
		t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
	}
	for i, v := range observedData {
		if v != expectedData[i] {
			t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
		}
	}

	// now change the output slice to make sure the matrix didn't change
	observedData[1] = 3.14
	if !AreMatricesEqual(testMat, startingMatrix) {
		t.Error("Changing output slice changed matrix")
	}

	// now get an expected column
	expectedData = []float64{4.4, 5.5, 6.6}
	observedData, err = testMat.GetCol(1)
	if err != nil {
		t.Errorf("expected valid data retreival but got error %v", err)
	}
	if len(expectedData) != len(observedData) {
		t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
	}
	for i, v := range observedData {
		if v != expectedData[i] {
			t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
		}
	}
}
func TestMatrix_GetRow(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)

	// test to make sure we pick up and impossible matrix
	_, err := invalidMatrix.GetRow(0)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// try to get a negative column (invalid index)
	dummy, err := testMat.GetRow(-1)
	if err != InvalidIndex {
		t.Errorf("Expected to get Invalid index, but got this result instead: %v and %v", dummy, err)
	}
	// try to get a column that too big (index out of bounds)
	dummy, err = testMat.GetRow(3)
	if err != IndexOutOfBounds {
		t.Errorf("Expected to get Index Out of Bounds, but got this result instead: %v and %v", dummy, err)
	}

	// now get an expected column
	expectedData := []float64{2.2, 5.5}
	observedData, err := testMat.GetRow(1)
	if err != nil {
		t.Errorf("expected valid data retreival but got error %v", err)
	}
	if len(expectedData) != len(observedData) {
		t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
	}
	for i, v := range observedData {
		if v != expectedData[i] {
			t.Errorf("Expected result and observed result don't match: %v expected, got %v", expectedData, observedData)
		}
	}

	// now change the output slice to make sure the matrix didn't change
	observedData[1] = 3.14
	if !AreMatricesEqual(testMat, startingMatrix) {
		t.Error("Changing output slice changed matrix")
	}
}
func TestMatrix_GetInd(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)

	// test to make sure we pick up and impossible matrix
	_, err := invalidMatrix.GetInd(0, 0)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// try to get an invalid index 3 ways
	f, err := testMat.GetInd(-1, 1)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Invalid Index error but got this instead: %f and %v", f, err)
	}
	f, err = testMat.GetInd(1, -1)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Invalid Index error but got this instead: %f and %v", f, err)
	}
	f, err = testMat.GetInd(-1, -1)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Invalid Index error but got this instead: %f and %v", f, err)
	}

	// try to get an out of bounds index
	f, err = testMat.GetInd(3, 3)
	if err != IndexOutOfBounds {
		t.Errorf("was supposed to get Index Out Of Bounds, but got this instead: %f and %v", f, err)
	}

	// now get a res value
	res, err := testMat.GetInd(2, 1)
	if err != nil {
		t.Errorf("was supposed to get res value but got this error instead: %v", err)
	}
	if res != 6.6 {
		t.Errorf("was supposed to get res value of 6.6 but got %f instead", res)
	}

	// change f and make sure the matrix doesn't change
	res = 12.25
	real2, err := testMat.GetInd(2, 1)
	if err != nil {
		t.Errorf("was supposed to get res value but got this error instead: %v", err)
	}
	if real2 != 6.6 {
		t.Errorf("Changing received value changed matrix")
	}
}

func TestMatrix_SetCol(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	finalMat := Matrix{Nrow: 3, Ncol: 2, Data: []float64{0, 7.7, 8.8, 4.4, 5.5, 6.6}}
	goodColToSet := []float64{0, 7.7, 8.8}
	badColToSet1 := []float64{0, 7.7}           // this is too small
	badColToSet2 := []float64{0, 7.7, 8.8, 9.9} // this is too big

	// test to make sure we pick up and impossible matrix
	err := invalidMatrix.SetCol(0, goodColToSet)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// test the column that is too small
	err = testMat.SetCol(0, badColToSet1)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test the column that is too big
	err = testMat.SetCol(0, badColToSet2)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test with an invalid index
	err = testMat.SetCol(-1, goodColToSet)
	if err != InvalidIndex {
		t.Error("expected to get an invalid index error but didn't get one")
	}
	// test with an out of bounds index
	err = testMat.SetCol(2, goodColToSet)
	if err != IndexOutOfBounds {
		t.Error("expected to get an index out of bounds error but didn't get one")
	}

	// test where it works
	err = testMat.SetCol(0, goodColToSet)
	if err != nil {
		t.Error("Expected column append to work but got an error instaed")
	}
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("Expected equal matrices, but they aren't: %v and %v", testMat, finalMat)
	}

	// change the inserted slice and make sure the matrix doesn't change
	goodColToSet[1] = 3.14
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("matrix changed after setting column when it shouldn't have")
	}
}
func TestMatrix_SetRow(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	finalMat := Matrix{Nrow: 3, Ncol: 2, Data: []float64{1.1, 2.2, 7.7, 4.4, 5.5, 8.8}}
	goodRowToSet := []float64{7.7, 8.8}
	badRowToSet1 := []float64{7.7}           // this is too small
	badRowToSet2 := []float64{7.7, 8.8, 9.9} // this is too big

	// test to make sure we pick up and impossible matrix
	err := invalidMatrix.SetRow(0, goodRowToSet)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// test the column that is too small
	err = testMat.SetRow(0, badRowToSet1)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test the column that is too big
	err = testMat.SetRow(0, badRowToSet2)
	if err != SizeMismatch {
		t.Error("Expected a size mismatch but didn't get one")
	}

	// test with an invalid index
	err = testMat.SetRow(-1, goodRowToSet)
	if err != InvalidIndex {
		t.Error("expected to get an invalid index error but didn't get one")
	}
	// test with an out of bounds index
	err = testMat.SetRow(3, goodRowToSet)
	if err != IndexOutOfBounds {
		t.Error("expected to get an index out of bounds error but didn't get one")
	}

	// test where it works
	err = testMat.SetRow(2, goodRowToSet)
	if err != nil {
		t.Error("Expected column append to work but got an error instaed")
	}
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("Expected equal matrices, but they aren't: %v and %v", testMat, finalMat)
	}

	// change the inserted slice and make sure the matrix doesn't change
	goodRowToSet[1] = 3.14
	if AreMatricesEqual(testMat, startingMatrix) {
		t.Errorf("matrix changed after setting column when it shouldn't have")
	}
}
func TestMatrix_SetInd(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	finalMat := Matrix{Nrow: 3, Ncol: 2, Data: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 3.14}}

	// test to make sure we pick up and impossible matrix
	err := invalidMatrix.SetInd(0, 0, 3.14)
	if err != ImpossibleMatrix {
		t.Error("Expected impossible matrix error but didn't get one")
	}

	// try to get an invalid index 3 ways
	err = testMat.SetInd(-1, 1, 3.14)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Index Out Of Bounds, but got this instead: %v", err)
	}
	err = testMat.SetInd(1, -1, 3.14)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Index Out Of Bounds, but got this instead: %v", err)
	}
	err = testMat.SetInd(-1, -1, 3.14)
	if err != InvalidIndex {
		t.Errorf("was supposed to get Index Out Of Bounds, but got this instead: %v", err)
	}

	// try to get an out of bounds index
	err = testMat.SetInd(3, 3, 3.14)
	if err != IndexOutOfBounds {
		t.Errorf("was supposed to get Index Out Of Bounds, but got this instead: %v", err)
	}

	// now get a real value
	err = testMat.SetInd(2, 1, 3.14)
	if err != nil {
		t.Errorf("was supposed to get real value but got this error instead: %v", err)
	}
	if !AreMatricesEqual(testMat, finalMat) {
		t.Errorf("did not set index correctly. got: %v, expected: %v", testMat, finalMat)
	}
}
