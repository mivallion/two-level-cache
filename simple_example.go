package main

import (
	"./two_level_cache"
	"bytes"
	"fmt"
)

type Person struct {
	name string
	age  int
}

func (p *Person) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, p.name, p.age)
	return b.Bytes(), nil
}

func (p *Person) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &p.name, &p.age)
	return err
}

func main() {
	c, err := two_level_cache.NewTwoLevelCache(5, 5, &Point{})
	if err != nil {
		panic(err)
	}
	pInp := &Person{name: "Boris", age: 42}
	err = c.Cache("Boris", pInp)
	if err != nil {
		panic(err)
	}
	p, _ := c.Get("Boris")
	err = c.Clear()
	if err != nil {
		panic(err)
	}
}
