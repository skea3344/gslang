// @file 	visitor.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	visitor

package ast

import "errors"

var (
	// ErrVisit 掉用访问方法时的错误
	ErrVisit = errors.New("invalid call to visit method")
)

// Visitor 访问者接口 用来访问各种节点表达式结构
type Visitor interface {
	VisitPackage(*Package) Node     // 访问 包
	VisitScript(*Script) Node       // 访问代码
	VisitEnum(*Enum) Node           // 访问 枚举
	VisitEnumVal(*EnumVal) Node     // 访问 单个枚举值
	VisitTable(*Table) Node         // 访问 表
	VisitField(*Field) Node         // 访问 域
	VisitContract(*Contract) Node   // 访问 协议
	VisitMethod(*Method) Node       // 访问 函数
	VisitParam(*Param) Node         // 访问 参数
	VisitTypeRef(*TypeRef) Node     // 访问 类型引用
	VisitAttr(*Attr) Node           // 访问 属性
	VisitArray(*Array) Node         // 访问 数组
	VisitList(*List) Node           // 访问 链表
	VisitArgs(*Args) Node           // 访问 参数列表
	VisitNamedArgs(*NamedArgs) Node // 访问 命名参数列表
	VisitString(*String) Node       // 访问 字符串
	VisitFloat(*Float) Node         // 访问 浮点数
	VisitInt(*Int) Node             // 访问 整数
	VisitBool(*Bool) Node           // 访问 布尔值
	VisitBinaryOp(*BinaryOp) Node   // 访问 二元运算
	VisitMap(*Map) Node             // 访问 字典
}

// 访问者模式
// 为每一种节点类型构造一个Accept方法 使其能够实现Node接口
// 每一种节点的接受一个访问者参数 然后以自身为参数调用访问者对应自己类型的访问方法

// Accept 为二元运算实现 Node接口
func (node *BinaryOp) Accept(visitor Visitor) Node {
	return visitor.VisitBinaryOp(node)
}

// Accept 为 参数 实现 Node接口
func (node *Param) Accept(visitor Visitor) Node {
	return visitor.VisitParam(node)
}

// Accept 为 字面量字符串 实现 Node接口
func (node *String) Accept(visitor Visitor) Node {
	return visitor.VisitString(node)
}

// Accept 为 字面量浮点数 实现 Node接口
func (node *Float) Accept(visitor Visitor) Node {
	return visitor.VisitFloat(node)
}

// Accept 为 字面量整形数 实现 Node接口
func (node *Int) Accept(visitor Visitor) Node {
	return visitor.VisitInt(node)
}

// Accept 为 字面量布尔值 实现 Node接口
func (node *Bool) Accept(visitor Visitor) Node {
	return visitor.VisitBool(node)
}

// Accept 为 命名参数列表 实现 Node接口
func (node *NamedArgs) Accept(visitor Visitor) Node {
	return visitor.VisitNamedArgs(node)
}

// Accept 为 匿名参数列表 实现 Node接口
func (node *Args) Accept(visitor Visitor) Node {
	return visitor.VisitArgs(node)
}

// Accept 为 包节点 实现 Node接口
func (node *Package) Accept(visitor Visitor) Node {
	return visitor.VisitPackage(node)
}

// Accept 为 代码节点 实现 Node接口
func (node *Script) Accept(visitor Visitor) Node {
	return visitor.VisitScript(node)
}

// Accept 为 一组枚举 实现 Node接口
func (node *Enum) Accept(visitor Visitor) Node {
	return visitor.VisitEnum(node)
}

// Accept 为 单个枚举值 实现 Node接口
func (node *EnumVal) Accept(visitor Visitor) Node {
	return visitor.VisitEnumVal(node)
}

// Accept 为 表 实现 Node接口
func (node *Table) Accept(visitor Visitor) Node {
	return visitor.VisitTable(node)
}

// Accept 为 域 实现 Node接口
func (node *Field) Accept(visitor Visitor) Node {
	return visitor.VisitField(node)
}

// Accept 为 单个协议 实现 Node接口
func (node *Contract) Accept(visitor Visitor) Node {
	return visitor.VisitContract(node)
}

// Accept 为  函数 实现 Node接口
func (node *Method) Accept(visitor Visitor) Node {
	return visitor.VisitMethod(node)
}

// Accept 为 类型引用 实现 Node接口
func (node *TypeRef) Accept(visitor Visitor) Node {
	return visitor.VisitTypeRef(node)
}

// Accept 为 属性 实现 Node接口
func (node *Attr) Accept(visitor Visitor) Node {
	return visitor.VisitAttr(node)
}

// Accept 为 数组 实现 Node接口
func (node *Array) Accept(visitor Visitor) Node {
	return visitor.VisitArray(node)
}

// Accept 为 链表 实现 Node接口
func (node *List) Accept(visitor Visitor) Node {
	return visitor.VisitList(node)
}

// Accept 为 字典 实现Node接口
func (node *Map) Accept(visitor Visitor) Node {
	return visitor.VisitMap(node)
}

// EmptyVisitor 一个空的什么都不做的访问者
type EmptyVisitor struct{}

// VisitString 实现访问者接口
func (visitor *EmptyVisitor) VisitString(*String) Node {
	return nil
}

// VisitFloat 实现访问者接口
func (visitor *EmptyVisitor) VisitFloat(*Float) Node {
	return nil
}

// VisitInt 实现访问者接口
func (visitor *EmptyVisitor) VisitInt(*Int) Node {
	return nil
}

// VisitBool 实现访问者接口
func (visitor *EmptyVisitor) VisitBool(*Bool) Node {
	return nil
}

// VisitPackage 实现访问者接口
func (visitor *EmptyVisitor) VisitPackage(*Package) Node {
	return nil
}

// VisitScript 实现访问者接口
func (visitor *EmptyVisitor) VisitScript(*Script) Node {
	return nil
}

// VisitEnum 实现访问者接口
func (visitor *EmptyVisitor) VisitEnum(*Enum) Node {
	return nil
}

// VisitEnumVal 实现访问者接口
func (visitor *EmptyVisitor) VisitEnumVal(*EnumVal) Node {
	return nil
}

// VisitTable 实现访问者接口
func (visitor *EmptyVisitor) VisitTable(*Table) Node {
	return nil
}

// VisitField 实现访问者接口
func (visitor *EmptyVisitor) VisitField(*Field) Node {
	return nil
}

// VisitContract 实现访问者接口
func (visitor *EmptyVisitor) VisitContract(*Contract) Node {
	return nil
}

// VisitMethod 实现访问者接口
func (visitor *EmptyVisitor) VisitMethod(*Method) Node {

	return nil
}

// VisitTypeRef 实现访问者接口
func (visitor *EmptyVisitor) VisitTypeRef(*TypeRef) Node {
	return nil
}

// VisitAttr 实现访问者接口
func (visitor *EmptyVisitor) VisitAttr(*Attr) Node {

	return nil
}

// VisitArray 实现访问者接口
func (visitor *EmptyVisitor) VisitArray(*Array) Node {
	return nil
}

// VisitList 实现访问者接口
func (visitor *EmptyVisitor) VisitList(*List) Node {
	return nil
}

// VisitMap 实现访问者接口
func (visitor *EmptyVisitor) VisitMap(*Map) Node {
	return nil
}

// VisitArgs 实现访问者接口
func (visitor *EmptyVisitor) VisitArgs(*Args) Node {
	return nil
}

// VisitNamedArgs 实现访问者接口
func (visitor *EmptyVisitor) VisitNamedArgs(*NamedArgs) Node {
	return nil
}

// VisitParam 实现访问者接口
func (visitor *EmptyVisitor) VisitParam(*Param) Node {
	return nil
}

// VisitBinaryOp 实现访问者接口
func (visitor *EmptyVisitor) VisitBinaryOp(*BinaryOp) Node {
	return nil
}
