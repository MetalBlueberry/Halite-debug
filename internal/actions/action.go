package actions

type Circler interface {
	Circle() (x, y, radius float64)
}

type Liner interface {
	Line() (x1, y1, x2, y2 float64)
}

type Action map[string]interface{}

func (a Action) Method() (name string) {
	return a["Method"].(string)
}

func (a Action) Class() (classes []string) {
	classInterface, ok := a["Class"].([]interface{})
	if !ok {
		panic("Class field is not pressent or type is not []interface{}")
	}
	class := make([]string, len(classInterface))
	for i, item := range classInterface {
		class[i], ok = item.(string)
		if !ok {
			panic("Fields inside class list are not the type string")
		}
	}
	return class
}

func (a Action) Circle() (x, y, radius float64) {
	return a["X"].(float64), a["Y"].(float64), a["R"].(float64)
}

func (a Action) Line() (x1, y1, x2, y2 float64) {
	return a["X1"].(float64), a["Y1"].(float64), a["X2"].(float64), a["Y2"].(float64)
}
