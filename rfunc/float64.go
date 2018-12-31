package rfunc

import(
	"sort"
)

func FloatMean(Vec []float64) (float64) {
	Sum := float64(0)
	Nvals := 0
	for _,val := range Vec {
		//if val > float64(50) {
			Sum += val
			Nvals++
		//}
	}
	Mean := Sum / float64(Nvals)
	return Mean
}

func FloatMedian(Vec []float64) (float64) {
	//start by sorting it
	sort.Float64s(Vec)
	//Find the indexes for the median
	l := len(Vec)
	var Median float64
	if l <= 1 {
		return Vec[0]
	} else if l%2 == 0 {
		//fmt.Println(l/1, l/2-1, l/2)
		Median = (Vec[l/2 - 1] + Vec[l/2]) / float64(2)
	} else {
		Median = Vec[l/2]
	}
	return Median
}

func FloatMin (InSlice []float64) (float64) {
	MinValue := InSlice[0]
	for _,val := range InSlice {
		if val < MinValue {
			MinValue = val
		}
	}
	return MinValue
}

func FloatMax (InSlice []float64) (float64) {
	MaxValue := InSlice[0]
	for _,val := range InSlice {
		if val > MaxValue {
			MaxValue = val
		}
	}
	return MaxValue
}

func FloatSum (InSlice []float64) (float64) {
	Sum := float64(0)
	for _,val := range InSlice {
		Sum += val
	}
	return Sum
}
