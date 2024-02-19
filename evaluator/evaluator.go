package evaluator

import (
	"CricLang/ast"
	"CricLang/object"
)

var (
	NOT_OUT   = &object.Boolean{Value: true}
	OUT       = &object.Boolean{Value: false}
	DEAD_BALL = &object.DeadBallNull{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
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
		return DEAD_BALL
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
		return DEAD_BALL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
