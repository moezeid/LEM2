package main

import (
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


	t := LEM2.Tuple{
		Attribute: "class",
		Value: "no-recurrence-events",
	}
	list := e.Algorithm(t)
	list.String()
	/*


	r, i1 := e.ColumnScan(t,i)
	fmt.Println(r)
	r,_ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	r, _ = e.ColumnScan(t,*i1)
	fmt.Println(r)

	 */

}
