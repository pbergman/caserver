package util

type Errors []error

func (e *Errors) Append(err error) {
	*e = append(*e, err)
}

func (e Errors) Error() string {
	var str string
	for i, c := 0, len(e); i < c; i++ {
		if i > 0 {
			str += "\n"
		}
		str += e[i].Error()
	}
	return str
}
