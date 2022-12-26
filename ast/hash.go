// @file 	hash.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	hash

package ast

// Map .
type Map struct {
	BaseExpr      // 内嵌基本表达式节点
	Key      Expr // 字典key类型
	Value    Expr // 字段Value类型
}

// NewMap 在代码节点内新建字典节点
func (node *Script) NewMap(key Expr, value Expr) *Map {

	expr := &Map{
		Key:   key,
		Value: value,
	}
	expr.Init(expr.Key.Name()+":"+expr.Value.Name(), node)
	expr.Key.SetParent(expr)
	expr.Value.SetParent(expr)
	return expr
}
