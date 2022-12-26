// @file 	attr.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	attr

package ast

// Attr 属性节点
type Attr struct {
	BaseExpr
	Type *TypeRef
	Args Expr
}

// NewAttr 创建属性
func (node *Script) NewAttr(attrType *TypeRef) *Attr {
	expr := &Attr{
		Type: attrType,
	}
	expr.Init(attrType.Name(), node)
	expr.Type.SetParent(expr)
	return expr
}
