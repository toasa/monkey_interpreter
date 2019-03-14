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

func Eval(node ast.Node, env *object.Env) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env)

    case *ast.LetStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)

    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue, env)
        if isError(val) {
            return val
        }
        return &object.ReturnValue{Value: val}

    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)

    case *ast.BlockStatement:
        return evalBlockStatement(node, env)

    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}

    case *ast.StringLiteral:
        return &object.String{Value: node.Value}

    case *ast.Identifier:

        if val, ok := env.Get(node.Value); ok {
            return val
        }

        if b, ok := builtins[node.Value]; ok {
            return b
        }

        return newError("identifier not found: " + node.Value)

    case *ast.FunctionLiteral:
        return &object.Function{Params: node.Params, Body: node.Body, Env: env}

    case *ast.FunctionCall:
        f := Eval(node.Func, env)
        if isError(f) {
            return f
        }
        args := evalExpressions(node.Args, env)
        if len(args) == 1 && isError(args[0]) {
            return args[0]
        }

        return applyFunction(f, args)

    case *ast.Boolean:
        if node.Value {
            return TRUE
        }
        return FALSE

    case *ast.IfExpression:
        cond := Eval(node.Cond, env)
        if isError(cond) {
            return cond
        }

        if isTruthly(cond) {
            return Eval(node.Cons, env)
        } else {
            if node.Alt != nil {
                return Eval(node.Alt, env)
            }
        }
        return NULL

    case *ast.ArrayLiteral:
        a := &object.Array{}
        elems := evalExpressions(node.Elems, env)
        a.Elems = elems
        return a

    case *ast.PrefixExpression:
        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }
        return evalPrefixExpression(node.Operator, right)

    case *ast.InfixExpression:
        left := Eval(node.Left, env)
        if isError(left) {
            return left
        }

        right := Eval(node.Right, env)
        if isError(right) {
            return right
        }

        op := node.Operator
        return evalInfixExpression(op, left, right)
    }

    return nil
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
    var res object.Object

    for _, stmt := range program.Statements {
        res = Eval(stmt, env)

        switch res := res.(type) {
        case *object.ReturnValue:
            return res.Value
        case *object.Error:
            return res
        }
    }

    return res
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Env) object.Object {
    var res object.Object

    for _, stmt := range bs.Statements {
        res = Eval(stmt, env)

        if res != nil {
            rt := res.Type()
            if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
                return res
            }
        }
    }
    return res
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
    args := []object.Object{}
    for _, exp := range exps {
        evaled := Eval(exp, env)
        if isError(evaled) {
            return []object.Object{evaled}
        }
        args = append(args, evaled)
    }
    return args
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
    case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
        return evalStringInfixExpression(op, left, right)
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

func evalStringInfixExpression(op string, left, right object.Object) object.Object {
    lStr := left.(*object.String).Value
    rStr := right.(*object.String).Value

    switch op {
    case "+":
        return &object.String{Value: lStr + rStr}
    case "==":
        return nativeBoolToBooleanObject(lStr == rStr)
    case "!=":
        return nativeBoolToBooleanObject(lStr != rStr)
    default:
        return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
    }
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

func applyFunction(fn object.Object, args []object.Object) object.Object {

    switch fn := fn.(type) {
    case *object.Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaled := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaled)
    case *object.Builtin:
        return fn.Fn(args...)
    default:
        return newError("not a function: %s", fn.Type())
    }
}

func extendFunctionEnv(f *object.Function, args []object.Object) *object.Env {
    // f.Envに包まれた新しい環境envを生成
    env := object.NewEnclosedEnv(f.Env)

    // 新しい環境envに引数の名前と値をsetする
    for i, param := range f.Params {
        env.Set(param.Value, args[i])
    }

    // f.Envに包まれた小さな環境envを返す
    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if rv, ok := obj.(*object.ReturnValue); ok {
        return rv.Value
    }

    return obj
}

func nativeBoolToBooleanObject(b bool) object.Object {
    if b {
        return TRUE
    }
    return FALSE
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
