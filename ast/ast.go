package ast

import (
    "monkey_interpreter/token"
)

type Node interface {
    // nodeに関連付けられたトークンのリテラル値を返す
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

// ASTのroot node
// すべてのmonkeyプログラムは文の集まり
type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

type LetStatement struct {
    // let <identifier> = <expression>;
    // ex. let a = 5 * 5;
    Name *Identifier
    Value Expression
    Token token.Token
}

// 下2つのメソッドをもって、LetStatement構造体はそれぞれStatementインターフェースと
// Nodeインターフェースを満たす
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

// identifierは値を生成するため式
type Identifier struct {
    Token token.Token
    Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
    return id.Token.Literal
}
