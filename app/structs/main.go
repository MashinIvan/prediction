package structs

type DataSet struct {
	T []float64
	Y []float64
	X [][]float64
}

type ResponseRow struct {
	PredictedValue float64 `json:"predicted_value"`
	LowerBound     float64 `json:"lower_bound"`
	UpperBound     float64 `json:"upper_bound"`
}

type Model interface {
	Fit(set *DataSet) []float64
	Predict(X [][]float64) []ResponseRow
}
