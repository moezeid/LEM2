package LEM2

import (
	"fmt"
)

type Rule struct {
	Attributes []Tuple
	Decision   Tuple
}

type RuleList []Rule

func (r RuleList) String() {

	for i, v := range r {

		fmt.Printf("Rule %d\n", i)

		for j, v1 := range v.Attributes {
			 v1.String()
			if len(v.Attributes) > 1 && j < len(v.Attributes)-1 {
				fmt.Print(" & ")
			}
			if j == len(v.Attributes)-1 {
				fmt.Print(" -----> ")
				v.Decision.String()
				fmt.Print("\n")
			}

		}
	}
}

func Union(s1, s2 []int) []int {
	if s1 == nil || s2 == nil {
		return nil
	}
	m := make(map[int]bool)
	result := make([]int, 0)

	for _, val := range s2 {
			m[val] = true
		}
	for _, val := range s1 {
			if _, ok := m[val]; !ok {
				result = append(result, val)
			}
		}






	if len(result) == 0{
		return s1
	}
	return result
}

func Inter(s1, s2 []int) []int {
	if s1 == nil || s2 == nil {
		return nil
	}
	m := make(map[int]bool)
	result := make([]int, 0)
	for _, val := range s1 {
		m[val] = true
	}

	for _, val := range s2 {
		if _, ok := m[val]; ok {
			result = append(result, val)
		}
	}
	if len(result) == 0{
		fmt.Println("empty intersection")
	}
	return result
}

//maps a decision to a list of a,v pairs
type DecisionToPairs map[Tuple][]Tuple

type IntersectionList []DecisionToPairs

type InitIntersection map[Tuple][]int

func (e *Env) InitialIntersections(decisionSet []int, goal string) InitIntersection {

	//get initial intersections
	l := make(InitIntersection)
	for i, v := range e.AttributeValueBlock {
		inter := Inter(v, decisionSet)
		if inter != nil   {
			l[i] = inter
		}
	}

	return l
}

func (e *Env) ColumnScan(goal Tuple, initialSets InitIntersection, decisionSet []int) (Tuple, InitIntersection) {

	max := 1
	var t Tuple
	var i Tuple
	var v []int
	for i, v = range initialSets {
		if i.Attribute != "" && v != nil && len(v) > 0 {
			inter := v
			if len(inter) > max && len(inter) > 0 {
				max = len(inter)
				t = i
			} else if len(inter) == max {
				if len(e.AttributeValueBlock[i]) < len(e.AttributeValueBlock[t]) {
					if i.Attribute != "" {
						fmt.Printf("%s %d was selected over %s %d \n",i,len(e.AttributeValueBlock[i]),t,len(e.AttributeValueBlock[t]))

						t = i
					}
				}
			}
		}
	}
	if max == 0 {
		fmt.Println("Could not find a set that works!")
		return t, initialSets
	}
/*
	//delete(*initialSets,t)
	for v := range initialSets {
		if v.Attribute != t.Attribute && v.Value != t.Value {
			tmp := make([]int,len(initialSets[v]))
			listCopy := copy(tmp,initialSets[v])
			newSets[v] = []int{listCopy}

		}
	}

 */
/*
	for i := range initialSets {

		if t.Attribute != "" && i.Attribute != t.Attribute && i.Value != t.Value{
			newSets[i] = Inter(e.AttributeValueBlock[i],decisionSet)
		}
	}

 */

	if t.Attribute == "" {
//		panic(errors.New("what the fuck"))
	}
	return t, initialSets
}

func (e *Env) Algorithm(goal Tuple) RuleList {

	ruleList := make(RuleList, 0)
	e.Parse()



	decisionSet := e.DecisionMap[goal.Value]

	var selectedAttribute Tuple
	selectedAttributeList := make([]Tuple, 0)
	testSet := make([]int, 0)
	i := e.InitialIntersections(decisionSet, goal.Value)

	for len(decisionSet) > 0 {
		selectedAttribute, i = e.ColumnScan(goal, i, decisionSet)
		if selectedAttribute.Attribute != "" {
			selectedAttributeList = append(selectedAttributeList, selectedAttribute)
			if len(selectedAttributeList) == 1{
				testSet = e.AttributeValueBlock[selectedAttribute]
			}else{
				testSet = Inter(testSet, e.AttributeValueBlock[selectedAttribute])
			}
			for v := range i {
				if v.Attribute == selectedAttribute.Attribute{
					i[v] = make([]int, 0)
				}
			}


		}

		tmp := e.IsSubset(testSet, goal.Value)
		if len(testSet) > 0 && tmp {
			oldSet := decisionSet
			decisionSet = e.reduceDecisionsSet(Inter(testSet, oldSet), oldSet)
			tupleList := make([]Tuple, len(selectedAttributeList))
			copy(tupleList, selectedAttributeList)

			rule := Rule{
				tupleList,
				goal,
			}
			ruleList = append(ruleList, rule)
			//i = e.InitialIntersections(decisionSet, goal.Value)
			if len(decisionSet) == 0{
				break

			}
			i = e.InitialIntersections(decisionSet,goal.Value)
			for _,v := range selectedAttributeList{
				for t := range i{
					if t == v{
						i[v] = make([]int,0)
					}
				}


			}
			selectedAttributeList = make([]Tuple, 0)
			testSet = make([]int, 0)

			//ruleList.String()

		}

	}
	return ruleList

}

func (e *Env) IsSubset(testSet []int, concept string) bool {

	desiredSet := e.DecisionMap[concept]

	if ok := len(testSet) > len(desiredSet); ok {
		return false
	}
	set := make(map[int]bool)
	for _, v := range desiredSet {
		set[v] = true
	}

	for _, v := range testSet {

		if _, ok := set[v]; !ok {
			return false
		}

	}
	return true
}

func (e *Env) reduceDecisionsSet(intersection, decisionSet []int) []int {


	set := make(map[int]bool)
	newSet := make([]int, 0)
	for _, v := range intersection {
		set[v] = true
	}

	for _, v := range decisionSet {

		if _, ok := set[v]; !ok {
			newSet = append(newSet, v)
		}
	}
	return newSet



}

func (e *Env) IntersectOverList(a []Tuple) []int {

	//	testSet := make([]int,0)
	if len(a) == 0 {
		return make([]int, 0)
	}
	if len(a) == 1 {
		attribute := a[0]
		return e.AttributeValueBlock[attribute]
	}
	intersection := make([]int, 0)

	if len(a) > 1 {
		for i := range a {
			if i == 0 {
				intersection = Inter(e.AttributeValueBlock[a[0]], e.AttributeValueBlock[a[1]])
			} else if i == len(a)-1 {
				break
			} else {
				intersection = Inter(intersection, e.AttributeValueBlock[a[i+1]])

			}
		}
	}

	return intersection
}
/*
func (e *Env) SimplifyInterval(a1, a2 Tuple) (Tuple,[]int){

	lb := strings.Split(a1.Value,"..",)
	fmt.Println(lb)
	num1, _ := strconv.Atoi(lb[0])
	num2, _ := strconv.Atoi(lb[1])

	lb2 := strings.Split(a2.Value,"..")

	otherNum,_ := strconv.Atoi(lb2[0])
	otherNum2,_ := strconv.Atoi(lb2[1])

	var newLow string
	if num1  > otherNum{
		newLow = lb[0]
	}else{
		newLow = lb2[0]
	}
	var newHigh string

	if num2 > otherNum2{
		newHigh = lb2[1]
	}else{
		newHigh = lb[1]
	}

	newValue := newLow + ".." + newHigh

	return Tuple{
		Attribute: a1.Attribute,
		Value: newValue,
	},nil
}

 */

