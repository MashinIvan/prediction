package prediction

type Pipeline interface {
	GetData()
	PrepareData()
	ChooseModel()
	FitModel()
	TuneModel()
	Predict()
}
