package main

type Person struct {
	ID   string
	Name string
	Age  int32
}

func (p *Person) GetKind() string {
	return "person"
}

func (p *Person) GetID() string {
	return p.ID
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) SetID(id string) {
	p.ID = id
}

func (p *Person) SetName(name string) {
	p.Name = name
}
