package main

import (
	"fmt"
	"log"
	"main/app/models"
	"main/app/preprocessing"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func main() {
	f, err := os.Open("samples/weather.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	df := dataframe.ReadCSV(f)
	y := df.Select("Temp")

	for {
		if !preprocessing.DFTest(y) {
			preprocessing.Diff(&y)
		} else {
			break
		}
	}

	AR := 2

	data := preprocessing.MakeAR(y, AR)

	preprocessing.DFTest(y)

	yTensor, xTensor := preprocessing.GetYX(data)

	model, err := models.NewLinearModel(yTensor, xTensor)
	if err != nil {
		fmt.Println(err)
	}

	model.Visualize()
	model.Fit()

	prediction := model.Predict(10)
	fmt.Println(prediction)
}
