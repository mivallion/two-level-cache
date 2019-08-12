package main

import (
	"./two_level_cache"
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type Point struct {
	x int
	y int
}

func (p *Point) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, p.x, p.y)
	return b.Bytes(), nil
}

func (p *Point) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &p.x, &p.y)
	return err
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	c, err := two_level_cache.NewTwoLevelCache(5, 5, &Point{})
	if err != nil {
		panic(err)
	}
	var keys []string
	for i := 0; i < 10; i++ {
		point := &Point{x: rand.Intn(100), y: rand.Intn(100)}
		err := c.Cache(fmt.Sprintf("%v.%v", point.x, point.y), point)
		keys = append(keys, fmt.Sprintf("%v.%v", point.x, point.y))
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < 1000; i++ {
		key := keys[rand.Intn(10)]
		_, ok := c.Get(key)
		if !ok {
			fmt.Printf("%v does not exists", key)
		}
	}
	err = c.Clear()
	if err != nil {
		panic(err)
	}
}
