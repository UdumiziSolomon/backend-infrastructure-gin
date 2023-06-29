package main

import (
	"log"
	"encoding/json"
)

// type Rectangle struct {   // rectangle struct
// 	width, height  float64
// }

// func area(r Rectangle) float64 {  
// 	return r.width * r.height
// }

// // Method
// func (r Rectangle) Area() float64 {
// 	return r.width * r.height
// }
// instance one
// rectOne := Rectangle{1,2}
// rectTwo := Rectangle{30,40}
// fmt.Printf("%v\n%v\n", area(rectOne), area(rectTwo))

// //  instance
// methOne := Rectangle{10,20}
// fmt.Println(methOne.Area())

type Person struct {
	Name	string `json:"name"`
}

func main() {
	p := Person{Name: "solomon"}
	pBytes, err := json.Marshal(p)

	log.Print(err)
	log.Print(p, string(pBytes))
}
