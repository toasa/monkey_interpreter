package parser

import (
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/lexer"
    "monkey_interpreter/token"
)

type Parser struct {
    l *lexer.Lexer
    curToken token.Token
    peepToken token.Token
    errors []string
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

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
        return nil
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

    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
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
