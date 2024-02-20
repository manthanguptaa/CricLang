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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isMisfield(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isMisfield(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isMisfield(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.AppealIfExpression:
		return evalAppealIfExpression(node, env)
	case *ast.SignalDecisionStatement:
		val := Eval(node.SignalDecisionValue, env)
		if isMisfield(val) {
			return val
		}
		return &object.SignalDecisionReturnValue{Value: val}
	case *ast.PlayerStatement:
		val := Eval(node.Value, env)
		if isMisfield(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FieldLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Field{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isMisfield(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isMisfield(args[0]) {
			return args[0]
		}
		return applyField(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.SignalDecisionReturnValue:
			return result.Value
		case *object.Misfield:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

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

func evalAppealIfExpression(ie *ast.AppealIfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isMisfield(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newMisfield("identifier not found: " + node.Value)
	}
	return val
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isMisfield(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyField(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Field)
	if !ok {
		return newMisfield("not a field: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapSignalDecisionValue(evaluated)
}

func extendFunctionEnv(fn *object.Field, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapSignalDecisionValue(obj object.Object) object.Object {
	if signalDecisionValue, ok := obj.(*object.SignalDecisionReturnValue); ok {
		return signalDecisionValue.Value
	}
	return obj
}
