package lexer

import "monkey_interpreter/token"

type Lexer struct {
    input string
    position int // 現在読む位置
    readPosition int // 次に読む位置
    ch byte // 現在検査中の文字
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

// Lexer構造体のメソッド, *がついているので参照渡しで、メソッドに渡される
func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
}

func newToken(tokentype token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokentype, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    for isSpace(l.ch) {
        l.readChar()
    }

    switch l.ch {
    case '=':
        if l.readPeep() == '=' {
            l.readChar()
            tok.Type = token.EQ
            tok.Literal = "=="
        } else {
            tok = newToken(token.ASSGIN, l.ch)
        }
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '-':
        tok = newToken(token.MINUS, l.ch)
    case '*':
        tok = newToken(token.MUL, l.ch)
    case '/':
        tok = newToken(token.DIV, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case '!':
        if l.readPeep() == '=' {
            l.readChar()
            tok.Type = token.NQ
            tok.Literal = "!="
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case '[':
        tok = newToken(token.LBRACKET, l.ch)
    case ']':
        tok = newToken(token.RBRACKET, l.ch)
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
    case 0:
        tok.Type = token.EOF
        tok.Literal = ""
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            // 早めのreturnは最後のl.readChar()を回避するため
            return tok
        } else if isDigit(l.ch) {
            tok.Literal = l.readNum()
            tok.Type = token.INT
            return tok
        } else {
            tok = newToken(token.ILLGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) readIdentifier() string {
    var start int = l.position
    var i int = start
    l.readChar()
    for isLetter(l.ch) || isDigit(l.ch) {
        i++
        l.readChar()
    }
    var ident string = l.input[start:i+1]
    return ident
}

func (l *Lexer) readNum() string {
    var start int = l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[start: l.position]
}

func (l *Lexer) readString() string {
    var i int = l.position + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
    }
    return l.input[i : l.position]
}

func (l *Lexer) readPeep() byte {
    if l.readPosition >= len(l.input) {
        return 0
    }
    return l.input[l.readPosition]
}

func isLetter(c byte) bool {
    if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'{
        return true
    }
    return false
}

func isSpace(c byte) bool {
    return c == ' ' || c == '\t' || c == '\r' || c== '\n'
}

func isDigit(c byte) bool {
    return '0' <= c && c <= '9'
}
