package parser

import (
    "testing"
    "fmt"
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

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("not enough lengtn of statement")
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("stmt type assertion invalid")
    }

    ident, ok := stmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("ident type assertion invalid")
    }

    if ident.String() != "foobar" {
        t.Fatalf("false identifier paring")
    }

    if ident.TokenLiteral() != "foobar" {
        t.Fatalf("incorrect TokenLiteral()")
    }

}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "46;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

    if !ok {
        t.Fatalf("stmt type assertion invalid")
    }

    il, ok := stmt.Expression.(*ast.IntergerLiteral)
    if !ok {
        t.Fatalf("type assetion invalid")
    }

    if il.Value != 46 {
        t.Fatalf("incorrect number")
    }

    if il.TokenLiteral() != "46" {
        t.Fatalf("incorrect number")
    }
}

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        integerValue int64
    }{
        {"!5;", "!", 5},
        {"-46;", "-", 46},
    }

    for _, test := range prefixTests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()

        if len(program.Statements) != 1 {
            t.Fatalf("incorrect length of program.Statements")
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

        if !ok {
            t.Fatalf("stmt type assertion invalid")
        }

        pe, ok := stmt.Expression.(*ast.PrefixExpression)

        if !ok {
            t.Fatalf("type assertion invalid")
        }

        if pe.Operator != test.operator {
            t.Fatalf("error")
        }

        if !testIntegerLiteral(t, pe.Right, test.integerValue) {
            t.Fatalf("error")
        }
    }
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) bool {
    il, ok := exp.(*ast.IntergerLiteral)
    if !ok {
        t.Fatalf("invalid type assertion")
        return false
    }

    if il.Value != val {
        t.Fatalf("error")
        return false
    }

    if il.TokenLiteral() != fmt.Sprintf("%d", val) {
        t.Fatalf("error")
        return false
    }

    return true
}
