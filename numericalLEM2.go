package LEM2

/*


import (
"fmt"
"strconv"
"strings"
)

func (e *Env) Algorithm(goal Tuple) RuleList {

	e.Parse()

	decisionSet := e.DecisionMap[goal.Value]

	var selectedAttribute Tuple


	mainGoal := e.DecisionMap[goal.Value]
	ruleList := make(RuleList, 0)
	for len(decisionSet) > 0 {
		selectedAttributeList := make([]Tuple, 0)
		testSet := make([]int, 0)
		i := e.InitialIntersections(decisionSet, goal.Value)

		for len(selectedAttributeList) == 0 || !e.IsSubset(testSet, goal.Value) {

			selectedAttribute, i = e.ColumnScan(goal, i, decisionSet)

			if selectedAttribute.Attribute != "" {

				selectedAttributeList = append(selectedAttributeList, selectedAttribute)
				mainGoal = Inter(e.AttributeValueBlock[selectedAttribute], mainGoal)
				//mainGoal = e.reduceDecisionsSet(e.AttributeValueBlock[selectedAttribute],mainGoal)
				i = e.InitialIntersections(mainGoal,goal.Value)
				for index := range selectedAttributeList{
					for v := range i{
						if selectedAttributeList[index] == v{
							i[v] = make([]int,0)
						}
					}
				}

				//if it is an interval, check and see if it can be simplified with other attributes
				//	if e.isInterval(selectedAttribute) {
				//selectedAttribute = e.CheckAndAdjust(selectedAttribute, selectedAttributeList, i)
				//		for v := range i {
				//			if v.Attribute == selectedAttribute.Attribute && e.IntervalContained(selectedAttribute, v) || v == selectedAttribute {
				//				_ = Inter(e.AttributeValueBlock[v], e.AttributeValueBlock[selectedAttribute])
				//				i[v] = make([]int, 0)
				//			}
				//		}
				//	}

				if len(selectedAttributeList) == 1 {
					testSet = e.AttributeValueBlock[selectedAttribute]
				} else if len(selectedAttributeList) > 1 && len(testSet) != 0 {
					testSet = Inter(testSet, e.AttributeValueBlock[selectedAttribute])
				}
				//if !e.isInterval(selectedAttribute) {
				//	for v := range i {
				//		if v.Attribute == selectedAttribute.Attribute {
				//			i[v] = make([]int, 0)
				//		}
				//	}
				//}

			} else {
				if len(testSet) == 0 {
					break
				}
			}
		}
		tmp := e.IsSubset(testSet, goal.Value)
		if len(testSet) > 0 && tmp {
			//oldSet := decisionSet
			casesCovered := Inter(testSet, decisionSet)
			/*
				if len(casesCovered) > 0 {
					decisionSet = e.reduceDecisionsSet(casesCovered, oldSet)


			tupleList := make([]Tuple, len(selectedAttributeList))
			copy(tupleList, selectedAttributeList)

			rule := Rule{
				tupleList,
				goal,
				casesCovered,
			}
			ruleList = append(ruleList, rule)
			mainGoal = e.reduceDecisionsSet(casesCovered,decisionSet)
			decisionSet = e.reduceDecisionsSet(casesCovered,mainGoal)
			//i = e.InitialIntersections(decisionSet, goal.Value)
			/*
					i = e.InitialIntersections(decisionSet, goal.Value)

					selectedAttributeList = make([]Tuple, 0)
					testSet = make([]int, 0)
				} else if len(testSet) == 0 {
					ogSet := e.DecisionMap[goal.Value]
					for i := len(ogSet) - 1; i > 0; {
						decisionSet = append(decisionSet, ogSet[i])
						i--
						if i == len(ogSet)-6 {
							break
						}
					}
					selectedAttributeList = make([]Tuple, 0)
					testSet = make([]int, 0)
					i = e.InitialIntersections(decisionSet, goal.Value)



		}
		//ruleList.String()

	}



	return ruleList

}



func (e *Env) IsSubset(testSet []int, concept string) bool {

	desiredSet := e.DecisionMap[concept]
	if testSet == nil{
		return false
	}
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
	//decisionSet = e.DecisionMap["2"]
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
				if a[i+1].Attribute != "" {
					tmp := Inter(intersection, e.AttributeValueBlock[a[i+1]])
					if len(tmp) != 0 {
						intersection = tmp
					}else{
						a[i+1] = Tuple{}
					}


				}

			}
		}
	}

	return intersection
}

func (e *Env) SimplifyInterval(a1, a2 Tuple) (Tuple, []int) {

	lb := strings.Split(a1.Value, "..")
	num1, _ := strconv.ParseFloat(lb[0], 64)
	num2, _ := strconv.ParseFloat(lb[1], 64)

	lb2 := strings.Split(a2.Value, "..")

	otherNum, _ := strconv.ParseFloat(lb2[0], 64)
	otherNum2, _ := strconv.ParseFloat(lb2[1], 64)

	if num2 == otherNum {
		newTuple := e.SmallerInterval(a1, a2)
		newSet := e.AttributeValueBlock[newTuple]
		e.AttributeValueBlock[newTuple] = newSet

		return newTuple, newSet

	} else {

		var newLow string
		if num1 > otherNum {
			newLow = lb[0]
		} else {
			newLow = lb2[0]
		}
		var newHigh string

		if num2 > otherNum2 {
			newHigh = lb2[1]
		} else {
			newHigh = lb[1]
		}


		if newHigh == newLow  {
			tuple := e.SmallerInterval(a1,a2)

			//newValue := newLow + ".." + newHigh

			newSet := e.AttributeValueBlock[tuple]

			e.AttributeValueBlock[tuple] = newSet
			return tuple,newSet


		}
		newValue := newLow + ".." + newHigh

		newSet := Inter(e.AttributeValueBlock[a1], e.AttributeValueBlock[a2])
		if len(newSet) == 0{
			t := e.SmallerInterval(a1,a2)
			return t,e.AttributeValueBlock[t]
		}

		newTuple := Tuple{
			Attribute: a1.Attribute,
			Value:     newValue,
		}

		e.AttributeValueBlock[newTuple] = newSet

		return newTuple, newSet
	}

}

func (e *Env) RuleCheck(ruleList RuleList, goal string) bool {

	set := make(map[int]bool)
	for _, v := range e.DecisionMap[goal] {
		set[v] = true
	}

	for _, v := range ruleList {

		for _, v1 := range v.CasesCovered {

			if ok := set[v1]; !ok {
				fmt.Printf("Rule covers %d while the decision set does not contain %d", v1, v1)
				return false

			}

		}

	}

	return true
}

func (e *Env) isInterval(att Tuple) bool {
	lb := strings.Split(att.Value, "..")

	_, err := strconv.ParseFloat(lb[0], 64)
	if err == nil {
		return true
	}
	return false

}

func (e *Env) CheckAndAdjust(selectedAttribute Tuple, selectedAttributeList []Tuple, i InitIntersection) Tuple {
	for index := range selectedAttributeList {
		if index < len(selectedAttributeList) && selectedAttributeList[index].Attribute == selectedAttribute.Attribute {
			att := selectedAttributeList[index]
			tuple, list := e.SimplifyInterval(selectedAttribute, att)
			if tuple != selectedAttribute && list != nil{
				for v := range i {
					if v == selectedAttribute{
						i[v] = make([]int,0)
					}
				}
			}
			if list != nil {
				//originalAtt := selectedAttribute
				selectedAttribute = tuple
				length := len(selectedAttributeList)
				selectedAttributeList[index] = selectedAttributeList[length-1]
				selectedAttributeList[length-1] = Tuple{}
				selectedAttributeList = selectedAttributeList[:length-1]
				for v := range i {
					vF := e.GetFirstNum(v)
					oF := e.GetFirstNum(selectedAttribute)
					if vF == oF {

						if t := e.SmallerInterval(v,selectedAttribute); t==selectedAttribute {
							i[v] = make([]int, 0)
						}
					}
				}

				for v := range i{
					if v == att{
						i[v] = make([]int,0)
					}
				}

			} else {
				for v := range i {
					if v == att && e.GetFirstNum(v) <= e.GetFirstNum(selectedAttribute) && e.GetSecondNum(v) <= e.GetSecondNum(selectedAttribute) {
						i[v] = make([]int, 0)
					}
				}
				return e.SmallerInterval(selectedAttribute, att)
			}

		}
	}

	return selectedAttribute
}

func (e *Env) SmallerInterval(t1, t2 Tuple) Tuple {
	lb := strings.Split(t1.Value, "..")
	num1, _ := strconv.ParseFloat(lb[0],64)
	num2, _ := strconv.ParseFloat(lb[1],64)

	lb2 := strings.Split(t2.Value, "..")

	otherNum, _ := strconv.ParseFloat(lb2[0],64)
	otherNum2, _ := strconv.ParseFloat(lb2[1],64)

	diff1 := num2 - num1
	if diff1 < 0 {
		diff1 = -diff1
	}
	diff2 := otherNum2 - otherNum
	if diff2 < 0 {
		diff2 = -diff2
	}

	if diff1 > diff2 {
		return t2
	} else {
		return t1
	}
}

func (e *Env) isFirstNumSame(t1, t2 Tuple) bool {
	lb := strings.Split(t1.Value, "..")
	num1, _ := strconv.ParseFloat(lb[0],64)

	lb2 := strings.Split(t2.Value, "..")

	otherNum, _ := strconv.ParseFloat(lb2[0],64)
	if num1 == otherNum {
		return true
	}
	return false

}

func (e *Env) GetFirstNum(t Tuple) float64 {
	lb := strings.Split(t.Value, "..")
	num1, _ := strconv.ParseFloat(lb[0],64)
	return num1
}
func (e *Env) GetSecondNum(t Tuple) float64 {
	lb := strings.Split(t.Value, "..")
	num2, _ := strconv.ParseFloat(lb[1],64)
	return num2
}


func (e *Env) IntervalContained(t1, t2 Tuple) bool{
	lb := strings.Split(t1.Value, "..")
	i11, _ := strconv.ParseFloat(lb[0],64)
	i12, _ := strconv.ParseFloat(lb[1],64)

	lb2 := strings.Split(t2.Value, "..")

	i21, _ := strconv.ParseFloat(lb2[0],64)
	i22, _ := strconv.ParseFloat(lb2[1],64)



	if i11 <= i21 && i12 <= i22 {
		return true

	}else if i11 >= i21 && i12 >= i22{
		return true
	}
	return false




}
*/
