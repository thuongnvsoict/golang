package main

import "fmt"

type greet interface {
	Greet()
}

type Person struct{
	Name string
	Age int
}

// Contructor
func NewPerson(name string, age int) Person{
	return Person{
		Name:name,
		Age:age,
	}
}

func NewPersonPoiter(name string, age int) *Person{
	return &Person{
		Name:name,
		Age:age,
	}
}

//Method
func (p Person) Greet()  {
	fmt.Printf("Hello I am %s\n", p.Name)
}

func main()  {

	p1 := Person{"Thuong", 23}
	p2 := NewPerson("Thuong", 23)

	p1.Greet()
	p2.Greet()
}