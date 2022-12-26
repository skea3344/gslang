// @file 	contract.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	contract

package ast

import (
	"fmt"

	"github.com/skea3344/gserrors"
)

// Param 函数参数节点
type Param struct {
	BaseExpr // 内嵌基本表达式实现
	ID       int
	Type     Expr
}

// Method 函数节点
type Method struct {
	BaseExpr          // 内嵌基本表达式实现
	ID       uint16   // 函数ID
	Return   []*Param // 返回参数列表
	Params   []*Param // 输入参数列表
}

// InputParams 函数输入参数个数
func (method *Method) InputParams() uint16 {
	return uint16(len(method.Params))
}

// ReturnParams 函数返回参数个数
func (method *Method) ReturnParams() uint16 {
	return uint16(len(method.Return))
}

// NewReturn 在函数节点上新建返回参数 并加入到此函数返回参数列表
func (method *Method) NewReturn(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前函数返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Return),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此函数节点所属的代码节点
	param.Init(fmt.Sprintf("return_arg(%d)", param.ID), method.Script())
	// 参数节点的父节点为此函数节点
	param.SetParent(method)
	// 加入到此函数返回参数列表
	method.Return = append(method.Return, param)
	return param
}

// NewParam 在函数节点上新建输入参数 并加入到此函数输入参数列表
func (method *Method) NewParam(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前函数返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Params),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此函数节点所属的代码节点
	param.Init(fmt.Sprintf("arg(%d)", param.ID), method.Script())
	// 参数节点的父节点为此函数节点
	param.SetParent(method)
	// 加入到此函数返回参数列表
	method.Params = append(method.Params, param)
	return param
}

// Contract 协议节点 一个协议内包含有多个函数节点
type Contract struct {
	BaseExpr                    // 内嵌基本表达式实现
	Methods  map[string]*Method // 协议内函数列表
	Bases    []*TypeRef         // 父协议列表
}

// NewContract 在代码节点内新建协议节点
func (node *Script) NewContract(name string) (expr *Contract) {
	// 确保协议的函数列表被初始化
	defer gserrors.Ensure(func() bool {
		return expr.Methods != nil
	}, "make sure alloc Contract's Methods field")
	expr = &Contract{
		Methods: make(map[string]*Method),
	}
	// 设置协议节点为给定名字 设置所属代码节点为此代码节点
	expr.Init(name, node)
	return expr
}

// NewBase 为此协议节点添加指定的类型引用
func (expr *Contract) NewBase(base *TypeRef) (ref *TypeRef, ok bool) {
	// 如果类型引用列表已存在同名引用 则返回
	for _, old := range expr.Bases {
		if base.Name() == old.Name() {
			ref = old
			return
		}
	}
	// 设置此类型引用的节点为此协议节点
	base.SetParent(expr)
	// 添加此类型引用到协议的类型引用列表
	expr.Bases = append(expr.Bases, base)
	ref = base
	ok = true
	return
}

// NewMethod 在协议内生成函数
func (expr *Contract) NewMethod(name string) (method *Method, ok bool) {
	// 如果已有同名函数 则直接返回
	if method, ok = expr.Methods[name]; ok {
		ok = false
		return
	}
	// 根据协议内函数列表长度获得函数ID
	method = &Method{
		ID: uint16(len(expr.Methods)),
		// Params: make([]*Param),
		// Return: make([]*Param),
	}
	// 设置函数给给定名字 设置所属代码节点为 协议所属的代码节点
	method.Init(name, expr.Script())
	// 设置函数的父节点为此协议节点
	method.SetParent(expr)
	// 将函数添加到协议的函数列表
	expr.Methods[name] = method
	ok = true
	return
}
