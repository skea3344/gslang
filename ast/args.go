// @file 	args.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	args

package ast

// Args 参数列表节点
type Args struct {
	BaseExpr        // 内嵌基本表达式实现
	Items    []Expr // 参数列表
}

// NewArgs 在代码内新建参数列表节点 此参数列表所属代码节点为此代码节点
func (node *Script) NewArgs() *Args {
	expr := &Args{}
	expr.Init("args", node)
	return expr
}

// NewArg 在参数列表节点内 保存对应表达式对应的参数 此参数的父节点为此参数列表节点
func (node *Args) NewArg(arg Expr) Expr {
	// 添加到参数列表
	node.Items = append(node.Items, arg)
	// 设置参数的父节点
	arg.SetParent(node)
	return arg
}

// NamedArgs 命名参数列表节点
type NamedArgs struct {
	BaseExpr                 // 内嵌基本表达式实现
	Items    map[string]Expr // 用字典保存命名参数列表
}

// NewNamedArgs 在代码节点内新建命名参数列表 此命名参数列表名字args 所属代码节点为此代码节点
func (node *Script) NewNamedArgs() *NamedArgs {
	expr := &NamedArgs{}
	expr.Init("args", node)
	return expr
}

// NewArg 用指定的名字和表达式在命名参数列表内添加参数 并返回此参数表达式和添加结果
func (node *NamedArgs) NewArg(name string, arg Expr) (Expr, bool) {
	// 先检查是否有同名参数 有则返回此参数 及 新建失败标志
	if item, ok := node.Items[name]; ok {
		return item, false
	}
	node.Items[name] = arg
	// 设置此参数的父节点为此命名参数列表
	arg.SetParent(node)
	return arg, true
}
