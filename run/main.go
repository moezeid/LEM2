package main

import (
	"fmt"

	"github.com/LEM2"
)

func main() {
	e := &LEM2.Env{

		AttributeMap:        make(map[string][]*LEM2.AttributeObject),
		DecisionMap:         make(map[string][]int),
		AttributeList:       make([]string, 0),
		AttributeValueBlock: make(map[LEM2.Tuple][]int),
	}
	e.Parse()
	i := e.InitialIntersections()
	for i, v := range i{
		i.String()
		fmt.Print(" ",v,"\n")
	}
	t := LEM2.Tuple{
		Attribute: "D",
		Value: "medium",
	}
	r := e.ColumnScan(t,i)
	fmt.Println(r)
}
