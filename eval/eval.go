package eval

import (
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
        return evalStatements(node.Statements)

    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue)
        return &object.ReturnValue{Value: val}

    case *ast.ExpressionStatement:
        return Eval(node.Expression)

    case *ast.BlockStatement:
        return evalStatements(node.Statements)

    case *ast.IntergerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.Boolean:
        if node.Value {
            return TRUE
        }
        return FALSE

    case *ast.IfExpression:
        cond := Eval(node.Cond)

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
        return evalPrefixExpression(node.Operator, right)

    case *ast.InfixExpression:
        left := Eval(node.Left)
        right := Eval(node.Right)
        op := node.Operator
        return evalInfixExpression(op, left, right)
    }

    return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
    var res object.Object

    for _, stmt := range stmts {
        res = Eval(stmt)
        rv, ok := res.(*object.ReturnValue)
        if ok {
            return rv.Value
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
    return NULL
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
    lval := left.(*object.Integer).Value
    rval := right.(*object.Integer).Value
    if (op == "+" || op == "-" || op == "*" || op == "/") {
        var val int64
        if (op == "+") {
            val = lval + rval
        } else if (op == "-") {
            val = lval - rval
        } else if (op == "*") {
            val = lval * rval
        } else if (op == "/") {
            val = lval / rval
        }

        return &object.Integer{Value: val}
    }

    var res bool
    if (op == "==") {
        res = (lval == rval)
    } else if (op == "!=") {
        res = (lval != rval)
    } else if (op == "<") {
        res = (lval < rval)
    } else if (op == ">") {
        res = (lval > rval)
    }

    if res {
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
        return NULL
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
