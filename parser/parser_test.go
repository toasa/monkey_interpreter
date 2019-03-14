package parser

import (
    "testing"
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
    tests := []struct {
        input string
        expectedIdentifier string
        expectedValue interface{}
    }{
        {"let x = 5;", "x", 5},
        {"let y = true;", "y", true},
        {"let foo = y;", "foo", "y"},
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        // package がparserのため, parser.go内のNew関数を呼ぶ
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if program == nil {
            t.Fatalf("parseProgram() returns nil")
        }

        testLetStatement(t, program.Statements[0], test.expectedIdentifier)

        ls := program.Statements[0].(*ast.LetStatement)
        testLiteralExpression(t, ls.Value, test.expectedValue)
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct {
        input string
        expectedVal interface{}
    }{
        {"return 5;", 5},
        {"return true;", true},
        {"return y;", "y"},
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        // package がparserのため, parser.go内のNew関数を呼ぶ
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if program == nil {
            t.Fatalf("parseProgram() returns nil")
        }

        rs, ok := program.Statements[0].(*ast.ReturnStatement)
        if !ok {
            t.Fatalf("type assertion to *ast.ReturnStatement is incorrect")
        }

        testLiteralExpression(t, rs.ReturnValue, test.expectedVal)
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

    il, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("type assertion invalid")
    }

    if il.Value != 46 {
        t.Fatalf("incorrect number")
    }

    if il.TokenLiteral() != "46" {
        t.Fatalf("incorrect number")
    }
}

func TestStringLiteralExpression(t *testing.T) {
    input := `"hello world";`

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    stmt := program.Statements[0].(*ast.ExpressionStatement)
    sl, ok := stmt.Expression.(*ast.StringLiteral)
    if !ok {
        t.Fatalf("incorrect type assertion")
    }

    if sl.Value != "hello world" {
        t.Errorf("expected hello world, but got %s", sl.Value)
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
            t.Fatalf("type assertion to *ast.ExpressionStatement is invaild")
        }

        // テストを整理した
        if !testInfixExpression(t, stmt.Expression, test.lval, test.op, test.rval) {
            t.Fatalf("test of InfixExpression incorrect")
        }
    }
}

func TestIfExpression(t *testing.T) {
    input := "if (x < y) { x }"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("length of parsed program.Statements is invalid")
    }

    es, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    ie, ok := es.Expression.(*ast.IfExpression)
    if ie.TokenLiteral() != "if" {
        t.Fatalf("token literal error")
    }

    if !testInfixExpression(t, ie.Cond, "x", "<", "y") {
        t.Fatalf("ie.Cond invalid form")
    }

    cons, ok := ie.Cons.Statements[0].(*ast.ExpressionStatement)

    if !ok {
        t.Fatalf("type assertion error")
    }

    if !testIdentifier(t, cons.Expression, "x") {
        t.Fatalf("consequence parsing failure")
    }

    if ie.Alt != nil {
        t.Errorf("error")
    }
}

func TestIfElseExpression(t *testing.T) {
    input := "if (x < y) { x } else { y }"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("length of parsed program.Statements is invalid")
    }

    es, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    ie, ok := es.Expression.(*ast.IfExpression)
    if ie.TokenLiteral() != "if" {
        t.Fatalf("token literal error")
    }

    if !testInfixExpression(t, ie.Cond, "x", "<", "y") {
        t.Fatalf("ie.Cond invalid form")
    }

    cons, ok := ie.Cons.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    if !testIdentifier(t, cons.Expression, "x") {
        t.Fatalf("parsing of consequence statements is failure")
    }

    alt, ok := ie.Alt.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    if !testIdentifier(t, alt.Expression, "y") {
        t.Fatalf("parsing of alternative statements is failure")
    }
}

func TestFunctionLiteralParsing(t *testing.T) {
    input := "fn(x, y) { x + y; }"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    fl, ok := stmt.Expression.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("type assertion error")
    }

    param0 := fl.Params[0]
    if !testIdentifier(t, param0, "x") {
        t.Fatalf("parsing parameter error")
    }

    param1 := fl.Params[1]
    if !testIdentifier(t, param1, "y") {
        t.Fatalf("parsing parameter error")
    }

    es, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    testInfixExpression(t, es.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
    tests := []struct {
        input string
        expectParams []string
    }{
        {
            input: "fn() {};",
            expectParams: []string{},
        },
        {
            input: "fn(x) {};",
            expectParams: []string{"x"},
        },
        {
            input: "fn(x, y, z) {};",
            expectParams: []string{"x", "y", "z"},
        },
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        es, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("type assertion error")
        }

        fl, ok := es.Expression.(*ast.FunctionLiteral)
        if !ok {
            t.Fatalf("type assertion error")
        }

        if len(fl.Params) != len(test.expectParams) {
            t.Fatalf("params num incorrect")
        }

        for i, ident := range test.expectParams {
            testLiteralExpression(t, fl.Params[i] ,ident)
        }
    }
}

func TestFunctionCallParsing(t *testing.T) {
    //input := "fn(x, y) { x + y; }(3, 6)"
    input := "add(1, 2 * 3, 4 + 5)"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    es, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("type assertion error")
    }

    fc, ok := es.Expression.(*ast.FunctionCall)
    if !ok {
        t.Fatalf("type assertion error")
    }

    testIdentifier(t, fc.Func, "add")

    testIntegerLiteral(t, fc.Args[0], 1)
    testInfixExpression(t, fc.Args[1], 2, "*", 3)
    testInfixExpression(t, fc.Args[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
    input := "[1, 2 * 2, 3 + 3]"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    al, ok := stmt.Expression.(*ast.ArrayLiteral)
    if !ok {
        t.Fatalf("invalid assertion error")
    }

    if len(al.Elems) != 3 {
        t.Fatalf("len not 3")
    }

    testIntegerLiteral(t, al.Elems[0], int64(1))
    testInfixExpression(t, al.Elems[1], 2, "*", 2)
    testInfixExpression(t, al.Elems[2], 3, "+", 3)
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) bool {
    il, ok := exp.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("invalid type assertion")
        return false
    }

    if il.Value != val {
        t.Fatalf("%d expected, but got %d", val, il.Value)
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
        t.Fatalf("incorrect identifier value: %s expected, but got %s", val, id.Value)
        return false
    }

    if id.TokenLiteral() != val {
        t.Fatalf("incorrect identifier token literal")
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
        {
            "1 + (2 + 3) + 4",
            "((1 + (2 + 3)) + 4)",
        },
        {
            "(1 + 3) * 4",
            "((1 + 3) * 4)",
        },
        {
            "2 / (2 + 3)",
            "(2 / (2 + 3))",
        },
        {
            "-(2 + 3)",
            "(-(2 + 3))",
        },
        {
            "!(true == true)",
            "(!(true == true))",
        },
        {
            "a + add(b * c) + d",
            "((a + add((b * c))) + d)",
        },
        {
            "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
            "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
        },
        {
            "add(a + b + c * d / f + g)",
            "add((((a + b) + ((c * d) / f)) + g))",
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
