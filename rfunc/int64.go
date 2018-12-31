package rfunc

func Int64Min (InSlice []int64) (int64) {
	MinValue := InSlice[0]
	for _,val := range InSlice {
		if val < MinValue {
			MinValue = val
		}
	}
	return MinValue
}

func Int64Max (InSlice []int64) (int64) {
	MaxValue := InSlice[0]
	for _,val := range InSlice {
		if val > MaxValue {
			MaxValue = val
		}
	}
	return MaxValue
}

func Int64Mean (Vec []int64) (int64) {
	Sum := int64(0)
	Nvals := int64(0)
	for _,val := range Vec {
			Sum += val
			Nvals++
	}
	Mean := Sum / Nvals
	return Mean
}

func Int64Sum (InSlice []int64) (int64) {
	Sum := int64(0)
	for _,val := range InSlice {
		Sum += val
	}
	return Sum
}

func Int64In(refint int64, checkslice []int64) (bool) {
	for _,i := range checkslice {
		if i == refint {
			return true
		}
	}
	return false
}