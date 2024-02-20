package evaluator

import (
	"CricLang/lexer"
	"CricLang/object"
	"CricLang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"notout", true},
		{"out", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"notout == notout", true},
		{"out == out", true},
		{"notout == out", false},
		{"notout != out", true},
		{"out != notout", true},
		{"(1 < 2) == notout", true},
		{"(1 < 2) == out", false},
		{"(1 > 2) == notout", false},
		{"(1 > 2) == out", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!notout", false},
		{"!out", true},
		{"!5", false},
		{"!!notout", true},
		{"!!out", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestAppealAppealRejectedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"appeal (notout) { 10 }", 10},
		{"appeal (out) { 10 }", nil},
		{"appeal (1) { 10 }", 10},
		{"appeal (1 < 2) { 10 }", 10},
		{"appeal (1 > 2) { 10 }", nil},
		{"appeal (1 > 2) { 10 } appealrejected { 20 }", 20},
		{"appeal (1 < 2) { 10 } appealrejected { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != DEAD_BALL {
		t.Errorf("object is not null. got=%T(%+v)", obj, obj)
		return false
	}
	return true
}

func TestSignalDecisionReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"signaldecision 10;", 10},
		{"signaldecision 10; 9", 10},
		{"signaldecision 2 * 5; 9", 10},
		{"9; signaldecision 10; 9", 10},
		{`appeal (10 > 1) {
			appeal (10 > 1) {
				signaldecision 10;
			}
			signaldecision 1;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMisfieldHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + notout;", "player type mismatch: INTEGER + BOOLEAN"},
		{"5 + notout; 5;", "player type mismatch: INTEGER + BOOLEAN"},
		{"-notout", "unknown operator team: -BOOLEAN"},
		{"notout + out;", "unknown operator team: BOOLEAN + BOOLEAN"},
		{"5; notout + out; 5", "unknown operator team: BOOLEAN + BOOLEAN"},
		{"appeal (10 > 1) { notout + out; }", "unknown operator team: BOOLEAN + BOOLEAN"},
		{`
			appeal (10 > 1) {
				appeal (10 > 1) {
					signaldecision notout + out;
				}
				signaldecision 1;
			}`, "unknown operator team: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "Hello"`, "unknown operator team: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Misfield)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestPlayerStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"player a = 5; a;", 5},
		{"player a = 5 * 5; a;", 25},
		{"player a = 5; player b = a; b;", 5},
		{"player a = 5; player b = a; player c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFieldObject(t *testing.T) {
	input := "field(x) {x + 2};"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Field)
	if !ok {
		t.Fatalf("object isn't field. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("field has wrong parameters. Paramaters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter isn't 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFieldApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"player identity = field(x) {x;}; identity(5);", 5},
		{"player identity = field(x) {signaldecision x;}; identity(5);", 5},
		{"player double = field(x) {x*2;}; double(5);", 10},
		{"player add = field(x, y) {x + y;}; add(5, 5);", 10},
		{"player add = field(x, y) {x + y;}; add(5+5, add(5, 5));", 20},
		{"field(x) {x;}(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		player newAdder = field(x){
			field(y){x + y;}
		};

		player addTwo = newAdder(2);
		addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object isn't String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// func TestBuiltinFields(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected interface{}
// 	}{
// 		{`thala("")`, &object.String{Value: "Captain Cool: 0"}},
// 		{`thala("four")`, &object.String{Value: "Captain Cool: 4"}},
// 		{`thala("hello world")`, &object.String{Value: "Captain Cool: 11"}},
// 		{`thala(1)`, "girlfriend se raat mei baat kar lena, pehle sahi type ka argument toh daal de"},
// 		{`thala("one", "two")`, "girlfriend se raat mei baat kar lena, pehle 2 ki jagah 1 argument daal de"},
// 		{`thala("seventh")`, &object.String{Value: "Thala for a reason: 7"}},
// 	}

// 	for _, tt := range tests {
// 		evaluated := testEval(tt.input)

// 		misfieldObj, ok := evaluated.(*object.Misfield)
// 		if !ok {
// 			t.Errorf("object is not misfield. got=%T (%+v)", evaluated, evaluated)
// 			continue
// 		}
// 		if misfieldObj.Message != tt.expected.Value {
// 			t.Errorf("wrong message. expected=%q, got=%q", tt.expected, misfieldObj.Message)
// 		}
// 	}
// }
