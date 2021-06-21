package estimators

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"io/ioutil"
	"log"
	"main/app/models"
)

type leastSquares struct {
	model *models.Model

	cost         *G.Node
	coefficients *G.Node

	ErrorEst   *G.Value
	Covariance *G.Value
}

func NewLSEstimator(model *models.Model) *leastSquares {
	e := leastSquares{model: model}

	e.coefficients = G.NewVector(
		model.Graph,
		G.Float64,
		G.WithName("coefficients"),
		G.WithShape(e.model.X.Shape()[1]),
		G.WithInit(G.Uniform(0, 1)),
	)
	G.Read(e.coefficients, e.model.Coefficients)

	pred := G.Must(G.Mul(e.model.X, e.coefficients))

	squaredError := G.Must(G.Square(
		G.Must(G.Sub(
			pred, e.model.Y))))
	G.WithName("squaredError")(squaredError)

	e.cost = G.Must(G.Mean(squaredError))
	G.WithName("cost")(e.cost)
	G.Read(e.cost, e.ErrorEst)

	//pInverse := G.Must(G.Inverse(G.Must(G.Mul(G.Must(G.Transpose(e.model.X)), e.model.X))))
	//covariance := G.Must(G.Mul(squaredError, pInverse))
	//G.Read(covariance, e.Covariance)

	return &e
}

func (e *leastSquares) Visualize(fileName string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s.dot", fileName), []byte(e.model.Graph.ToDot()), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *leastSquares) Fit() {
	if _, err := G.Grad(e.cost, e.coefficients); err != nil {
		log.Fatalf("Failed to backpropagate: %v", err)
	}

	machine := G.NewTapeMachine(e.model.Graph, G.BindDualValues(e.coefficients))
	defer machine.Close()

	model := []G.ValueGrad{e.coefficients}
	solver := G.NewVanillaSolver(G.WithLearnRate(0.001))

	//if err := solver.Step(model); err != nil {
	//	log.Fatal(err)
	//}

	iter := 100000
	var err error
	for i := 0; i < iter; i++ {
		if err = machine.RunAll(); err != nil {
			log.Printf("Error during iteration: %v: %v\n", i, err)
			break
		}

		if err = solver.Step(model); err != nil {
			log.Fatal(err)
		}
		machine.Reset()
	}
}
