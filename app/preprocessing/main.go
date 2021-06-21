package preprocessing

type DataSet struct {
	Columns []Column
	Rows    int64
}

type Column struct {
	Name   string
	Type   string
	Values []interface{}
}
