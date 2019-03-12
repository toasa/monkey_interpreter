package object

type Env struct {
    store map[string]Object
    outer *Env
}

func NewEnv() *Env {
    s := make(map[string]Object)
    return &Env{store: s}
}

func NewEnclosedEnv(outer *Env) *Env {
    env := NewEnv()
    env.outer = outer
    return env
}

func (e *Env) Get(name string) (Object, bool) {
    obj, ok := e.store[name]
    if !ok && e.outer != nil {
        obj, ok = e.outer.Get(name)
    }
    return obj, ok
}

func (e *Env) Set(name string, obj Object) Object {
    e.store[name] = obj
    return obj
}
