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

	return Eval(program)
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
			if (10 > 1) {
				if (10 > 1) {
					signaldecision notout + out;
				}
				signaldecision 1;
			}`, "unknown operator team: BOOLEAN + BOOLEAN"},
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
