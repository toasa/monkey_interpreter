package eval

import (
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node)

    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue)
        if isError(val) {
            return val
        }
        return &object.ReturnValue{Value: val}

    case *ast.ExpressionStatement:
        return Eval(node.Expression)

    case *ast.BlockStatement:
        return evalBlockStatement(node)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.Boolean:
        if node.Value {
            return TRUE
        }
        return FALSE

    case *ast.IfExpression:
        cond := Eval(node.Cond)
        if isError(cond) {
            return cond
        }

        if isTruthly(cond) {
            return Eval(node.Cons)
        } else {
            if node.Alt != nil {
                return Eval(node.Alt)
            }
        }
        return NULL

    case *ast.PrefixExpression:
        right := Eval(node.Right)
        if isError(right) {
            return right
        }
        return evalPrefixExpression(node.Operator, right)

    case *ast.InfixExpression:
        left := Eval(node.Left)
        if isError(left) {
            return left
        }

        right := Eval(node.Right)
        if isError(right) {
            return right
        }

        op := node.Operator
        return evalInfixExpression(op, left, right)
    }

    return nil
}

func evalProgram(program *ast.Program) object.Object {
    var res object.Object

    for _, stmt := range program.Statements {
        res = Eval(stmt)

        switch res := res.(type) {
        case *object.ReturnValue:
            return res.Value
        case *object.Error:
            return res
        }
    }

    return res
}

func evalBlockStatement(bs *ast.BlockStatement) object.Object {
    var res object.Object

    for _, stmt := range bs.Statements {
        res = Eval(stmt)

        if res != nil {
            rt := res.Type()
            if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
                return res
            }
        }
    }
    return res
}

func evalPrefixExpression(op string, right object.Object) object.Object {
    if (op == "!") {
        return evalBangOperatorExpression(right)
    }
    if (op == "-") {
        return evalMinusPrefixOperatorExpression(right)
    }
    return newError("unknown operator: %s%s", op, right.Type())
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(op, left, right)
    case left.Type() != right.Type():
        return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
    default:
        return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
    }
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
    lval := left.(*object.Integer).Value
    rval := right.(*object.Integer).Value
    switch op {
    case "+":
        return &object.Integer{Value: lval + rval}
    case "-":
        return &object.Integer{Value: lval - rval}
    case "*":
        return &object.Integer{Value: lval * rval}
    case "/":
        return &object.Integer{Value: lval / rval}
    case "==":
        return nativeBoolToBooleanObject(lval == rval)
    case "!=":
        return nativeBoolToBooleanObject(lval != rval)
    case "<":
        return nativeBoolToBooleanObject(lval < rval)
    case ">":
        return nativeBoolToBooleanObject(lval > rval)
    default:
        return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
    }
}

func nativeBoolToBooleanObject(b bool) object.Object {
    if b {
        return TRUE
    }
    return FALSE
}

func evalBangOperatorExpression(exp object.Object) object.Object {
    switch exp {
    case TRUE:
        return FALSE
    case FALSE:
        return TRUE
    case NULL:
        return TRUE
    default:
        return FALSE
    }
}

func evalMinusPrefixOperatorExpression(exp object.Object) object.Object {
    i, ok := exp.(*object.Integer)
    if !ok {
        return newError("unknown operator: -%s", exp.Type())
    }

    i.Value = -i.Value
    return i
}

func isTruthly(cond object.Object) bool {

    // monkey言語においてNULLオブジェクトとFALSEオブジェクト以外は、trueとなる
    switch cond {
    case NULL:
        return false
    case FALSE:
        return false
    case TRUE:
        return true
    default:
        return true
    }
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        return (obj.Type() == object.ERROR_OBJ)
    }
    return false
}
