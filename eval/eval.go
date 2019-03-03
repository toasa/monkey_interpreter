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

    case *ast.ExpressionStatement:
        return Eval(node.Expression)

    case *ast.IntergerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.PrefixExpression:
        right := Eval(node.Right)
        return evalPrefixExpression(node.Operator, right)

    case *ast.Boolean:
        if node.Value {
            return TRUE
        }
        return FALSE
    }

    return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
    for _, stmt := range stmts {
        switch stmt := stmt.(type) {
        case *ast.LetStatement:
        case *ast.ReturnStatement:
        case *ast.ExpressionStatement:
            return Eval(stmt)
        // case *ast BlockStatement:
        //     return evalStatements(stmt.Statements)
        }
    }
    return nil
}

func evalPrefixExpression(op string, right object.Object) object.Object {
    if (op == "!") {
        return evalBangOperatorExpression(right)
    }
    return NULL
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
