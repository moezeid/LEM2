package LEM2






func Inter(s1, s2 []int) (result []int) {
	m := make(map[int]bool)

	for _, val := range s1 {
		m[val] = true
	}

	for _, val := range s2 {
		if _, ok := m[val]; ok {
			result = append(result, val)
		}
	}
	return
}

//maps a decision to a list of a,v pairs
type Intersection map[Tuple][]Tuple

type IntersectionList []Intersection

type initIntersection map[Tuple][]int


func (e *Env) InitialIntersections() map[Tuple][]int {


	//get initial intersections

	l := make(initIntersection)
	for i, v := range e.AttributeValueBlock{
		decisionSet := e.DecisionMap["medium"]
		inter := Inter(v,decisionSet)
		l[i] = inter
	}

	return l
}

func (e *Env) ColumnScan(goal Tuple, initalSets initIntersection) Tuple{

	decisionSet := e.DecisionMap[goal.Value]

	max := -1
	var t Tuple
	for i, v := range initalSets{

		inter := Inter(decisionSet,v)
		if len(inter) > max{
			max = len(inter)
			t = i
		}else if len(inter) == max{
			if len(e.AttributeValueBlock[i]) < len(e.AttributeValueBlock[t]){
				t = i
			}
		}

	}

	return t
}




