package models

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type Preprocessing interface {
	PrepareData()
}

type LossFunction interface {
	Visualize()
	Fit()
}

type Model struct {
	Graph *G.ExprGraph

	Y *G.Node
	X *G.Node
	N int

	Loss         LossFunction
	Coefficients *G.Value
}

func NewModel(Y, X *tensor.Dense) *Model {
	m := Model{}
	m.Graph = G.NewGraph()

	err := Y.Reshape(Y.Shape()[0])
	if err != nil {
		fmt.Println("Reshape error", err)
	}
	m.Y = G.NodeFromAny(m.Graph, Y, G.WithName("y"))
	m.X = G.NodeFromAny(m.Graph, X, G.WithName("x"))
	m.N = X.Shape()[0]

	return &m
}
