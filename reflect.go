// @file 	reflect.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	reflect

package gslang

import (
	"github.com/skea3344/gserrors"
	"github.com/skea3344/gslang/ast"
)

// Enum 将Enum节点内的EnumVal字典解析对应成 名字和int64值的字典
func Enum(enum *ast.Enum) map[string]int64 {
	values := make(map[string]int64)
	for _, val := range enum.Values {
		values[val.Name()] = val.Value
	}
	return values
}

// EvalFieldInitArg 在参数列表中找到与指定域对应的参数表达式并返回 该表达式及查找结果
func EvalFieldInitArg(field *ast.Field, expr ast.Expr) (ast.Expr, bool) {
	eval := &evalArg{
		field: field,
	}
	expr.Accept(eval)
	if eval.expr != nil {
		return eval.expr, true
	}
	return nil, false
}

// EvalEnumVal 访问枚举值节点的val值
func EvalEnumVal(expr ast.Expr) int64 {
	visitor := &evalEnumVal{}
	expr.Accept(visitor)
	return visitor.val
}

// IsAttrUsage 判断是不是内置AttrUsage结构
func IsAttrUsage(expr *ast.Table) bool {
	if expr.Package().Name() == GSLangPackage && expr.Name() == "AttrUsage" {
		return true
	}
	return false
}

// IsStruct 检查表节点是不是表示一个结构体
func IsStruct(expr *ast.Table) bool {
	_, ok := expr.Extra("isStruct")
	return ok
}

// markAsStruct 设置一个表节点表示一个结构体
func markAsStruct(expr *ast.Table) {
	expr.NewExtra("isStruct", true)
}

// IsError 检查枚举是不是表示错误声明
func IsError(expr *ast.Enum) bool {
	_, ok := expr.Extra("isError")
	return ok
}

// markAsError 设置一个枚举节点表示一组错误声明
func markAsError(expr *ast.Enum) {
	expr.NewExtra("isError", true)
}

// EvalAttrUsage 执行属性
func (cs *CompileS) EvalAttrUsage(attr *ast.Attr) int64 {
	// 属性的类型引用必须先连接到对应类型
	gserrors.Require(attr.Type.Ref != nil, "attr(%s) must linked first :\n\t%s", attr, Pos(attr))
	// 只有Table才能被作为属性的类型引用
	table, ok := attr.Type.Ref.(*ast.Table)
	if !ok {
		gserrors.Panicf(
			ErrCompileS,
			"only table can be used as attribute type :\n\tattr def: %s\n\ttype def: %s", Pos(attr), Pos(attr.Type.Ref))
	}
	// 轮询属性的类型引用的属性列表
	for _, metattr := range table.Attrs() {
		// 属性的类型引用必须是Table
		usage, ok := metattr.Type.Ref.(*ast.Table)
		gserrors.Require(ok, "attr(%s) must linked first :\n\t%s", metattr, Pos(attr))
		// 如果是内置的AttrUsage
		if IsAttrUsage(usage) {
			// 取AttrUsage的Target域且必须有
			field, ok := usage.Field("Target")
			if !ok {
				gserrors.Panicf(ErrCompileS, "inner error: gslang AttrUsage must declare Target Field \n\ttype def:%s", Pos(usage))
			}
			// 在属性的参数列表中查找对应域的参数并求值并返回
			if target, ok := EvalFieldInitArg(field, metattr.Args); ok {
				return EvalEnumVal(target)
			}
			// AttrUsage 属性必须有
			gserrors.Panicf(ErrCompileS, "AttrUsage attribute initlist expect target val \n\tattr def:%s", Pos(metattr))
		}
	}
	// 能作为属性的Table必须有一个属性@AttrUsage
	gserrors.Panicf(ErrCompileS, "target table can no be used as attribute type:\n\tattr def: %s\n\ttype def: %s", Pos(attr), Pos(attr.Type.Ref))
	return 0
}
