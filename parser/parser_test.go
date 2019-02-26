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
    p := New(l)

    program := p.ParseProgram()

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

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "let" {
        return false
    }

    letstmt, ok := s.(*ast.LetStatement)
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

// func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdent string) bool {
//     if stmt.TokenLiteral() != "let" {
//         return false
//     }
//
//     letstmt, ok := stmt.(*ast.LetStatement)
//     if !ok {
//         return false
//     }
//
//     if letstmt.Name.Value != expectedIdent {
//         return false
//     }
//
//     if letstmt.Name.TokenLiteral() != expectedIdent {
//         return false
//     }
//
//     return true
// }
