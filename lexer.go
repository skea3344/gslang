// @file 	lexer.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	lexer

package gslang

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/skea3344/gserrors"
	"github.com/skea3344/logger"
)

var (
	// ErrLexer 词法分析器错误
	ErrLexer = errors.New("lexer error")
)

const (
	TokenEOF        rune = -(iota + 1) // TokenEOF 文件结束符
	TokenID                            // TokenID 标识符
	TokenINT                           // TokenINT 整形
	TokenFLOAT                         // TokenFLOAT 浮点
	TokenTrue                          // TokenTrue 真
	TokenFalse                         // TokenFalse 假
	TokenSTRING                        // TokenSTRING 字符串字面量
	TokenCOMMENT                       // TokenCOMMENT 注释
	TokenLABEL                         // TokenLABEL 标签
	TokenArrowRight                    // TokenArrowRight ->
	KeyByte                            // KeyByte byte
	KeySByte                           // KeySByte sbyte
	KeyUInt16                          // KeyUInt16 uint16
	KeyInt16                           // KeyInt16 int16
	KeyUInt32                          // KeyUInt32 uint32
	KeyInt32                           // KeyInt32 int32
	KeyUInt64                          // KeyUInt64 uint64
	KeyInt64                           // KeyInt64 int64
	KeyFloat32                         // KeyFloat32 float32
	KeyFloat64                         // KeyFloat64 float64
	KeyBool                            // KeyBool bool
	KeyEnum                            // KeyEnum enum
	KeyString                          // KeyString string
	KeyStruct                          // KeyStruct struct
	KeyTable                           // KeyTable table
	KeyContract                        // KeyContract contract
	KeyImport                          // KeyImport import
	KeyMap                             // KeyMap map
)

var tokenName = map[rune]string{
	TokenEOF:   "EOF",
	TokenID:    "ID",
	TokenINT:   "INT",
	TokenFLOAT: "FLOAT",
	// TokenTrue:
	// TokenFalse:
	TokenSTRING:     "STRING",
	TokenCOMMENT:    "COMMENT",
	TokenLABEL:      "LABEL",
	TokenArrowRight: "->",
	KeyByte:         "byte",
	KeySByte:        "sbyte",
	KeyUInt16:       "uint16",
	KeyInt16:        "int16",
	KeyUInt32:       "uint32",
	KeyInt32:        "int32",
	KeyUInt64:       "uint64",
	KeyInt64:        "int64",
	KeyFloat32:      "float32",
	KeyFloat64:      "float64",
	KeyBool:         "bool",
	KeyEnum:         "enum",
	KeyString:       "string",
	KeyStruct:       "struct",
	KeyTable:        "table",
	KeyContract:     "contract",
	KeyImport:       "import",
	KeyMap:          "map",
}

var keyMap = map[string]rune{
	"byte":     KeyByte,
	"sbyte":    KeySByte,
	"int16":    KeyInt16,
	"uint16":   KeyUInt16,
	"int32":    KeyInt32,
	"uint32":   KeyUInt32,
	"int64":    KeyInt64,
	"uint64":   KeyUInt64,
	"float32":  KeyFloat32,
	"float64":  KeyFloat64,
	"string":   KeyString,
	"bool":     KeyBool,
	"enum":     KeyEnum,
	"struct":   KeyStruct,
	"table":    KeyTable,
	"contract": KeyContract,
	"import":   KeyImport,
	"map":      KeyMap,
}

// TokenName 取Token类型rune对应的字符串表示 大于0的为字符本身 小于0的为内置类型
func TokenName(token rune) string {
	if token > 0 {
		return string(token)
	}
	return tokenName[token]
}

// Token 一个gslang符号对象
type Token struct {
	Type  rune        // 符号类型
	Value interface{} // 符号值
	Pos   Position    // 符号在代码中的位置
}

// NewToken 新建一个符号对象
func NewToken(t rune, val interface{}) *Token {
	return &Token{
		Type:  t,
		Value: val,
	}
}

// String 符号对象的字符串输出显示
func (token *Token) String() string {
	if token.Value != nil {
		return fmt.Sprintf("token[%s]\n\tval:%v\n\tpos:%s", TokenName(token.Type), token.Value, token.Pos)
	}
	return fmt.Sprintf("token[%s]\n\tpos:%s", TokenName(token.Type), token.Pos)
}

// Lexer 词法分析器
type Lexer struct {
	logger.ILog // 内嵌通用日志接口
	reader      *bufio.Reader
	position    Position
	token       *Token
	buff        [utf8.UTFMax]byte
	buffPos     int
	offset      int
	ws          uint64
	curr        rune
}

// NewLexer 新建一个词法分析器
func NewLexer(filename string, reader io.Reader) *Lexer {
	return &Lexer{
		ILog:   logger.Get("gslang[lexer]"),
		reader: bufio.NewReader(reader),
		ws:     1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ',
		position: Position{
			Filename: filename,
			Line:     1,
			Column:   1,
		},
		curr: TokenEOF,
	}
}

// newerror 返回一个yferrors.YFError
func (lexer *Lexer) newerror(fmtstring string, args ...interface{}) error {
	return gserrors.Newf(ErrLexer, "[lexer] %s\n\t%s", fmt.Sprintf(fmtstring, args...), lexer.position)
}

// nextChar 读取下一个utf8字符
func (lexer *Lexer) nextChar() error {
	// 从reader中读取一个字节
	c, err := lexer.reader.ReadByte()
	if err != nil {
		// 如果是文件结束标志 则返回
		if err == io.EOF {
			lexer.curr = TokenEOF
			return nil
		}
		return err
	}
	// 偏移量加1
	lexer.offset++
	if c >= utf8.RuneSelf {
		// 保存字节到buff buffPos加1
		lexer.buff[0] = c
		lexer.buffPos = 1
		// 循环直到buff中以一个码值的完整utf8编码开始
		for !utf8.FullRune(lexer.buff[0:lexer.buffPos]) {
			c, err = lexer.reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					lexer.curr = TokenEOF
					return nil
				}
				return err
			}
			lexer.buff[lexer.buffPos] = c
			lexer.buffPos++
			// buffPos必须小于buff的长度
			gserrors.Assert(lexer.buffPos < len(lexer.buff), "utf8.UTFMax must < len(lexer.buff)")
		}
		// 对buff进行utf8解码 得到一个utf8编码的rune
		c, width := utf8.DecodeRune(lexer.buff[0:lexer.buffPos])
		if c == utf8.RuneError && width == 1 {
			return lexer.newerror("illegal utf8 character")
		}
		lexer.curr = c
	} else { // ASCII
		lexer.curr = rune(c)
	}
	// 列号加1
	lexer.position.Column++
	return nil
}

// isDecimal 判断是不是0-9的小数字
func isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// lexer 分析并生成下一个Token
func (lexer *Lexer) next() (token *Token, err error) {
	// 如果当前的是EOF
	if TokenEOF == lexer.curr {
		// 读取下一个rune
		if err = lexer.nextChar(); err != nil {
			return
		}
	}
	// 忽略 \t \r 空格  回车的话 处理位置
	for lexer.ws&(1<<uint(lexer.curr)) != 0 {
		if lexer.curr == '\n' {
			lexer.position.Column = 0
			lexer.position.Line++
		}
		if err = lexer.nextChar(); err != nil {
			return
		}
	}
	// 如果是TokenEOF 返回EOF Token
	if lexer.curr == TokenEOF {
		token = NewToken(TokenEOF, nil)
		token.Pos = lexer.position
		return
	}
	// 拷贝一份位置信息
	position := lexer.position
	// 按类型处理
	switch {
	// 如果是字母或者下划线 那就认为接下来的是一个 标识符
	case unicode.IsLetter(lexer.curr) || lexer.curr == '_':
		// 扫描获得一个 标识符 Token, Type 为TokenID
		token, err = lexer.scanID()
		if err == nil {
			// 为特殊的标识符设置特殊类型
			id := token.Value.(string)
			if id == "true" { // true -> TokenTrue
				token.Type = TokenTrue
			} else if id == "false" { // false -> TokenFalse
				token.Type = TokenFalse
			} else if key, ok := keyMap[id]; ok { // keyMap 中的类型 标志
				token.Type = key // keyMap中的key对应的标识符 类型为其对应val
			} else { // 例如  (lang: 认为lang是一个标签,但是保存的value仅仅是lang,:被丢弃了,类型更改为TokenLABEL)
				if lexer.curr == ':' {
					token.Type = TokenLABEL
					_ = lexer.nextChar()
				}
			}
		}
	case isDecimal(lexer.curr): // 如果是一个 0-9 的数字, 则分类为数字扫描
		token, err = lexer.scanNum()
	case lexer.curr == '"': // " 进入字符串字面量扫描
		token, err = lexer.scanString('"')
	case lexer.curr == '\'': // '  进入字符串字面量扫描
		token, err = lexer.scanString('\'')
	case lexer.curr == '/': // / 判断是不是注释 单行或者块注释
		err = lexer.nextChar()
		if err == nil {
			if lexer.curr == '/' || lexer.curr == '*' {
				token, err = lexer.scanComment(lexer.curr)
			} else { // 不是注释 则 以/后的那一个rune 作为类型返回 nil值的Token
				token = NewToken(lexer.curr, nil)
			}
		}
	case lexer.curr == '-': // 如果是- 则判断是不是->
		err = lexer.nextChar()
		if err == nil {
			if lexer.curr == '>' {
				token = NewToken(TokenArrowRight, nil)
				err = lexer.nextChar()
			} else {
				token = NewToken('-', nil)
			}
		}
	default: // 其他情况返回 rune 本身作为类型 值为nil 的Token
		token = NewToken(lexer.curr, nil)
		lexer.curr = TokenEOF
	}
	if err == nil {
		token.Pos = position
	}
	return
}

// scanComment 判断接下来的块是不是注释  返回TokenCOMMENT的Token
func (lexer *Lexer) scanComment(ch rune) (*Token, error) {
	// buff存储注释内容
	var buff bytes.Buffer
	// 单行注释
	if ch == '/' {
		// 取下一个rune
		err := lexer.nextChar()
		if err != nil {
			return nil, err
		}
		// 只要rune!=\n且rune>=0就保存为注释内容
		for lexer.curr != '\n' && lexer.curr >= 0 {
			buff.WriteRune(lexer.curr)
			err := lexer.nextChar()
			if err != nil {
				return nil, err
			}
		}
		// 返回TokenCOMMENT　Token
		return NewToken(TokenCOMMENT, buff.String()), nil
	}
	err := lexer.nextChar()
	if err != nil {
		return nil, err
	}
	// 块注释　判断　　/* */
	for {
		if lexer.curr < 0 {
			return nil, lexer.newerror("comment not terminated")
		}
		if lexer.curr == '\n' {
			lexer.position.Column = 0
			lexer.position.Line++
		}
		ch0 := lexer.curr
		err := lexer.nextChar()
		if err != nil {
			return nil, err
		}
		// 循环出口为检查到 */
		if ch0 == '*' && lexer.curr == '/' {
			err := lexer.nextChar()
			if err != nil {
				return nil, err
			}
			break
		}
		// 合法的rune均保存为注释内容
		buff.WriteRune(ch0)
	}
	return NewToken(TokenCOMMENT, buff.String()), nil
}

// scanEscape 判断是不是字符串中需要转义的括号内容　\" \'
func (lexer *Lexer) scanEscape(buff *bytes.Buffer, quote rune) (err error) {
	err = lexer.nextChar()
	if err != nil {
		return
	}
	switch lexer.curr {
	case quote:
		buff.WriteRune(lexer.curr)
		err = lexer.nextChar()
		if err != nil {
			return
		}
	default:
		err = lexer.newerror("illegal char escape")
	}
	return
}

// scanString 字符串字面量判断 quote 指明是哪种引号
func (lexer *Lexer) scanString(quote rune) (token *Token, err error) {
	var buff bytes.Buffer
	err = lexer.nextChar()
	if err != nil {
		return nil, err
	}
	for lexer.curr != quote {
		// 字符串不能换行
		if lexer.curr == '\n' || lexer.curr < 0 {
			err = lexer.newerror("literal not terminated")
			return
		}
		if lexer.curr == '\\' {
			// 判断转义内容 \" \' 需要处理
			_ = lexer.scanEscape(&buff, quote)
		} else {
			// 其余作为字符串内容写入
			buff.WriteRune(lexer.curr)
			err = lexer.nextChar()
			if err != nil {
				return nil, err
			}
		}
	}
	err = lexer.nextChar()
	if err != nil {
		return nil, err
	}
	// 返回TokenSTRING类Token
	token = NewToken(TokenSTRING, buff.String())
	return
}

// scanID 判断标识符
func (lexer *Lexer) scanID() (token *Token, err error) {
	var buff bytes.Buffer
	// 以_或者字母开头　以_或者字母或者十进制数字后续　的标识符
	for lexer.curr == '_' || unicode.IsLetter(lexer.curr) || unicode.IsDigit(lexer.curr) {
		buff.WriteRune(lexer.curr)
		if err = lexer.nextChar(); err != nil {
			return nil, err
		}
	}
	// 返回一个TokenID 类Token
	token = NewToken(TokenID, buff.String())
	return
}

// digitVal 返回十六进制rune表示的十进制数值
func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16
}

// scanNum 判断数字字面量
func (lexer *Lexer) scanNum() (*Token, error) {
	var buff bytes.Buffer
	if lexer.curr == '0' {
		buff.WriteRune(lexer.curr)
		_ = lexer.nextChar()
		// 判断是是不是十六进制数值
		if lexer.curr == 'x' || lexer.curr == 'X' {
			buff.WriteRune(lexer.curr)
			_ = lexer.nextChar()
			for digitVal(lexer.curr) < 16 {
				buff.WriteRune(lexer.curr)
				_ = lexer.nextChar()
			}
			if buff.Len() < 3 {
				return nil, lexer.newerror("illegal hexadecimal num")
			}
			val, err := strconv.ParseInt(buff.String(), 0, 64)
			if err != nil {
				return nil, lexer.newerror(err.Error())
			}
			return NewToken(TokenINT, val), nil
		}
	}
	lexer.scanMantissa(&buff)
	switch lexer.curr {
	case '.', 'e', 'E':
		lexer.scanFraction(&buff)
		lexer.scanExponent(&buff)
		// 返回浮点数
		val, err := strconv.ParseFloat(buff.String(), 64)
		if err != nil {
			return nil, lexer.newerror(err.Error())
		}
		return NewToken(TokenFLOAT, val), nil
	}
	// 返回整数
	val, err := strconv.ParseInt(buff.String(), 0, 64)
	if err != nil {
		return nil, lexer.newerror(err.Error())
	}
	return NewToken(TokenINT, val), nil
}

// scanMantissa 扫描尾数　只要是连续0-9数字就合法
func (lexer *Lexer) scanMantissa(buff *bytes.Buffer) {
	for isDecimal(lexer.curr) {
		buff.WriteRune(lexer.curr)
		_ = lexer.nextChar()
	}
}

// scanFraction 扫描分数 小数部分
func (lexer *Lexer) scanFraction(buff *bytes.Buffer) {
	if lexer.curr == '.' {
		buff.WriteRune(lexer.curr)
		_ = lexer.nextChar()
		lexer.scanMantissa(buff)
	}
}

// scanExponent 扫描指数
func (lexer *Lexer) scanExponent(buff *bytes.Buffer) {
	if lexer.curr == 'e' || lexer.curr == 'E' {
		buff.WriteRune(lexer.curr)
		_ = lexer.nextChar()
		if lexer.curr == '-' || lexer.curr == '+' {
			buff.WriteRune(lexer.curr)
			_ = lexer.nextChar()
		}
		lexer.scanMantissa(buff)
	}
}

// Peek 看一眼 返回分析器当面的token 如果为nil则获取下一个Token保存并返回
// 特点是无论如何lexer的token不会为空
func (lexer *Lexer) Peek() (token *Token, err error) {
	if lexer.token != nil {
		token = lexer.token
		return
	}
	token, err = lexer.next()
	if err == nil {
		lexer.token = token
	}
	return
}

// Next 返回分析器当前token 并清空 如果为nil则直接返回下一个token
// 特点是lexer的token一定会被置为nil
func (lexer *Lexer) Next() (token *Token, err error) {
	if lexer.token != nil {
		token, lexer.token = lexer.token, nil
		return
	}
	token, err = lexer.next()
	return
}
