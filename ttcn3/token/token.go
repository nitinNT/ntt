// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).
//
package token

import "strconv"

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT
	PREPROC

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT   // main
	INT     // 12345
	FLOAT   // 123.45
	STRING  // "abc"
	BSTRING // '101?F'H
	MODIF   // @fuzzy
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	SHL    // <<
	ROL    // <@
	SHR    // >>
	ROR    // @>
	CONCAT // &

	REDIR  // ->
	DECODE // =>
	ANY    // ?
	EXCL   // !
	RANGE  // ..
	ASSIGN // :=

	EQ // ==
	NE // !=
	LT // <
	LE // <=
	GT // >
	GE // >=

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	DOT    // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// Keywords
	MOD // mod
	REM // rem

	AND // and
	OR  // or
	XOR // xor
	NOT // not

	AND4B // and4b
	OR4B  // or4b
	XOR4B // xor4b
	NOT4B // not4b

	ADDRESS
	ALIVE
	ALL
	ALT
	ALTSTEP
	ANYKW
	BREAK
	CASE
	CHARSTRING
	COMPONENT
	CONST
	CONTINUE
	CONTROL
	DECMATCH
	DISPLAY
	DO
	ELSE
	ENCODE
	ENUMERATED
	ERROR
	EXCEPT
	EXCEPTION
	EXTENDS
	EXTENSION
	EXTERNAL
	FAIL
	FALSE
	FOR
	FRIEND
	FROM
	FUNCTION
	GOTO
	GROUP
	IF
	IFPRESENT
	IMPORT
	IN
	INCONC
	INOUT
	INTERLEAVE
	LABEL
	LANGUAGE
	LENGTH
	MAP
	MESSAGE
	MIXED
	MODIFIES
	MODULE
	MODULEPAR
	MTC
	NAN
	NOBLOCK
	NONE
	NULL
	OF
	OMIT
	ON
	OPTIONAL
	OUT
	OVERRIDE
	PARAM
	PASS
	PATTERN
	PORT
	PRESENT
	PRIVATE
	PROCEDURE
	PUBLIC
	RECORD
	REPEAT
	RETURN
	RUNS
	SELECT
	SENDER
	SET
	SIGNATURE
	STEPSIZE
	SYSTEM
	TEMPLATE
	TESTCASE
	TIMER
	TIMESTAMP
	TO
	TRUE
	TYPE
	UNION
	UNIVERSAL
	UNMAP
	VALUE
	VAR
	VARIANT
	WHILE
	WITH
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",
	PREPROC: "PREPROC",

	IDENT:   "IDENT",
	INT:     "INT",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
	BSTRING: "BSTRING",
	MODIF:   "MODIF",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	SHL:    "<<",
	ROL:    "<@",
	SHR:    ">>",
	ROR:    "@>",
	CONCAT: "&",

	REDIR:  "->",
	DECODE: "=>",
	ANY:    "?",
	EXCL:   "!",
	RANGE:  "..",

	ASSIGN: ":=",
	EQ:     "==",
	NE:     "!=",
	LT:     "<",
	LE:     "<=",
	GT:     ">",
	GE:     ">=",

	LPAREN:    "(",
	LBRACK:    "[",
	LBRACE:    "{",
	COMMA:     ",",
	DOT:       ".",
	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	MOD: "mod",
	REM: "rem",

	AND: "and",
	OR:  "or",
	XOR: "xor",
	NOT: "not",

	AND4B: "and4b",
	OR4B:  "or4b",
	XOR4B: "xor4b",
	NOT4B: "not4b",

	ADDRESS:    "address",
	ALIVE:      "alive",
	ALL:        "all",
	ALT:        "alt",
	ALTSTEP:    "altstep",
	ANYKW:      "any",
	BREAK:      "break",
	CASE:       "case",
	CHARSTRING: "charstring",
	COMPONENT:  "component",
	CONST:      "const",
	CONTINUE:   "continue",
	CONTROL:    "control",
	DECMATCH:   "decmatch",
	DISPLAY:    "display",
	DO:         "do",
	ELSE:       "else",
	ENCODE:     "encode",
	ENUMERATED: "enumerated",
	ERROR:      "error",
	EXCEPT:     "except",
	EXCEPTION:  "exception",
	EXTENDS:    "extends",
	EXTENSION:  "extension",
	EXTERNAL:   "external",
	FAIL:       "fail",
	FALSE:      "false",
	FOR:        "for",
	FRIEND:     "friend",
	FROM:       "from",
	FUNCTION:   "function",
	GOTO:       "goto",
	GROUP:      "group",
	IF:         "if",
	IFPRESENT:  "ifpresent",
	IMPORT:     "import",
	IN:         "in",
	INCONC:     "inconc",
	INOUT:      "inout",
	INTERLEAVE: "interleave",
	LABEL:      "label",
	LANGUAGE:   "language",
	LENGTH:     "length",
	MAP:        "map",
	MESSAGE:    "message",
	MIXED:      "mixed",
	MODIFIES:   "modifies",
	MODULE:     "module",
	MODULEPAR:  "modulepar",
	MTC:        "mtc",
	NAN:        "not_a_number",
	NOBLOCK:    "noblock",
	NONE:       "none",
	NULL:       "null",
	OF:         "of",
	OMIT:       "omit",
	ON:         "on",
	OPTIONAL:   "optional",
	OUT:        "out",
	OVERRIDE:   "override",
	PARAM:      "param",
	PASS:       "pass",
	PATTERN:    "pattern",
	PORT:       "port",
	PRESENT:    "present",
	PRIVATE:    "private",
	PROCEDURE:  "procedure",
	PUBLIC:     "public",
	RECORD:     "record",
	REPEAT:     "repeat",
	RETURN:     "return",
	RUNS:       "runs",
	SELECT:     "select",
	SENDER:     "sender",
	SET:        "set",
	SIGNATURE:  "signature",
	STEPSIZE:   "stepsize",
	SYSTEM:     "system",
	TEMPLATE:   "template",
	TESTCASE:   "testcase",
	TIMER:      "timer",
	TIMESTAMP:  "timestamp",
	TO:         "to",
	TRUE:       "true",
	TYPE:       "type",
	UNION:      "union",
	UNIVERSAL:  "universal",
	UNMAP:      "unmap",
	VALUE:      "value",
	VAR:        "var",
	VARIANT:    "variant",
	WHILE:      "while",
	WITH:       "with",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 15
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (tok Token) Precedence() int {
	switch tok {
	case ASSIGN:
		return 1
	case DECODE:
		return 2
	case RANGE:
		return 3
	case EXCL:
		return 4
	case OR:
		return 5
	case XOR:
		return 6
	case AND:
		return 7
	case NOT:
		return 8
	case EQ, NE:
		return 9
	case LT, LE, GT, GE:
		return 10
	case SHR, SHL, ROR, ROL:
		return 11
	case OR4B:
		return 12
	case XOR4B:
		return 13
	case AND4B:
		return 14
	case NOT4B:
		return 15
	case ADD, SUB, CONCAT:
		return 16
	case MUL, DIV, REM, MOD:
		return 17
	}
	return LowestPrec
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }
