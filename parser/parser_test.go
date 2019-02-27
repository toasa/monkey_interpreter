package parser

import (
    "testing"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
    input := `
let x = 5;
let y = 10;
let foo = 46;`

    l := lexer.New(input)
    // package がparserのため, parser.go内のNew関数を呼ぶ
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if program == nil {
        t.Fatalf("parseProgram() returns nil")
    }

    tests := []struct {
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foo"},
    }

    for i, tt := range tests {
        stmt := program.Statements[i]
        if !testLetStatement(t, stmt, tt.expectedIdentifier) {
            return
        }
    }
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
    if stmt.TokenLiteral() != "let" {
        return false
    }

    letstmt, ok := stmt.(*ast.LetStatement)
    if !ok {
        return false
    }

    if letstmt.Name.Value != name {
        return false
    }

    if letstmt.Name.TokenLiteral() != name {
        return false
    }

    return true
}

func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()
    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}
