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
		CasesCovered: 0,
	}



	fmt.Print("Enter the path to the input file: ")

	var input string
/*
	_, err := fmt.Scanln(&input)
	for err != nil{
		fmt.Print("Please enter a valid filepath: ")
		_, err = fmt.Scanln(&input)
	}


 */

	e.FilePath = input
	//fmt.Println(et)
	list := e.Algorithm()
	list.String()
	//e.RuleCheck(list,"no-recurrence-events")

	/*


	r, i1 := e.ColumnScan(t,i)
	fmt.Println(r)
	r,_ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	r, _ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	 */

}
