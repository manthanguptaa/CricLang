package evaluator

import (
	"CricLang/object"
	"fmt"
	"math/rand"
)

var builtins = map[string]*object.Builtin{
	"thala": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newMisfield("girlfriend se raat mei baat kar lena, pehle %d ki jagah 1 argument daal de", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.String{Value: calculateLength(arg)}
			default:
				return newMisfield("girlfriend se raat mei baat kar lena, pehle sahi type ka argument toh daal de")
			}
		},
	},
	"gambhir": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newMisfield("*Gautam Gambhir Stares Angrily* got=%d arguments, want=2 arguments", len(args))
			}
			fmt.Printf("Interviewer: %v or %v\n", args[0].Inspect(), args[1].Inspect())
			return returnRandomValue()
		},
	},
}

func calculateLength(arg object.Object) string {
	len := int(len(arg.(*object.String).Value))
	if len == 7 {
		return fmt.Sprintf("Thala for a reason: %d", len)
	}
	if findDigitSum(len) == 7 {
		return fmt.Sprintf("Thala for a reason: %d", len)
	}
	return fmt.Sprintf("Captain Cool: %d", len)
}

func findDigitSum(num int) int {
	res := 0
	for num > 0 {
		res += num % 10
		num /= 10
	}
	return res
}

func returnRandomValue() object.Object {
	options := []object.Object{
		&object.String{Value: "Gautam Gambhir: baingan"},
		&object.String{Value: "Gautam Gambhir: shaktimaan"},
		&object.String{Value: "Gautam Gambhir: sachin tendulkar"},
		&object.String{Value: "Gautam Gambhir: 23"},
		&object.String{Value: "Gautam Gambhir: spider-man"},
	}
	randomIndex := rand.Intn(len(options))
	pick := options[randomIndex]
	return pick
}
