package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"log"
	"main/app/estimators"
	"main/app/models"
	"main/app/preprocessing/arima"
	"os"
)

// Uncomment this for working example
//func main() {
//	f, err := os.Open("samples/weather.csv")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	df := dataframe.ReadCSV(f)
//	y := df.Select("Temp")
//
//	for {
//		if !preprocessing.DFTest(y) {
//			preprocessing.Diff(&y)
//		} else {
//			break
//		}
//	}
//
//	AR := 1
//
//	data := preprocessing.MakeAR(y, AR)
//
//	preprocessing.DFTest(y)
//
//	yTensor, xTensor := preprocessing.GetYX(data)
//
//	model, err := models.NewLinearModel(yTensor, xTensor)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	model.Fit()
//
//	prediction := model.Predict(10)
//	fmt.Println(prediction)
//}

// Refactored example
func main() {
	f, err := os.Open("samples/weather.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)
	yTensor, xTensor := arima.PrepareDataARI(&df, "Temp", 1)

	model := models.NewModel(yTensor, xTensor)

	estimator := estimators.NewLSEstimator(model)
	estimator.Visualize("simple_graph")
	estimator.Fit()

	fmt.Println(model.Coefficients)
}
