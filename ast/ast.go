// @file 	ast.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	ast

package ast

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/skea3344/gserrors"
)

var (
	ErrAst = errors.New("ast contract check error")
)

// Node gs文件中的解析出的对象
type Node interface {
	fmt.Stringer
	Name() string
	Path() string
	Parent() Node
	SetParent(parent Node) Node
	Package() *Package
	Attrs() []*Attr
	AddAttr(attr *Attr)
	RemoveAttr(attr *Attr)
	NewExtra(name string, data interface{})
	Extra(name string) (interface{}, bool)
	DelExtra(name string)
	Accept(visitor Visitor) Node
}

// Path 节点路径
func Path(node Node) (result []Node) {
	var nodes []Node
	current := node
	for current != nil {
		nodes = append(nodes, current)
		current = current.Parent()
	}
	for i := len(nodes) - 1; i > -1; i-- {
		result = append(result, nodes[i])
	}
	return
}

// GetAttrs 获取节点内给定属性类型的属性列表
func GetAttrs(node Node, attrType Expr) []*Attr {
	var attrs []*Attr
	for _, attr := range node.Attrs() {
		if attr.Type.Ref == attrType {
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

// BaseNode 基本节点
type BaseNode struct {
	name   string                 // 名字
	parent Node                   // 父节点
	attrs  []*Attr                // 属性列表
	extras map[string]interface{} // 额外数据
}

func (node *BaseNode) Init(name string, parent Node) {
	node.name = name
	node.parent = parent
}

func (node *BaseNode) Name() string {
	return node.name
}

func (node *BaseNode) String() string {
	return node.name
}

func (node *BaseNode) Path() string {
	var writer bytes.Buffer
	for _, node := range Path(node) {
		writer.WriteString(node.Name())
		writer.WriteRune('.')
	}
	return writer.String()
}

func (node *BaseNode) Package() *Package {
	if node.Parent() == nil {
		return nil
	}
	return node.Parent().Package()
}

func (node *BaseNode) Parent() Node {
	return node.parent
}

// SetParent 设置父节点
func (node *BaseNode) SetParent(parent Node) (old Node) {
	old, node.parent = node.parent, parent
	return
}

// getExtra 获取节点额外数据
func (node *BaseNode) getExtra() map[string]interface{} {
	if node.extras == nil {
		node.extras = make(map[string]interface{})
	}
	return node.extras
}

// Attrs 属性列表
func (node *BaseNode) Attrs() []*Attr {
	return node.attrs
}

// AddAttr 为节点添加属性
func (node *BaseNode) AddAttr(attr *Attr) {
	for _, old := range node.attrs {
		if old == attr {
			return
		}
	}
	attr.SetParent(node)
	node.attrs = append(node.attrs, attr)
}

// AddAttrs 添加多个属性
func (node *BaseNode) AddAttrs(attrs []*Attr) {
	for _, attr := range attrs {
		node.AddAttr(attr)
	}
}

// RemoveAttr 从节点的属性列表删除指定属性
func (node *BaseNode) RemoveAttr(attr *Attr) {
	var attrs []*Attr
	for _, old := range node.attrs {
		if old == attr {
			continue
		}
		attrs = append(attrs, old)
	}
	node.attrs = attrs
	attr.SetParent(nil)
}

func (node *BaseNode) NewExtra(name string, data interface{}) {
	node.getExtra()[name] = data
}

func (node *BaseNode) Extra(name string) (data interface{}, ok bool) {
	data, ok = node.getExtra()[name]
	return
}

func (node *BaseNode) DelExtra(name string) {
	delete(node.getExtra(), name)
}

func (node *BaseNode) Accept(visitor Visitor) Node {
	gserrors.Panicf(nil, "type(%s) not implement Accept", reflect.TypeOf(node))
	return nil
}

// Expr 表达式
type Expr interface {
	Node
	Script() *Script
}

// BaseExpr 基本表达式
type BaseExpr struct {
	BaseNode
	script *Script
}

// Init 初始化基本表达式 指明其所属代码节点
func (node *BaseExpr) Init(name string, script *Script) {
	defer gserrors.Ensure(func() bool {
		return script != nil
	}, "the param script can not be nil")
	node.BaseNode.Init(name, nil)
	node.script = script
}

// Script 返回基本表达式所属代码节点
func (node *BaseExpr) Script() *Script {
	gserrors.Require(node.script != nil, "the param script can not be nil")
	return node.script
}

// Package 返回基本表达式所属包
func (node *BaseExpr) Package() *Package {
	return node.Script().Package()
}
