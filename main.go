package main

import (
    "fmt"
    "os"
    "os/user"
    "io/ioutil"
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "monkey_interpreter/eval"
    "monkey_interpreter/repl"
)

func main() {

    user, err := user.Current()
    if err != nil {
        panic(err)
    }

    if len(os.Args) == 1 {
        fmt.Printf("howdy? %s\n", user.Username)
        repl.Start(os.Stdin, os.Stdout)
    } else if len(os.Args) == 2 {
        exec(os.Args[1])
    }
}

func exec(filename string) {
    f, err := os.Open("./examples/" + filename)
    if err != nil{
        fmt.Println("error")
    }
    defer f.Close()

    b, err := ioutil.ReadAll(f)
    input := string(b)

    l := lexer.New(input)
    p := parser.New(l)

    program := p.ParseProgram()

    env := object.NewEnv()
    eval.Eval(program, env)
}
