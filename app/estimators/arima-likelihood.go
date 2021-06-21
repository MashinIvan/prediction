package estimators

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"io/ioutil"
	"main/app/models"
	"math"
)

type arimaMaxLikelihood struct {
	model *models.Model

	graph        *G.ExprGraph
	coefficients *G.Node

	Likelihood *G.Value
	ErrorEst   *G.Value
	Covariance *G.Value
}

func NewArimaEstimator(model *models.Model, q int) *arimaMaxLikelihood {
	e := &arimaMaxLikelihood{model: model}
	graph := G.NewGraph()

	coefficients := G.NewVector(
		graph,
		G.Float64,
		G.WithName("theta"),
		G.WithShape(e.model.X.Shape()[1]+q),
		G.WithInit(G.Uniform(0, 1)),
	)
	G.Read(coefficients, e.model.Coefficients)

	// TODO
	squaredError := G.UniformRandomNode(model.Graph, G.Float64, 0.0, 1.0)
	covariance := G.UniformRandomNode(model.Graph, G.Float64, 0.0, 1.0)

	a := G.NewScalar(graph, G.Float64, G.WithValue(
		math.Log(math.Pi*2)*float64(model.N)*(-1/2)))
	b := G.Must(G.Mul(
		G.Must(G.Log(covariance)),
		G.NewScalar(graph, G.Float64, G.WithValue(float64(model.N)*(-0.5))),
	))
	c := G.Must(G.Sum(
		G.Must(G.Div(
			squaredError,
			G.Must(G.Mul(
				G.NewScalar(graph, G.Float64, G.WithValue(2.0)),
				covariance,
			)),
		))))

	logLikelihood := G.Must(G.Add(G.Must(G.Add(a, b)), c))
	G.Read(logLikelihood, e.Likelihood)

	e.graph = graph
	return e
}

func (e *arimaMaxLikelihood) Visualize(fileName string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s.dot", fileName), []byte(e.graph.ToDot()), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *arimaMaxLikelihood) Fit() {

}
