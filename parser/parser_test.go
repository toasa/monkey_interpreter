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

func testReturnStatements(t *testing.T) {
    input := `
return 0;
return 46;`
    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    for _, stmt := range program.Statements {
        rs, ok := stmt.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
            continue
        }
        if rs.TokenLiteral() != "return" {
            t.Errorf("rs.TokenLiteral() not 'return', got %q", rs.TokenLiteral())
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
