package LEM2

import (
	"fmt"
	"strconv"
	"strings"
)

type Rule struct {
	Attributes      []Tuple
	Decision        Tuple
	CasesCovered    []int
	Specificity     int
	Strength        int
	NumCasesCovered int
}

//maps a decision to a list of a,v pairs
type DecisionToPairs map[Tuple][]Tuple

type IntersectionList []DecisionToPairs

type InitIntersection map[Tuple][]int

type LocalCovering map[string]RuleList

type RuleList []Rule

func (l LocalCovering) String() {

	x := 0
	var decisionName string

	for i, v := range l {
		fmt.Printf("\n")
		if x == 0 {
			for _, j := range v {
				decisionName = j.Decision.Attribute
				break
			}
		}
		fmt.Printf("Ruleset for (%s, %s) \n", decisionName, i)
		v.String()
		x++
	}

}

func (r RuleList) String() {
	for i, v := range r {

		fmt.Printf("Rule %d\n", i)
		fmt.Printf("%d,  %d,  %d \n", v.Specificity, v.Strength, v.NumCasesCovered)

		for j, v1 := range v.Attributes {
			if v1.Value != "" {
				v1.String()
				if len(v.Attributes) > 1 && j < len(v.Attributes)-1 {
					fmt.Print(" & ")
				}
			}
			if j == len(v.Attributes)-1 {
				fmt.Print(" -----> ")
				v.Decision.String()
				fmt.Print("\n")
			}

		}

	}
}

//Inter takes two integer slices and returns a slice of their intersection
func Inter(s1, s2 []int) []int {
	if s1 == nil || s2 == nil {
		return nil
	}
	if len(s1) == 0 || len(s2) == 0 {
		return make([]int, 0)
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

	return result
}

//InitialIntersections takes a decision set and creates the "columns" i.e the intersections
//between each attribute-value block and the current decision set (goal)
func (e *Env) InitialIntersections(decisionSet []int) InitIntersection {

	//get initial intersections
	l := make(InitIntersection)
	for i, v := range e.AttributeValueBlock {
		inter := Inter(v, decisionSet)
		if inter != nil {
			l[i] = inter
		}
	}
	return l
}

//ColumnScan scans the "column" created by InitialIntersections and chooses the attribute-value block
//that has the largest size intersection with the decision set. If there is a tie,
//choose the block who's attribute-value block size is the minimum amongst the ties.  If there is
//another tie, choose the first one.
func (e *Env) ColumnScan(goal Tuple, initialSets InitIntersection, decisionSet []int, selectedAttributeList []Tuple) (Tuple, InitIntersection) {

	//eliminate attributes if they're already selected
	if len(selectedAttributeList) > 0 {
		for index := range selectedAttributeList {
			for v := range initialSets {
				if len(initialSets[v]) > 0 {
					if selectedAttributeList[index].Attribute == v.Attribute && selectedAttributeList[index].Value == v.Value {
						initialSets[v] = make([]int, 0)
					}
				}
			}
		}
	}

	max := 0
	//var t Tuple
	var i Tuple
	var v []int
	//Get the max length
	for i, v = range initialSets {
		if i.Attribute != "" && v != nil && len(v) > 0 {
			inter := v
			if len(inter) > max {
				max = len(inter)
				//t = append(t,i)
			}
		}
	}
	//get the number of sets with that max value if there are multiple initial sets that have that length
	maxCounter := 0
	for i, v = range initialSets {
		if i.Attribute != "" && v != nil && len(v) > 0 {
			inter := v
			if len(inter) == max {
				maxCounter++
			}
		}
	}
	if maxCounter == 1 {
		for i, v = range initialSets {
			if i.Attribute != "" && v != nil && len(v) > 0 {
				inter := v
				if len(inter) == max {
					return i, initialSets
				}
			}
		}

	} else if maxCounter > 1 {

		min := 100000000000
		for i, v = range initialSets {
			if i.Attribute != "" && v != nil && len(v) > 0 {
				inter := v
				//find the Tuple with the max length for its intersection
				// and now test the corresponding (a,v) block's length against the minimum length
				if len(inter) == max && len(e.AttributeValueBlock[i]) < min {
					min = len(e.AttributeValueBlock[i])
				}
			}
		}

		//now that we have the min value, lets grab the block
		for i, v = range initialSets {
			if i.Attribute != "" && v != nil && len(v) > 0 {
				inter := v
				//return the first set whose intersection matches max length and (a,v) block matches min length
				if len(inter) == max && len(e.AttributeValueBlock[i]) == min {
					return i, initialSets
				}
			}
		}

	}
	return Tuple{}, nil

}

//Algorithm runs through the entire decision set and returns a local covering. i.e
// a rule list for each concept in the decision set
func (e *Env) Algorithm() LocalCovering {

	e.Parse()
	localCovering := make(LocalCovering)
	for val, cases := range e.DecisionMap {

		decisionSet := cases
		goal := Tuple{
			Attribute: e.AttributeList[len(e.AttributeList)-1],
			Value:     val,
		}
		var selectedAttribute Tuple

		mainGoal := e.DecisionMap[goal.Value]
		ruleList := make(RuleList, 0)
		for len(decisionSet) > 0 {
			selectedAttributeList := make([]Tuple, 0)
			testSet := make([]int, 0)
			i := e.InitialIntersections(decisionSet)

			for len(selectedAttributeList) == 0 || !e.IsSubset(testSet, goal.Value) {

				selectedAttribute, i = e.ColumnScan(goal, i, decisionSet, selectedAttributeList)

				if selectedAttribute.Attribute != "" {
					selectedAttributeList = append(selectedAttributeList, selectedAttribute)
					mainGoal = Inter(e.AttributeValueBlock[selectedAttribute], mainGoal)
					//mainGoal = e.reduceDecisionsSet(e.AttributeValueBlock[selectedAttribute],mainGoal)
					i = e.InitialIntersections(mainGoal)
					if e.isInterval(selectedAttribute) {
						//remove entries in selected attribute list that match exactly attributes we already have
						for index := range selectedAttributeList {
							if e.isInterval(selectedAttributeList[index]) {
								sA := selectedAttributeList[index]
								for v := range i {
									if sA.Attribute == v.Attribute && (e.GetFirstNum(sA) == e.GetFirstNum(v)) && v != e.SmallerInterval(sA, v) || selectedAttributeList[index] == v {
										i[v] = make([]int, 0)
									} else if sA.Attribute == v.Attribute && e.isInterval(v) {
										if (e.GetSecondNum(sA) == e.GetFirstNum(v)) || (e.GetFirstNum(sA) == e.GetSecondNum(v)) {
											i[v] = make([]int, 0)
										}
									}

								}
							}
						}

					} else {
						for index := range selectedAttributeList {
							for v := range i {
								if selectedAttributeList[index].Attribute == v.Attribute && selectedAttributeList[index].Value == v.Value {
									i[v] = make([]int, 0)
								}
							}
						}
					}

					if len(selectedAttributeList) == 1 {
						testSet = e.AttributeValueBlock[selectedAttribute]
					} else if len(selectedAttributeList) > 1 && len(testSet) != 0 {
						testSet = Inter(testSet, e.AttributeValueBlock[selectedAttribute])
					}

				}
			}
			tmp := e.IsSubset(testSet, goal.Value)
			if len(testSet) > 0 && tmp {
				casesCovered := Inter(testSet, decisionSet)
				tupleList := make([]Tuple, len(selectedAttributeList))
				copy(tupleList, selectedAttributeList)

				for index := 0; index < len(tupleList)-1; index++ {
					if e.isInterval(tupleList[index]) {
						av := Tuple{}
						if tupleList[index].Attribute != "" {
							av = tupleList[index]
							for j := range tupleList {
								if (av.Attribute != "" && tupleList[index].Attribute != "") && (av.Attribute == tupleList[index].Attribute) && av != tupleList[index] && e.IntervalContained(av, tupleList[index]) {
									attribute, _ := e.SimplifyInterval(av, tupleList[j])
									tupleList[index] = Tuple{}
									tupleList[j] = attribute
									newList := make([]Tuple, 0)
									for k := range tupleList {
										if tupleList[k].Attribute != "" {
											newList = append(newList, tupleList[k])
										}
									}
									tupleList = newList

								}
							}

						}
					}
				}

				rule := Rule{
					tupleList,
					goal,
					casesCovered,
					len(tupleList),
					len(casesCovered),
					len(casesCovered),
				}
				ruleList = append(ruleList, rule)
				mainGoal = e.reduceDecisionsSet(casesCovered, decisionSet)
				decisionSet = e.reduceDecisionsSet(casesCovered, mainGoal)

			}

		}
		localCovering[goal.Value] = ruleList

	}
	return localCovering
}

func (e *Env) IsSubset(testSet []int, concept string) bool {

	desiredSet := e.DecisionMap[concept]
	if testSet == nil {
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
					} else {
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

		if newHigh == newLow {
			tuple := e.SmallerInterval(a1, a2)

			//newValue := newLow + ".." + newHigh

			newSet := e.AttributeValueBlock[tuple]

			e.AttributeValueBlock[tuple] = newSet
			return tuple, newSet

		}
		newValue := newLow + ".." + newHigh

		newSet := Inter(e.AttributeValueBlock[a1], e.AttributeValueBlock[a2])
		if len(newSet) == 0 {
			t := e.SmallerInterval(a1, a2)
			return t, e.AttributeValueBlock[t]
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
			if tuple != selectedAttribute && list != nil {
				for v := range i {
					if v == selectedAttribute {
						i[v] = make([]int, 0)
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

						if t := e.SmallerInterval(v, selectedAttribute); t == selectedAttribute {
							i[v] = make([]int, 0)
						}
					}
				}

				for v := range i {
					if v == att {
						i[v] = make([]int, 0)
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
	num1, _ := strconv.ParseFloat(lb[0], 64)
	num2, _ := strconv.ParseFloat(lb[1], 64)

	lb2 := strings.Split(t2.Value, "..")

	otherNum, _ := strconv.ParseFloat(lb2[0], 64)
	otherNum2, _ := strconv.ParseFloat(lb2[1], 64)

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
	num1, _ := strconv.ParseFloat(lb[0], 64)

	lb2 := strings.Split(t2.Value, "..")

	otherNum, _ := strconv.ParseFloat(lb2[0], 64)
	if num1 == otherNum {
		return true
	}
	return false

}

func (e *Env) GetFirstNum(t Tuple) float64 {
	lb := strings.Split(t.Value, "..")
	num1, _ := strconv.ParseFloat(lb[0], 64)
	return num1
}
func (e *Env) GetSecondNum(t Tuple) float64 {
	lb := strings.Split(t.Value, "..")
	num2, _ := strconv.ParseFloat(lb[1], 64)
	return num2
}

func (e *Env) IntervalContained(t1, t2 Tuple) bool {
	lb := strings.Split(t1.Value, "..")
	i11, _ := strconv.ParseFloat(lb[0], 64)
	i12, _ := strconv.ParseFloat(lb[1], 64)

	lb2 := strings.Split(t2.Value, "..")

	i21, _ := strconv.ParseFloat(lb2[0], 64)
	i22, _ := strconv.ParseFloat(lb2[1], 64)

	if i11 <= i21 && i12 <= i22 {
		return true

	} else if i11 >= i21 && i12 >= i22 {
		return true
	}
	return false

}
