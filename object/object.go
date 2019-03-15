package object

import (
    "fmt"
    "bytes"
    "strings"
    "hash/fnv"
    "monkey_interpreter/ast"
)

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
    STRING_OBJ = "STRING"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    FUNCTION_OBJ = "FUNCTION"
    BUILTIN_OBJ = "BUILTIN"
    ARRAY_OBJ = "ARRAY"
    HASH_OBJ = "HASH"
    ERROR_OBJ = "ERROR"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
    return INTEGER_OBJ
}

type String struct {
    Value string
}

func (s *String) Inspect() string {
    return s.Value
}
func (s *String) Type() ObjectType {
    return STRING_OBJ
}

type Function struct {
    Params []*ast.Identifier
    Body *ast.BlockStatement
    Env *Env
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
    var out bytes.Buffer

    params := []string{}
    for _, param := range f.Params {
        params = append(params, param.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n")

    return out.String()
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}
func (rc *ReturnValue) Inspect() string {
    return rc.Value.Inspect()
}

type Null struct {}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}
func (n *Null) Inspect() string {
    return "null"
}

type Error struct {
    Msg string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJ
}
func (e *Error) Inspect() string {
    return "ERROR: " + e.Msg
}

// 引数に0個以上のObject型をとり、返り値にObject型を返す関数
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
    Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType{
    return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
    return "builtin function"
}

type Array struct {
    Elems []Object
}

func (a *Array) Type() ObjectType {
    return ARRAY_OBJ
}
func (a *Array) Inspect() string {
    var out bytes.Buffer

    out.WriteString("[")
    for i, e := range a.Elems {
        out.WriteString(e.Inspect())
        if (i + 1 != len(a.Elems)) {
            out.WriteString(", ")
        }
    }
    out.WriteString("]")

    return out.String()
}

type HashKey struct {
    Type ObjectType
    Value uint64
}

type Hashable interface {
    HashKey() HashKey
}

func (b *Boolean) HashKey() HashKey {
    var val uint64
    if b.Value {
        val = 1
    } else {
        val = 0
    }

    return HashKey{Type: BOOLEAN_OBJ, Value: val}
}
func (i *Integer) HashKey() HashKey {
    return HashKey{Type: INTEGER_OBJ, Value: uint64(i.Value)}
}
func (s *String) HashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))
    return HashKey{Type: STRING_OBJ, Value: h.Sum64()}
}

type HashPair struct {
    Key Object
    Value Object
}

type Hash struct {
    Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
    var out bytes.Buffer
    out.WriteString("{")

    pairs := []string{}
    for _, pair := range h.Pairs {
        pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
    }
    out.WriteString(strings.Join(pairs, ", "))

    out.WriteString("}")

    return out.String()
}
