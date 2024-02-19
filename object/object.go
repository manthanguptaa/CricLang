package object

import (
	"CricLang/ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ                     = "INTEGER"
	BOOLEAN_OBJ                     = "BOOLEAN"
	DEAD_BALL_NULL_OBJ              = "DEAD_BALL"
	SIGNALDECISION_RETURN_VALUE_OBJ = "SIGNALDECISION"
	MISFIELD_ERROR_OBJECT           = "MISFIELD"
	FIELD_FUNCTION_OBJECT           = "FIELD"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type DeadBallNull struct{}

func (n *DeadBallNull) Type() ObjectType { return DEAD_BALL_NULL_OBJ }
func (n *DeadBallNull) Inspect() string  { return "deadball" }

type SignalDecisionReturnValue struct {
	Value Object
}

func (sdr *SignalDecisionReturnValue) Type() ObjectType { return SIGNALDECISION_RETURN_VALUE_OBJ }
func (sdr *SignalDecisionReturnValue) Inspect() string  { return sdr.Value.Inspect() }

type Misfield struct {
	Message string
}

func (m *Misfield) Type() ObjectType { return MISFIELD_ERROR_OBJECT }
func (m *Misfield) Inspect() string  { return "MISFIELD: " + m.Message }

type Field struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Field) Type() ObjectType { return FIELD_FUNCTION_OBJECT }
func (f *Field) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("field")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
