package rgo

import "testing"

func TestMatrix_AddConstant(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	addedMat := Matrix{Nrow: 3, Ncol: 2, Data: []float64{2.2, 3.3, 4.4, 5.5, 6.6, 7.7}}

	// add a positive number
	testMat.AddConstant(1.1)
	if !AreMatricesEqualTol(addedMat, testMat, 1e-15) {
		t.Errorf("failed on adding positive constant. Expected %v but got %v", addedMat, testMat)
	}

	// add a negative number
	testMat.AddConstant(-1.1)
	if !AreMatricesEqualTol(startingMatrix, testMat, 1e-15) {
		t.Errorf("failed on adding negative constant. Expected %v but got %v", startingMatrix, testMat)
	}

	// add 0
	testMat.AddConstant(0)
	if !AreMatricesEqualTol(startingMatrix, testMat, 1e-15) {
		t.Errorf("failed on adding 0. Expected %v but got %v", startingMatrix, testMat)
	}

}

func TestMatrix_MultiplyConstant(t *testing.T) {
	testMat := CopyMatrix(startingMatrix)
	multMat1 := Matrix{Nrow: 3, Ncol: 2, Data: []float64{2.75, 5.5, 8.25, 11, 13.75, 16.5}}
	multMat2 := Matrix{Nrow: 3, Ncol: 2, Data: []float64{-5.5, -11, -16.5, -22, -27.5, -33}}

	// multiply by a positive number
	testMat.MultiplyConstant(2.5)
	if !AreMatricesEqualTol(multMat1, testMat, 1e-14) {
		t.Errorf("failed on adding positive constant. Expected %v but got %v", multMat1, testMat)
	}

	// multiply by a negative number
	testMat.MultiplyConstant(-2)
	if !AreMatricesEqualTol(multMat2, testMat, 1e-14) {
		t.Errorf("failed on adding negative constant. Expected %v but got %v", multMat2, testMat)
	}

	// multiply by fraction (division)
	testMat.MultiplyConstant(-0.2)
	if !AreMatricesEqualTol(startingMatrix, testMat, 1e-13) {
		t.Errorf("failed on adding 0. Expected %v but got %v", startingMatrix, testMat)
	}

}

func TestMatrixAdd(t *testing.T) {
	addingMat := &Matrix{Nrow: 3, Ncol: 2, Data: []float64{1, 1, 1, 2, 2, 2}}
	addedMat := Matrix{Nrow: 3, Ncol: 2, Data: []float64{2.1, 3.2, 4.3, 6.4, 7.5, 8.6}}

	// try to add to an impossible matrix
	_, err := MatrixAdd(&invalidMatrix, addingMat)
	if err != ImpossibleMatrix {
		t.Error("expected to get an impossible matrix error but didn't")
	}

	// try to add incompatible matrices
	_, err = MatrixAdd(&startingTranspose, addingMat)
	if err != SizeMismatch {
		t.Error("expected to get a size mismatch error but got this instead:", err)
	}

	// try a real addition
	testMat, err := MatrixAdd(addingMat, &startingMatrix)
	if err != nil {
		t.Error("got unexpected error:", err)
	}

	if !AreMatricesEqualTol(*testMat, addedMat, 1e-15) {
		t.Errorf("matrix addition didn't work. Expected %v but got %v instead", addedMat, *testMat)
	}
}

func TestMatrixMultiply(t *testing.T) {
	identity, err := CreateIdentity(2)
	if err != nil {
		t.Error("got error creating identity matrix:", err)
	}

	// try to add to an impossible matrix
	_, err = MatrixMultiply(&invalidMatrix, &startingMatrix)
	if err != ImpossibleMatrix {
		t.Error("expected to get an impossible matrix error but didn't")
	}

	// try to add incompatible matrices
	_, err = MatrixMultiply(&startingMatrix, &startingMatrix)
	if err != SizeMismatch {
		t.Error("expected to get a size mismatch error but got this instead:", err)
	}

	// multiplying by the identity matrix should result in an identical matrix
	testMat2, err := MatrixMultiply(&startingMatrix, identity)
	if err != nil {
		t.Error("Got error when multiplying by identity matrix:", err)
	}
	if !(AreMatricesEqualTol(startingMatrix, *testMat2, 1e-15)) {
		t.Errorf("Multiplying matrix %v by identity matrix did not result in itself but some other matrix %v", startingMatrix, testMat2)
	}

	// now do a "real" multiplication
	multBy := &Matrix{Nrow: 2, Ncol: 2, Data: []float64{1.1, 1.1, 2.2, 2.2}}
	/* SHOW YOUR WORK
	we multiply the starting matrix by multBy, resulting in a 3 x 2 matrix
	[1,1] = 1.1 * 1.1 + 4.4 * 1.1 = 6.05
	[2,1] = 2.2 * 1.1 + 5.5 * 1.1 = 8.47
	[3,1] = 3.3 * 1.1 + 6.6 * 1.1 = 10.89
	[1,2] = 1.1 * 2.2 + 4.4 * 2.2 = 12.1
	[2,2] = 2.2 * 2.2 + 5.5 * 2.2 = 16.94
	[3,2] = 3.3 * 2.2 + 6.6 * 2.2 = 21.78
	*/
	multRes := Matrix{Nrow: 3, Ncol: 2, Data: []float64{6.05, 8.47, 10.89, 12.1, 16.94, 21.78}}
	testMat3, err := MatrixMultiply(&startingMatrix, multBy)
	if err != nil {
		t.Error("got an unexpected error:", err)
	}
	if !AreMatricesEqualTol(*testMat3, multRes, 1e-13) {
		t.Errorf("Expected to get matrix product %v but got %v instead", multRes, testMat3)
	}
}
