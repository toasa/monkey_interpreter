package parser

import (
    "fmt"
    "strconv"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
    "monkey_interpreter/token"
)

type Parser struct {
    l *lexer.Lexer
    errors []string

    curToken token.Token
    peepToken token.Token

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns map[token.TokenType]infixParseFn
}

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

const (
    _ int = iota
    LOWEST
    EQUALS // ==
    LESSGREATER // > or <
    SUM // +
    PRODUCT // *
    PREFIX // -x or !x
    CALL // func(x)
)

var precedences = map[token.TokenType]int {
    token.EQ: EQUALS,
    token.NQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.MUL: PRODUCT,
    token.DIV: PRODUCT,
    // 中置記法としての`(`. 関数呼び出しに用いられる
    token.LPAREN: CALL,
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.infixParseFns = make(map[token.TokenType]infixParseFn)

    // どうして識別子がprefixにparseされる(prefixParseFns内の関数でのparse)の？
    // => 理由の１つにparseExpression()内の
    //      pre_fn := p.prefixParseFns[p.curToken.Type]
    // でpre_fn == nilの場合errorとしたいから
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.TRUE, p.parseBoolean)
    p.registerPrefix(token.FALSE, p.parseBoolean)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
    p.registerPrefix(token.IF, p.parseIfExpression)
    p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.MUL, p.parseInfixExpression)
    p.registerInfix(token.DIV, p.parseInfixExpression)
    // 関数呼び出しは add(2, 5); となるが
    // addはprefixのIdentifierとしてparseされ、
    // `(`が来て、引数2へと続く。つまり`(`を中置記号とも見なせる
    p.registerInfix(token.LPAREN, p.parseFunctionCall)

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peepError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, but got %s instead",
    t, p.peepToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
    p.curToken = p.peepToken
    p.peepToken = p.l.NextToken()
}

func (p *Parser) expectToken(tt token.TokenType) bool {
    if p.peepToken.Type == tt {
        p.nextToken()
        return true
    }
    return false
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    // = の右辺値は型([]ast.Statement)ではなく、値を持つ式([]ast.Statement{})でないといけない
    // 例えば var n = int とはできない
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}
    if !p.expectPeep(token.IDENT) {
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeep(token.ASSGIN) {
        return nil
    }

    p.nextToken()
    stmt.Value = p.parseExpression(LOWEST)

    if p.peepTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    p.nextToken()
    stmt.ReturnValue = p.parseExpression(LOWEST)

    if p.peepTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}
    stmt.Expression = p.parseExpression(LOWEST)

    if p.peepTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    bs := &ast.BlockStatement{Token: p.curToken}
    // bs.Statements = []ast.Statement{}

    p.nextToken()

    for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        if stmt != nil {
            bs.Statements = append(bs.Statements, stmt)
        }
        p.nextToken()
    }

    return bs
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    pre_fn := p.prefixParseFns[p.curToken.Type]
    if pre_fn == nil {
        msg := "not fonnd prefix parse function"
        p.errors = append(p.errors, msg)
        return nil
    }
    leftExp := pre_fn()

    for !p.peepTokenIs(token.SEMICOLON) && precedence < p.peepPrecedence() {
        infix := p.infixParseFns[p.peepToken.Type]

        if infix == nil {
            return leftExp
        }
        p.nextToken()

        leftExp = infix(leftExp)
    }

    return leftExp
}

// 戻り値の型をast.Expressionで*ast.Identifierとしないのは、
// parseIdentifier()をprefixParseFn型と扱えるようにし、
// map prefixParseFnsに登録したいため。以降のparseExpression系関数も同様
func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("cannot parse %q as integer", p.curToken.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }
    return &ast.IntergerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseBoolean() ast.Expression {
    if !(p.curToken.Literal == "true" || p.curToken.Literal == "false") {
        p.errors = append(p.errors, "token appeares neither true or false")
        return nil
    }

    b := &ast.Boolean{Token: p.curToken}
    b.Value = (p.curToken.Literal == "true")

    return b
}

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.nextToken()

    // parseExpression()内で`)`にぶつかるまで通常の式としてparseされる。
    // `)`のprecedenceはlowestなので、`)`以降までparseされることはない。
    ex := p.parseExpression(LOWEST)

    if !p.expectPeep(token.RPAREN) {
        return nil
    }

    return ex
}

func (p *Parser) parseIfExpression() ast.Expression {
    ie := &ast.IfExpression{Token: p.curToken}

    p.expectPeep(token.LPAREN)
    ie.Cond = p.parseExpression(LOWEST)

    p.expectPeep(token.LBRACE)
    ie.Cons = p.parseBlockStatement()

    if !p.peepTokenIs(token.ELSE) {
        return ie
    }

    p.nextToken()

    p.expectPeep(token.LBRACE)
    ie.Alt = p.parseBlockStatement()

    return ie

}

func (p *Parser) parseFunctionLiteralParams() []*ast.Identifier {

    params := []*ast.Identifier{}
    p.nextToken()
    if p.curTokenIs(token.RPAREN) {
        p.nextToken()
        return params
    }

    for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
        id, ok := p.parseIdentifier().(*ast.Identifier)
        if !ok {
            p.errors = append(p.errors, "type assertion to (*ast.identifier) error")
        }
        params = append(params, id)

        if p.peepTokenIs(token.COMMA) {
            p.nextToken()
        }
        p.nextToken()
    }

    p.nextToken()

    return params
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
    fl := &ast.FunctionLiteral{Token: p.curToken}
    p.expectPeep(token.LPAREN)

    fl.Params = p.parseFunctionLiteralParams()

    if p.peepTokenIs(token.RBRACE) {
        p.nextToken()
        return fl
    }
    fl.Body = p.parseBlockStatement()

    return fl
}

func (p *Parser) parseFunctionCallArgs() []ast.Expression {
    fca := []ast.Expression{}

    p.nextToken()

    if p.curTokenIs(token.RPAREN) {
        return fca
    }

    for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
        exp := p.parseExpression(LOWEST)
        fca = append(fca, exp)
        if p.peepTokenIs(token.COMMA) {
            p.nextToken()
        }
        p.nextToken()
    }

    return fca
}

func (p *Parser) parseFunctionCall(f ast.Expression) ast.Expression {
    fc := &ast.FunctionCall{Token: p.curToken, Func: f}
    fc.Args = p.parseFunctionCallArgs()
    return fc
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    pe := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
    p.nextToken()
    pe.Right = p.parseExpression(PREFIX)

    return pe
}

func (p *Parser) parseInfixExpression(leftExp ast.Expression) ast.Expression {
    ie := &ast.InfixExpression{
        Token: p.curToken,
        Left: leftExp,
        Operator: p.curToken.Literal,
    }

    precedence := p.curPrecedence()
    p.nextToken()

    ie.Right = p.parseExpression(precedence)

    return ie
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peepTokenIs(t token.TokenType) bool {
    return p.peepToken.Type == t
}

func (p *Parser) expectPeep(t token.TokenType) bool {
    if p.peepTokenIs(t) {
        p.nextToken()
        return true
    } else {
        p.peepError(t)
        return false
    }
}

// map prefixParseFnsへエントリの追加
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) peepPrecedence() int {
    n, ok := precedences[p.peepToken.Type]
    if ok {
        return n
    }
    return LOWEST
}

func (p *Parser) curPrecedence() int {
    n, ok := precedences[p.curToken.Type]
    if ok {
        return n
    }
    return LOWEST
}
