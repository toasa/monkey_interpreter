package eval

import (
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
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

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true", true},
        {"false", false},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testBooleanObject(t, evaledObj, test.expected)
    }
}

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testBooleanObject(t, evaledObj, test.expected)
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
    b, ok := obj.(*object.Boolean)
    if !ok {
        t.Errorf("object is not Boolean, got = %T", obj)
        return false
    }

    if b.Value != expected {
        t.Errorf("object has incorrect value")
        return false
    }

    return true
}
