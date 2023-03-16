package main

type Animal struct {
	ID   string
	Name string
	Type string
}

func (a *Animal) GetKind() string {
	return "animal"
}

func (a *Animal) GetID() string {
	return a.ID
}

func (a *Animal) GetName() string {
	return a.Name
}

func (a *Animal) SetID(id string) {
	a.ID = id
}

func (a *Animal) SetName(name string) {
	a.Name = name
}
