package repl

import (
    "bufio"
    "fmt"
    "io"
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
)

const PROMPT = ">>> "


// バグあり
// 例えば、`let a 46`と入力した時、panic runtime errorが起こり、replが終了してしまう
func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        // scanが終わるとscannedはfalseになる
        if !scanned {
            return
        }

        line := scanner.Text()
        l := lexer.New(line)
        p := parser.New(l)

        program := p.ParseProgram()
        if len(p.Errors()) != 0 {
            printParserErrors(out, p.Errors())
        }

        io.WriteString(out, program.String())
        io.WriteString(out, "\n")
    }
}

func printParserErrors(out io.Writer, errors []string) {
    for _, msg := range errors {
        io.WriteString(out, "\t" + msg + "\n")
    }
}
