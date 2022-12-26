// @file 	package.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	package

package ast

// Package 代码包节点
type Package struct {
	BaseNode
	Scripts map[string]*Script
	Types   map[string]Expr
}

// NewPackage 生成一个新的代码包节点
func NewPackage(name string) *Package {
	node := &Package{
		Scripts: make(map[string]*Script),
		Types:   make(map[string]Expr),
	}
	node.Init(name, nil)
	return node
}

// NewType 添加类型
func (node *Package) NewType(expr Expr) (Expr, bool) {
	if old, ok := node.Types[expr.Name()]; ok {
		return old, false
	}
	node.Types[expr.Name()] = expr
	return expr, true
}

func (node *Package) Package() *Package {
	return node
}
