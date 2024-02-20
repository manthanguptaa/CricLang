package evaluator

import (
	"CricLang/object"
	"fmt"
	"log"
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
	"kohli": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			message := "shaam tak khelenge, inki G phatt jaayegi lekin abhi tera code phatt gaya"
			for _, arg := range args {
				message += " " + arg.Inspect()
			}
			log.Fatalf(message)
			return &object.String{}
		},
	},
	"rohit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newMisfield("mera gale ka vaat lag gaya chilla chilla ke ki 1 argument chahiye! tunne %d de diye", len(args))
			}
			arg, ok := args[0].(*object.String)
			if !ok {
				return newMisfield("mera gale ka vaat lag gaya chilla chilla ke ki sahi type ka argument daal de")
			}
			fmt.Printf("Reporter: %s ke birthday ke baare mei kuch boliye.\n", arg.Value)
			return &object.String{Value: "Rohit: Abhi birthday mei kya bola jata hai? Happy Birthday? Yahi bola jata hai."}
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
