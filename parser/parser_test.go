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
    checkParserErrors(t, p)
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

func TestBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        val bool
    }{
        {"true;", true},
        {"false;", false},
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

        if !ok {
            t.Fatalf("type assertion error")
        }

        b, ok := stmt.Expression.(*ast.Boolean)

        if !ok {
            t.Fatalf("type assertion to Boolean is error")
        }

        if b.Value != test.val {
            t.Fatalf("incorrect boolean value")
        }

        if b.TokenLiteral() != test.input[:len(test.input)-1] {
            t.Fatalf("incorrect boolean literal")
        }
    }
}

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        val interface{}
    }{
        {"!5;", "!", 5},
        {"-46;", "-", 46},
        {"!true;", "!", true},
        {"!false;", "!", false},
    }

    for _, test := range prefixTests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()
        checkParserErrors(t, p)

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

        if !testLiteralExpression(t, pe.Right, test.val) {
            t.Fatalf("error")
        }
    }
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        lval interface{}
        op string
        rval interface{}
    }{
        {"2 + 3;", 2, "+" , 3},
        {"2 - 3;", 2, "-" , 3},
        {"2 * 3;", 2, "*" , 3},
        {"2 / 3;", 2, "/" , 3},
        {"2 < 3;", 2, "<" , 3},
        {"2 > 3;", 2, ">" , 3},
        {"2 == 3;", 2, "==" , 3},
        {"2 != 3;", 2, "!=" , 3},
        {"true == true;", true, "==", true},
        {"false != true;", false, "!=", true},
        {"false == false;", false, "==", false},
    }

    for _, test := range infixTests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("invalid length of stmt")
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("type assetion to *ast.ExpressionStatement is invaild")
        }

        // テストを整理した
        if !testInfixExpression(t, stmt.Expression, test.lval, test.op, test.rval) {
            t.Fatalf("test of InfixExpression incorrect")
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

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
    id, ok := exp.(*ast.Identifier)
    if !ok {
        t.Fatalf("invalid type assertion of (*ast.Identifier)")
        return false
    }

    if id.Value != val {
        t.Fatalf("incorrect identifier value")
        return false
    }

    if id.TokenLiteral() != val {
        t.Fatalf("incorrect identifier value")
        return false
    }

    return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) bool {
    b, ok := exp.(*ast.Boolean)
    if !ok {
        t.Fatalf("invalid type assertion of (*ast.Boolean)")
        return false
    }

    if b.Value != val {
        t.Fatalf("incorrect boolean value")
        return false
    }

    return true
}

func testLiteralExpression(t * testing.T, exp ast.Expression, expected interface{}) bool {
    switch v:= expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp, v)
    case string:
        return testIdentifier(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
    }
    t.Errorf("type of expected no handled")
    return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
    op string, right interface{}) bool {

    ie, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Fatalf("exp is not (*ast.InfixExpression)")
        return false
    }

    if ie.Operator != op {
        t.Fatalf("incorrect operator")
        return false
    }

    if !testLiteralExpression(t, ie.Left, left) {
        t.Fatalf("incorrect left expression")
        return false
    }

    if !testLiteralExpression(t, ie.Right, right) {
        t.Fatalf("incorrect right expression")
        return false
    }

    return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {
            "-a * b",
            "((-a) * b)",
        },
        {
            "!-a",
            "(!(-a))",
        },
        {
            "a +        b + c",
            "((a + b) + c)",
        },
        {
            "a + b - c",
            "((a + b) - c)",
        },
        {
            "a + b * c",
            "(a + (b * c))",
        },
        {
            "a * b + c",
            "((a * b) + c)",
        },
        {
            "a + b / c",
            "(a + (b / c))",
        },
        {
            "5 > 4 == 3 < 4",
            "((5 > 4) == (3 < 4))",
        },
        {
            "2 + 4 * 5 == 6 + 7 * 8 + 9",
            "((2 + (4 * 5)) == ((6 + (7 * 8)) + 9))",
        },
        {
            "true",
            "true",
        },
        {
            "false",
            "false",
        },
        {
            "3 > 5 == false",
            "((3 > 5) == false)",
        },
        {
            "3 < 5 != true",
            "((3 < 5) != true)",
        },
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)
        if program.String() != test.expected {
            t.Fatalf("expected %s, but got %s", test.expected, program.String())
        }
    }
}
