package rgo

import "testing"

func TestCopyMatrix(t *testing.T) {
	startMat := startingMatrix
	endMat := CopyMatrix(startMat)
	if !AreMatricesEqual(startMat, endMat) {
		t.Error("Copy did not produce identical matrix")
	}

	// changing endMat should not change startMat
	endMat.AppendRow([]float64{3.14, 6.28})
	if !AreMatricesEqual(startMat, startingMatrix) {
		t.Errorf("copy resulted in a copy that changes the original too: %v and %v", startMat, startingMatrix)
	}
}

func TestCreateZeros(t *testing.T) {
	// try to get invalid indexes a number of ways
	_, err := CreateZeros(-1, 2)
	if err != InvalidIndex {
		t.Error("expected an invalid index error but didn't get one")
	}
	_, err = CreateZeros(1, -2)
	if err != InvalidIndex {
		t.Error("expected an invalid index error but didn't get one")
	}
	_, err = CreateZeros(-1, -2)
	if err != InvalidIndex {
		t.Error("expected an invalid index error but didn't get one")
	}

	realMat, err := CreateZeros(1, 2)
	checkMat := Matrix{Nrow: 1, Ncol: 2, Data: []float64{0, 0}}
	if err != nil {
		t.Errorf("Got error %v when trying to create valid zeros matrix", err)
	}
	if !AreMatricesEqual(*realMat, checkMat) {
		t.Errorf("expected zeros matrix %v, but got %v instead", checkMat, realMat)
	}

}

func TestCreateIdentity(t *testing.T) {
	checkMat := Matrix{Nrow: 3, Ncol: 3, Data: []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}}
	// there's only 1 invalid index to try
	_, err := CreateIdentity(-3)
	if err != InvalidIndex {
		t.Error("expected an invalid index error but didn't get one")
	}

	realMat, err := CreateIdentity(3)
	if err != nil {
		t.Errorf("Got error %v when trying to create valid identity matrix", err)
	}
	if !AreMatricesEqual(*realMat, checkMat) {
		t.Errorf("expected identity matrix %v, but got %v instead", checkMat, realMat)
	}

}

func TestNewMatrix(t *testing.T) {
	inputSlice := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}
	// get a size error in a few different ways
	_, err := NewMatrix(1, 2, inputSlice)
	if err != ImpossibleMatrix {
		t.Errorf("Expected to get a size mismatch error but got %v instead", err)
	}
	_, err = NewMatrix(3, 3, inputSlice)
	if err != ImpossibleMatrix {
		t.Errorf("Expected to get a size mismatch error but got %v instead", err)
	}

	// try a couple invalid sizes
	_, err = NewMatrix(-1, 3, inputSlice)
	if err != InvalidIndex {
		t.Errorf("Expected to get an invalid index error but got %v instead", err)
	}
	_, err = NewMatrix(3, -3, inputSlice)
	if err != InvalidIndex {
		t.Errorf("Expected to get an invalid index error but got %v instead", err)
	}

	// try it for real now
	realMat, err := NewMatrix(3, 2, inputSlice)
	if err != nil {
		t.Errorf("expceted no error but got %v instead", err)
	}
	if !AreMatricesEqual(*realMat, startingMatrix) {
		t.Errorf("expected to get %v, but got %v instead", startingMatrix, realMat)
	}

	// change the input slice and make sure the matrix didn't change
	inputSlice[3] = 3.14159
	if !AreMatricesEqual(*realMat, startingMatrix) {
		t.Error("changing the slice changed the matrix after creation")
	}
}

func TestMatrix_CreateTranspose(t *testing.T) {
	testMat := startingMatrix.CreateTranspose()
	finalMat := startingTranspose

	if !AreMatricesEqual(*testMat, finalMat) {
		t.Errorf("transpose failed, expected %v, but got %v", finalMat, testMat)
	}
}
