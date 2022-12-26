// @file 	op.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	op

package ast

// BinaryOp 二元运算
type BinaryOp struct {
	BaseExpr
	Left  Expr // 左操作数
	Right Expr // 右操作数
}

// BinaryOp 代码内创建二元运算
func (node *Script) NewBinaryOp(name string, left Expr, right Expr) *BinaryOp {
	op := &BinaryOp{
		Left:  left,
		Right: right,
	}
	op.Init(name, node)
	return op
}

// UnaryOp 一元运算
type UnaryOp struct {
	BaseExpr
	Right Expr
}

// NewUnaryOp 代码内创建一元运算
func (node *Script) NewUnaryOp(name string) *UnaryOp {
	op := &UnaryOp{}
	op.Init(name, node)
	return op
}
