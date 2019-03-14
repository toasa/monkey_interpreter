package eval

import "monkey_interpreter/object"

var builtins = map[string]*object.Builtin {
    "len": &object.Builtin {
        Fn: func(args ...object.Object) object.Object {
            l := len(args)
            if l != 1 {
                return newError("wrong number of arguments. got=%d, want=1", l)
            }
            switch arg := args[0].(type) {
            case *object.String:
                return &object.Integer{Value: int64(len(arg.Value))}
            default:
                return newError("argument to `len` not supported, got %s", arg.Type())
            }
        },
    },
}
