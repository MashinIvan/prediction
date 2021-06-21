package arima

import (
	"github.com/go-gota/gota/dataframe"
	"gorgonia.org/tensor"
)

func PrepareDataARI(df *dataframe.DataFrame, yName string, p int) (yTensor *tensor.Dense, xTensor *tensor.Dense) {
	y := df.Select(yName)

	for {
		if !DFTest(y) {
			Diff(&y)
		} else {
			break
		}
	}

	data := MakeAR(y, p)
	yTensor, xTensor = GetYX(data)

	return yTensor, xTensor
}
