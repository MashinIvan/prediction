package models

import (
	"fmt"
	"main/app/structs"
	"math"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type linearModel struct {
	g        *G.ExprGraph
	Y        *G.Node
	Pred     *G.Node
	X        *G.Node
	Theta    *G.Node
	Sigma    *G.Node
	ErrorStd *G.Node
}

func NewLinearModel(y *tensor.Dense, x *tensor.Dense) (*linearModel, error) {
	m := linearModel{}
	m.g = G.NewGraph()

	m.Y = G.NodeFromAny(m.g, y, G.WithName("y"))
	m.X = G.NodeFromAny(m.g, x, G.WithName("x"))

	xT, err := G.Transpose(m.X)

	xTx, err := G.Mul(xT, m.X)
	xTy, err := G.Mul(xT, m.Y)

	xTxInv, err := G.Inverse(xTx)
	theta, err := G.Mul(xTxInv, xTy)

	m.Theta = theta
	m.Pred, err = G.Mul(m.X, theta)

	estError, err := G.Sub(m.Y, m.Pred)
	m.ErrorStd, err = G.Mean(G.Must(G.Square(estError)))

	d := G.NewScalar(m.g, G.Float64, G.WithName("degrees of freedom"), G.WithValue(-1.0))
	a, err := G.Div(m.ErrorStd, d)
	m.Sigma, err = G.Mul(xTxInv, a)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *linearModel) Fit() {
	machine := G.NewTapeMachine(m.g)
	defer machine.Close()

	if machine.RunAll() != nil {
		fmt.Println("error running machine")
	}
}

func (m *linearModel) Predict(T int) []structs.ResponseRow {
	theta := m.Theta.Value().Data().([]float64)
	errorStd := m.ErrorStd.Value().Data().(float64)

	multiply := func(x [][]float64, theta []float64) []float64 {
		out := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			for j := 0; j < len(x[0]); j++ {
				out[i] += x[i][j] * theta[i]
			}
		}
		return out
	}

	y := m.Y.Value().Data().([]float64)
	AR := m.Theta.Shape()[0]

	getLastYs := func(y []float64) [][]float64 {
		var x [][]float64
		var row []float64
		for i := 0; i < AR; i++ {
			row = append(row, y[len(y)-i-1])
		}
		x = append(x, row)
		return x
	}
	x := getLastYs(y)

	getIntervalOffset := func(T int) float64 {
		// C = (sum(sum(theta_i)^2k)) * errorStd, i = 0..AR, k = 0..T
		C := 0.0
		for k := 0; k < T; k++ {
			thetaSum := 0.0
			for _, value := range theta {
				thetaSum += value
			}
			C += math.Pow(thetaSum, float64(2*k))
		}
		return math.Pow(C*errorStd, 0.5)
	}

	var response []structs.ResponseRow
	for i := 0; i < T; i++ {
		yPredicted := multiply(x, theta)
		IntervalOffset := getIntervalOffset(i)

		response = append(response, structs.ResponseRow{
			PredictedValue: yPredicted[0],
			LowerBound:     yPredicted[0] - IntervalOffset*1.96,
			UpperBound:     yPredicted[0] + IntervalOffset*1.96,
		})

		y = append(y, yPredicted[0])
		x = getLastYs(y)
	}

	return response
}
