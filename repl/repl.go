package repl

import (
    "bufio"
    "fmt"
    "io"
    "monkey_interpreter/token"
    "monkey_interpreter/lexer"
)

const PROMPT = ">>> "

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

        for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
            fmt.Printf("%+v\n", tok)
        }
    }
}
