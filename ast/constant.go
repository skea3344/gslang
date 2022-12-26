// @file 	constant.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	constant

package ast

// String 字面量字符串常量
type String struct {
	BaseExpr
	Value string
}

// NewString 在代码节点内新建字面量字符串常量
func (node *Script) NewString(val string) *String {
	expr := &String{
		Value: val,
	}
	expr.Init("string", node)
	return expr
}

// Float 字面量浮点数常量
type Float struct {
	BaseExpr
	Value float64
}

// NewFloat 在代码节点内新建字面量浮点数常量
func (node *Script) NewFloat(val float64) *Float {
	expr := &Float{
		Value: val,
	}
	expr.Init("float", node)
	return expr
}

// Int 字面量整形常量
type Int struct {
	BaseExpr
	Value int64
}

// NewInt 在代码节点内新建字面量整形常量
func (node *Script) NewInt(val int64) *Int {
	expr := &Int{
		Value: val,
	}
	expr.Init("int", node)
	return expr
}

// Bool 字面量布尔值常量
type Bool struct {
	BaseExpr      // 内嵌基本表达式
	Value    bool // 值
}

// NewBool 在代码节点内新建字面量布尔值常量
func (node *Script) NewBool(val bool) *Bool {
	expr := &Bool{
		Value: val,
	}
	expr.Init("bool", node)
	return expr
}
