package libwebawk

type Address struct {
	name []string
}

func NewAddress() *Address {
	a := new(Address)
	return a
}

func (a *Address) Insert(v string) {
	a.name = append(a.name, v)
}
