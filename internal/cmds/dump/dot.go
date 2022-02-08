package dump

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/nokia/ntt/ttcn3/ast"
)

func dot(n ast.Node) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	w.WriteString(`digraph {
	rankdir=LR
`)
	q := []ast.Node{n}
	toks := []ast.Token{}

	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		if !IsValid(n) {
			continue
		}
		if tok, ok := n.(ast.Token); ok {
			toks = append(toks, tok)
			continue
		}
		fmt.Fprintf(w, "\t%s %s;\n", nodeID(n), nodeProps(n))
		for _, child := range ast.Children(n) {
			if IsValid(child) {
				fmt.Fprintf(w, "\t%s -> %s;\n", nodeID(n), nodeID(child))
				q = append(q, child)
			}
		}
	}

	w.WriteString("	{ \n")
	for _, tok := range toks {
		fmt.Fprintf(w, "\t%s %s;\n", nodeID(tok), nodeProps(tok))
	}
	w.WriteString("	}")
	w.WriteString("}")
}

func IsValid(n ast.Node) bool {
	if n == nil {
		return false
	}
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		return false
	}
	if tok, ok := n.(ast.Token); ok {
		return tok.IsValid()
	}
	return true
}

func nodeID(n ast.Node) string {
	if tok, ok := n.(ast.Token); ok {
		return fmt.Sprintf("t%d", tok.Pos())
	}
	return fmt.Sprintf("n%p", n)
}

func nodeProps(n ast.Node) string {
	label := func(n ast.Node) string {
		if tok, ok := n.(ast.Token); ok {
			if tok.Kind.IsLiteral() {
				return fmt.Sprintf("%s", tok.Lit)
			}
			return fmt.Sprintf("%v", tok.Kind)
		}
		return strings.TrimPrefix(fmt.Sprintf("%T", n), "*ast.")
	}
	if tok, ok := n.(ast.Token); ok {
		return fmt.Sprintf("[label=<<B>%s</B>>; shape=box; style=filled; fillcolor=lightgrey]", label(tok))
	}
	return fmt.Sprintf("[label=\"%s\"]", label(n))
}
