package main

import (
	"fmt"
)
type person struct {}
func sayHi(p *person) { fmt.Println("hi") }
func (p *person) sayHi() { fmt.Println("hi",p) }
func main() {

	var p *person
	p.sayHi() // hi

	sayHi(p)
}
