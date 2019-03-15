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
            case *object.Array:
                return &object.Integer{Value: int64(len(arg.Elems))}
            default:
                return newError("argument to `len` not supported, got %s", arg.Type())
            }
        },
    },
    "first": &object.Builtin {
        Fn: func(args ...object.Object) object.Object {
            l := len(args)
            if l != 1 {
                return newError("wrong number of arguments. got=%d, want=1", l)
            }

            switch arg := args[0].(type) {
            case *object.Array:
                if len(arg.Elems) == 0 {
                    return NULL
                } else {
                    return arg.Elems[0]
                }
            default:
                return newError("argument to `first` not supported, got %s", arg.Type())
            }
        },
    },
    "last": &object.Builtin {
        Fn: func(args ...object.Object) object.Object {
            l := len(args)
            if l != 1 {
                return newError("wrong number of arguments. got=%d, want=1", l)
            }

            switch arg := args[0].(type) {
            case *object.Array:
                if len(arg.Elems) == 0 {
                    return NULL
                } else {
                    return arg.Elems[len(arg.Elems) - 1]
                }
            default:
                return newError("argument to `last` not supported, got %s", arg.Type())
            }
        },
    },
}
