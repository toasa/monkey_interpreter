package eval

import (
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct{
        input string
        expected int64
    }{
        {"5", 5},
        {"46", 46},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testIntegerObject(t, evaledObj, test.expected)
    }
}

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
    i, ok := obj.(*object.Integer)
    if !ok {
        t.Errorf("object is not Integer, got = %T", obj)
        return false
    }

    if i.Value != expected {
        t.Errorf("object has incorrect value")
        return false
    }
    
    return true
}
