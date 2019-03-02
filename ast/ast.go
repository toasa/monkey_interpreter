package ast

import (
    "bytes"
    "monkey_interpreter/token"
)

type Node interface {
    // nodeに関連付けられたトークンのリテラル値を返す
    TokenLiteral() string
    String() string
}

// let, return, 式文の３種のみ
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
func (p *Program) String() string {
    var out bytes.Buffer

    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type LetStatement struct {
    // let <identifier> = <expression>;
    // ex. let a = 5 * 5;
    Token token.Token // `let` token
    Name *Identifier
    Value Expression
}

// 下2つのメソッドをもって、LetStatement構造体はそれぞれStatementインターフェースと
// Nodeインターフェースを満たす
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}
func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")
    return out.String()
}

type ReturnStatement struct {
    // return <expression>;
    Token token.Token // `return` token
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")
    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }
    out.WriteString(";")
    return out.String()
}

type ExpressionStatement struct {
    // <expression>;
    Token token.Token
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}
func (es *ExpressionStatement) String() string {

    if es.Expression != nil {
        return es.Expression.String()
    }
    return ""
}

// identifierは値を生成するため式(expression)
type Identifier struct {
    Token token.Token
    Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
    return id.Token.Literal
}
func (id *Identifier) String() string {
    return id.Value
}

type IntergerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntergerLiteral) expressionNode() {}
func (il *IntergerLiteral) TokenLiteral() string {
    return il.Token.Literal
}
func (il *IntergerLiteral) String() string {
    return il.Token.Literal
}

type Boolean struct {
    Token token.Token
    Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}
func (b *Boolean) String() string {
    return b.Token.Literal
}

type PrefixExpression struct {
    // <prefix operator><expression>;
    Token token.Token // token of `!` or `-`
    Operator string
    Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.TokenLiteral())
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    // <expression><infix operator><expression>;
    Token token.Token
    Left Expression
    Operator string
    Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String() + " ")
    out.WriteString(ie.TokenLiteral())
    out.WriteString(" " + ie.Right.String())
    out.WriteString(")")

    return out.String()
}
