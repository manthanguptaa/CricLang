package evaluator

import (
	"CricLang/ast"
	"CricLang/object"
	"fmt"
)

var (
	NOT_OUT   = &object.Boolean{Value: true}
	OUT       = &object.Boolean{Value: false}
	DEAD_BALL = &object.DeadBallNull{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isMisfield(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isMisfield(left) {
			return left
		}
		right := Eval(node.Right)
		if isMisfield(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.AppealIfExpression:
		return evalAppealIfExpression(node)
	case *ast.SignalDecisionStatement:
		val := Eval(node.SignalDecisionValue)
		if isMisfield(val) {
			return val
		}
		return &object.SignalDecisionReturnValue{Value: val}
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.SignalDecisionReturnValue:
			return result.Value
		case *object.Misfield:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.SIGNALDECISION_RETURN_VALUE_OBJ || rt == object.MISFIELD_ERROR_OBJECT {
				return result
			}
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return NOT_OUT
	}
	return OUT
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newMisfield("unknown player type: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case NOT_OUT:
		return OUT
	case OUT:
		return NOT_OUT
	case DEAD_BALL:
		return NOT_OUT
	default:
		return OUT
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newMisfield("unknown operator team: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newMisfield("player type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newMisfield("unknown operator team: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newMisfield("unknown operator team: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalAppealIfExpression(ie *ast.AppealIfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isMisfield(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return DEAD_BALL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case DEAD_BALL:
		return false
	case NOT_OUT:
		return true
	case OUT:
		return false
	default:
		return true
	}
}

func newMisfield(format string, a ...interface{}) *object.Misfield {
	return &object.Misfield{Message: fmt.Sprintf(format, a...)}
}

func isMisfield(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.MISFIELD_ERROR_OBJECT
	}
	return false
}
