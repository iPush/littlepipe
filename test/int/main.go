package main

import "fmt"

type Person struct {
	name string
	city string
}

var foo = Person{
	"foo", "beijing",
}

var persons []*Person = []*Person{
	&Person{"bar", "shanghai"},
	&Person{"christ", "hongkong"},
}

func Display(payload interface{}) {
	items, ok := payload.([]interface{})
	if !ok {
		fmt.Printf("%+v failed to convertd to slice", payload)
		//return
	}
	for i, item := range items {
		fmt.Printf("%d: %v", i, item)
	}
}

func main() {
	//fmt.Println(persons)
	Display(persons)
}
