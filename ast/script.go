// @file 	package.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	script

package ast

import "github.com/skea3344/gserrors"

// Script 代码节点 代表一个源码文件
type Script struct {
	BaseNode                        // 内嵌基本节点实现 , 注意非基本表达式
	Imports  map[string]*PackageRef // 引用的代码包
	Types    []Expr                 // 包内类型 声明 用于防止重复声明类型
	pkg      *Package
}

// NewScript 在包节点内新建一个代码节点
func (node *Package) NewScript(name string) (script *Script, err error) {
	if old, ok := node.Scripts[name]; ok {
		script = old
		err = gserrors.Newf(ErrAst, "duplicate script named:%s", old.Name())
		return
	}
	defer gserrors.Ensure(func() bool {
		return script.pkg != nil
	}, "make sure set the pkg field")
	defer gserrors.Ensure(func() bool {
		return script.Imports != nil
	}, "make sure set the Imports field")
	script = &Script{
		pkg:     node,
		Imports: make(map[string]*PackageRef),
	}
	// 初始化代码节点 设置包节点为代码节点的父节点
	script.Init(name, node)
	// 将代码节点加入到包节点的代码列表
	node.Scripts[name] = script
	return
}

// PackageRef 包引用节点 代表一个源代码文件中引用的其他包
type PackageRef struct {
	BaseNode
	Ref *Package
}

// NewPackageRef 在代码节点中新建一个包引用节点
func (node *Script) NewPackageRef(name string, pkg *Package) (ref *PackageRef, ok bool) {
	// 检查已引用的包列表内是否有同名包 有的则返回该同名包 并设置新引用失败
	if ref, ok = node.Imports[name]; ok {
		ok = false
		return
	}
	// 新建包引用
	ref = &PackageRef{
		Ref: pkg,
	}
	// 设置包引用 名字 设置父节点为此代码节点
	ref.Init(name, node)
	// 将包引用加入到代码节点的包引用列表
	node.Imports[name] = ref
	return ref, true
}

// NewType 在代码节点中新建一个类型节点,类型节点在代码节点所属包节点中唯一. 包和代码节点分别以字典和切片保存此类型节点的引用
func (node *Script) NewType(expr Expr) (old Expr, ok bool) {
	old, ok = node.pkg.NewType(expr)
	if ok {
		node.Types = append(node.Types, expr)
		expr.SetParent(node)
	}
	return
}

// Package 获取代码节点所属的包节点
func (node *Script) Package() *Package {
	return node.pkg
}
