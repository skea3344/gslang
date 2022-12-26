// @file 	array.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	array

package ast

// Array 数组节点
type Array struct {
	BaseExpr        // 内嵌基本表达式节点
	Length   uint16 // 数组长度
	Element  Expr   // 数组元素类型
}

// NewArray 在代码节点内新建数组节点
func (node *Script) NewArray(length uint16, element Expr) *Array {
	expr := &Array{
		Length:  length,
		Element: element,
	}
	// 用数组的元素类型名字 命名此节点 并设置所属代码节点
	expr.Init(expr.Element.Name(), node)
	// 将此元素类型节点的父节点设置为此数组节点
	expr.Element.SetParent(expr)
	return expr
}

// List 链表节点
type List struct {
	BaseExpr      // 内嵌基本表达式节点
	Element  Expr // 链表元素类型
}

// NewList 在代码节点内新建链表节点
func (node *Script) NewList(element Expr) *List {
	expr := &List{
		Element: element,
	}
	expr.Init(expr.Element.Name(), node)
	expr.Element.SetParent(expr)
	return expr
}
