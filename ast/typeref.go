// @file 	typeref.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	typeref

package ast

import (
	"bytes"

	"github.com/skea3344/gserrors"
)

// Type 类型引用
type TypeRef struct {
	BaseExpr
	Ref      Expr
	NamePath []string
}

// NewTypeRef 在代码内创建类型引用
func (node *Script) NewTypeRef(namePath []string) *TypeRef {
	gserrors.Require(len(namePath) > 0, "namePath can not be nil")
	expr := &TypeRef{
		NamePath: namePath,
	}
	var buff bytes.Buffer
	for _, nodeName := range namePath {
		buff.WriteRune('.')
		buff.WriteString(nodeName)
	}
	// 类型引用的名字为namePath按.连接  所属代码节点为此代码节点
	// 如引用ast.TypeRef  namePath =["ast","TypeRef"] name="ast.TypeRef"
	expr.Init(buff.String(), node)
	return expr
}
