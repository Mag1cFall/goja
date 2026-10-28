package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mag1cFall/goja/file"
	"github.com/Mag1cFall/goja/token"
)

const (
	err_UnexpectedToken      = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
	err_UnexpectedEscape     = "Unexpected escape"
)

//    UnexpectedNumber:  'Unexpected number',
//    UnexpectedString:  'Unexpected string',
//    UnexpectedIdentifier:  'Unexpected identifier',
//    UnexpectedReserved:  'Unexpected reserved word',
//    NewlineAfterThrow:  'Illegal newline after throw',
//    InvalidRegExp: 'Invalid regular expression',
//    UnterminatedRegExp:  'Invalid regular expression: missing /',
//    InvalidLHSInAssignment:  'Invalid left-hand side in assignment',
//    InvalidLHSInForIn:  'Invalid left-hand side in for-in',
//    MultipleDefaultsInSwitch: 'More than one default clause in switch statement',
//    NoCatchOrFinally:  'Missing catch or finally after try',
//    UnknownLabel: 'Undefined label \'%0\'',
//    Redeclaration: '%0 \'%1\' has already been declared',
//    IllegalContinue: 'Illegal continue statement',
//    IllegalBreak: 'Illegal break statement',
//    IllegalReturn: 'Illegal return statement',
//    StrictModeWith:  'Strict mode code may not include a with statement',
//    StrictCatchVariable:  'Catch variable may not be eval or arguments in strict mode',
//    StrictVarName:  'Variable name may not be eval or arguments in strict mode',
//    StrictParamName:  'Parameter name eval or arguments is not allowed in strict mode',
//    StrictParamDupe: 'Strict mode function may not have duplicate parameter names',
//    StrictFunctionName:  'Function name may not be eval or arguments in strict mode',
//    StrictOctalLiteral:  'Octal literals are not allowed in strict mode.',
//    StrictDelete:  'Delete of an unqualified identifier in strict mode.',
//    StrictDuplicateProperty:  'Duplicate data property in object literal not allowed in strict mode',
//    AccessorDataProperty:  'Object literal may not have data and accessor property with the same name',
//    AccessorGetSet:  'Object literal may not have multiple get/set accessors with the same name',
//    StrictLHSAssignment:  'Assignment to eval or arguments is not allowed in strict mode',
//    StrictLHSPostfix:  'Postfix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictLHSPrefix:  'Prefix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictReservedWord:  'Use of future reserved word in strict mode'

// A SyntaxError is a description of an ECMAScript syntax error.

// An Error represents a parsing error. It includes the position where the error occurred and a message/description.
type Error struct {
	Position file.Position
	Message  string
}

// FIXME Should this be "SyntaxError"?

func (self Error) Error() string {
	filename := self.Position.Filename
	if filename == "" {
		filename = "(anonymous)"
	}
	return fmt.Sprintf("%s: Line %d:%d %s",
		filename,
		self.Position.Line,
		self.Position.Column,
		self.Message,
	)
}

func (self *_parser) error(place interface{}, msg string, msgValues ...interface{}) *Error {
	idx := file.Idx(0)
	switch place := place.(type) {
	case int:
		idx = self.idxOf(place)
	case file.Idx:
		if place == 0 {
			idx = self.idxOf(self.chrOffset)
		} else {
			idx = place
		}
	default:
		panic(fmt.Errorf("error(%T, ...)", place))
	}

	position := self.position(idx)
	msg = fmt.Sprintf(msg, msgValues...)
	self.errors.Add(position, msg)
	return self.errors[len(self.errors)-1]
}

func (self *_parser) errorUnexpected(idx file.Idx, chr rune) error {
	if chr == -1 {
		return self.error(idx, err_UnexpectedEndOfInput)
	}
	return self.error(idx, err_UnexpectedToken, token.ILLEGAL)
}

func (self *_parser) errorUnexpectedToken(tkn token.Token) error {
	if tkn == token.ESCAPED_RESERVED_WORD {
		return self.error(self.idx, "Keyword must not contain escaped characters")
	}
	return self.error(self.idx, "unexpected token: %s", self.describe(tkn))
}

// errorExpectedExpression 按 SpiderMonkey 报告当前位置缺少表达式，yield 与 await 报告为保留标识符
func (self *_parser) errorExpectedExpression() error {
	if self.token == token.YIELD || self.token == token.AWAIT {
		return self.error(self.idx, "%s is a reserved identifier", self.token.String())
	}
	return self.error(self.idx, "expected expression, got %s", self.describe(self.token))
}

// describe 返回 SpiderMonkey TokenKindToDesc 对词法单元 tkn 的描述，tkn 为当前词法单元时按其源码文本区分
func (self *_parser) describe(tkn token.Token) string {
	literal := tkn.String()
	if tkn == self.token {
		literal = self.literal
	}
	return tokenDescription(tkn, literal)
}

// tokenDescription 返回 SpiderMonkey TokenKindToDesc 对源码文本为 literal 的词法单元 tkn 的描述
func tokenDescription(tkn token.Token, literal string) string {
	switch tkn {
	case token.EOF:
		return "end of script"
	case token.IDENTIFIER:
		switch literal {
		case "as", "async", "await", "each", "from", "get", "let", "meta", "of", "set", "static", "target", "yield":
			return "'" + literal + "'"
		case "implements", "interface", "package", "private", "protected", "public":
			return "reserved word '" + literal + "'"
		}
		return "identifier"
	case token.PRIVATE_IDENTIFIER:
		return "private identifier"
	case token.NUMBER:
		if strings.HasSuffix(literal, "n") {
			return "bigint literal"
		}
		return "numeric literal"
	case token.STRING:
		return "string literal"
	case token.BACKTICK:
		return "template literal"
	case token.BOOLEAN:
		return "boolean literal '" + literal + "'"
	case token.NULL:
		return "null literal"
	case token.KEYWORD:
		if literal == "enum" {
			return "reserved word 'enum'"
		}
		return "keyword '" + literal + "'"
	case token.LET, token.STATIC, token.ASYNC, token.AWAIT, token.YIELD, token.OF:
		return "'" + tkn.String() + "'"
	}
	if tkn >= token.IF && tkn <= token.INSTANCEOF {
		return "keyword '" + tkn.String() + "'"
	}
	return "'" + tkn.String() + "'"
}

// ErrorList is a list of *Errors.
type ErrorList []*Error

// Add adds an Error with given position and message to an ErrorList.
func (self *ErrorList) Add(position file.Position, msg string) {
	*self = append(*self, &Error{position, msg})
}

// Reset resets an ErrorList to no errors.
func (self *ErrorList) Reset() { *self = (*self)[0:0] }

func (self ErrorList) Len() int      { return len(self) }
func (self ErrorList) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self ErrorList) Less(i, j int) bool {
	x := &self[i].Position
	y := &self[j].Position
	if x.Filename < y.Filename {
		return true
	}
	if x.Filename == y.Filename {
		if x.Line < y.Line {
			return true
		}
		if x.Line == y.Line {
			return x.Column < y.Column
		}
	}
	return false
}

func (self ErrorList) Sort() {
	sort.Sort(self)
}

// Error implements the Error interface.
func (self ErrorList) Error() string {
	switch len(self) {
	case 0:
		return "no errors"
	case 1:
		return self[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", self[0].Error(), len(self)-1)
}

// Err returns an error equivalent to this ErrorList.
// If the list is empty, Err returns nil.
func (self ErrorList) Err() error {
	if len(self) == 0 {
		return nil
	}
	return self
}
