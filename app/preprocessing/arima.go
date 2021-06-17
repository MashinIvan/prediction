package preprocessing

import (
	"fmt"
	"main/app/models"
	"math"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/mat"
	"gorgonia.org/tensor"
)

func MakeAR(y dataframe.DataFrame, order int) dataframe.DataFrame {
	var newY series.Series

	var values []float64
	for i := 0; i < y.Nrow(); i++ {
		values = append(values, y.Elem(i, 0).Float())
	}

	newY = series.Floats(values[order:])
	df := dataframe.New(newY)

	for i := order - 1; i >= 0; i-- {
		df = df.CBind(dataframe.New(series.Floats(values[i : len(values)-order+i])))
	}
	return df
}

func Diff(d *dataframe.DataFrame) {
	var df dataframe.DataFrame

	for j := 0; j < d.Ncol(); j++ {
		var values []float64
		for i := 1; i < d.Nrow(); i++ {
			values = append(values, d.Elem(i, j).Float()-d.Elem(i-1, j).Float())
		}
		df = df.CBind(dataframe.New(series.Floats(values)))
	}
	*d = df
}

func Shift(d dataframe.DataFrame) dataframe.DataFrame {
	var values []float64
	for t := 1; t < d.Nrow(); t++ {
		values = append(values, d.Elem(t-1, 0).Float())
	}
	df := dataframe.New(series.Floats(values))
	return df
}

func ToSeries(d dataframe.DataFrame, index int) series.Series {
	var values []float64
	for i := 0; i < d.Nrow(); i++ {
		values = append(values, d.Elem(i, index).Float())
	}
	return series.Floats(values)
}

func DFTest(y dataframe.DataFrame) bool {
	yPrev := Shift(y)
	Diff(&y)
	data := dataframe.New(ToSeries(y, 0), ToSeries(yPrev, 0))

	yTensor, xTensor := GetYX(data)
	model, err := models.NewLinearModel(yTensor, xTensor)
	if err != nil {
		fmt.Println(err)
	}

	model.Fit()
	theta := model.Theta.Value().Data().([]float64)
	sigma := model.Sigma.Value().Data().([]float64)

	return math.Abs(theta[0]/sigma[0]) > 1.96
}

func GetYX(d dataframe.DataFrame) (*tensor.Dense, *tensor.Dense) {
	yTensor := tensor.FromMat64(mat.DenseCopyOf(&matrix{d.Select(0)}))
	xTensor := tensor.FromMat64(mat.DenseCopyOf(&matrix{d.Drop(0)}))
	return yTensor, xTensor
}

type matrix struct {
	dataframe.DataFrame
}

func (m matrix) At(i, j int) float64 {
	return m.Elem(i, j).Float()
}

func (m matrix) T() mat.Matrix {
	return mat.Transpose{Matrix: m}
}
