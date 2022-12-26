// @file 	enum.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	enum

package ast

import (
	"github.com/skea3344/gserrors"
)

// EnumVal 枚举值 指一个枚举括号中的单个枚举值
type EnumVal struct {
	BaseExpr       // 内嵌基本表达式实现
	Value    int64 //  枚举值节点对应实际枚举数值
}

// Enum 枚举 指一个枚举声明中的所有内容
type Enum struct {
	BaseExpr                     // 内嵌基本表达式实现
	Values   map[string]*EnumVal // 枚举值表
	Default  *EnumVal            // 入口枚举值
	Length   uint                // 枚举类型长度
	Signed   bool                // 枚举值是否有符号
}

// NewEnum 在代码内新建枚举 此枚举节点的父节点为此代码节点
func (node *Script) NewEnum(name string, length uint, signed bool) (expr *Enum) {
	// 枚举类型的长度仅支持1,2,4字节 TODO
	gserrors.Require(func() bool {
		switch length {
		case 1, 2, 4:
			return true
		default:
			return false
		}
	}(), "the enum type length can only be 1,2,4,got :%d", length)
	// 确保生成的Enum对象的Values值要被初始化 不能为nil
	defer gserrors.Ensure(func() bool {
		return expr.Values != nil
	}, "make sure alloc Enum's Values field")
	expr = &Enum{
		Values: make(map[string]*EnumVal),
		Length: length,
		Signed: signed,
	}
	// 设置枚举名字 设置所属代码节点
	expr.Init(name, node)
	// 设置父节点 为此代码节点
	expr.SetParent(node)
	return expr
}

// NewVal 在枚举内生成一个枚举值
func (node *Enum) NewVal(name string, val int64) (result *EnumVal, ok bool) {
	defer gserrors.Ensure(func() bool {
		return ok == (node.Values[name] == result)
	}, "post condition check")
	// 检查枚举中是否已有同名枚举值 有则直接返回
	if result, ok = node.Values[name]; ok {
		ok = !ok
		return
	}
	// 新建枚举值
	result = &EnumVal{
		Value: val,
	}
	// 指定枚举值的名字及所属代码节点
	result.Init(name, node.Script())
	// 加入枚举节点枚举值字典
	node.Values[name] = result
	ok = true
	// 如果枚举还没有默认入口值 则将此值设为默认值
	if node.Default == nil {
		node.Default = result
	}
	return
}
