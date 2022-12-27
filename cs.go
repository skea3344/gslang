// @file 	cs.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	complieS

package gslang

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skea3344/gserrors"
	"github.com/skea3344/gslang/ast"
	"github.com/skea3344/logger"
)

var (
	// ErrCompileS 编译模块错误
	ErrCompileS = errors.New("CompileS error")
)

// setFilePath 设置代码节点的 绝对文件名
func setFilePath(script *ast.Script, fullPath string) {
	script.NewExtra("FilePath", fullPath)
}

// FilePath 返回代码节点的绝对文件名
func FilePath(script *ast.Script) (string, bool) {
	path, ok := script.Extra("FilePath")
	if ok {
		return path.(string), ok
	}
	return "", false
}

// CompileS 编译器
type CompileS struct {
	logger.ILog                         // 内嵌通用日志接口
	Loaded      map[string]*ast.Package // 已加载包节点字典
	loading     []*ast.Package          // 正在加载的包节点列表
	goPath      []string                // 系统golang路径
}

// NewCompileS 新建一个编译器
func NewCompileS() *CompileS {
	GOPATH := os.Getenv("GOPATH")
	if GOPATH == "" {
		gserrors.Panicf(ErrCompileS, "must set GOPATH first")
	}
	return &CompileS{
		ILog:   logger.Get("gslang"),
		Loaded: make(map[string]*ast.Package),
		goPath: strings.Split(GOPATH, string(os.PathListSeparator)),
	}
}

// searchPackage 从系统 $GOPATH/src下查找指定名字的代码包 需要唯一
func (cs *CompileS) searchPackage(packageName string) string {
	var found []string
	for _, path := range cs.goPath {
		fullpath := filepath.Join(path, "src", packageName)
		fi, err := os.Stat(fullpath)
		if err == nil && fi.IsDir() {
			found = append(found, fullpath)
		}
	}
	// 多于1个或者少于1个包均报错
	if len(found) < 1 {
		gserrors.Panicf(ErrCompileS, "found no package named %s", packageName)
	}
	if len(found) > 1 {
		var buff bytes.Buffer
		buff.WriteString(fmt.Sprintf("found more than one package named:%s", packageName))
		for i, path := range found {
			buff.WriteString(fmt.Sprintf("\n\t%d)%s", i, path))
		}
		gserrors.Panicf(ErrCompileS, buff.String())
	}
	// 返回目标代码包绝对路径
	return found[0]
}

// circularRefCheck 循环引用检查 指定名字的包
func (cs *CompileS) circularRefCheck(packageName string) {
	var buff bytes.Buffer
	// 如果当前正在loading的包中包含对应的包名 那么就认为循环引用
	for _, pkg := range cs.loading {
		if pkg.Name() == packageName || buff.Len() != 0 {
			buff.WriteString(fmt.Sprintf("\t%s import\n", pkg.Name()))
		}
	}
	if buff.Len() != 0 {
		panic(fmt.Errorf("circular package import :\n%s\t%s", buff.String(), packageName))
	}
}

// errorf 编译器报错
func (cs *CompileS) errorf(position Position, fmtstring string, args ...interface{}) {
	gserrors.Panicf(
		ErrParse,
		fmt.Sprintf(
			"parse %s error : %s",
			position,
			fmt.Sprintf(fmtstring, args...),
		),
	)
}

// Accept 实现访问者模式  编译器节点访问入口
func (cs *CompileS) Accept(visitor ast.Visitor) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(gserrors.GSError); ok {
				err = e.(error)
			} else {
				err = gserrors.New(e.(error))
			}
		}
	}()
	// 使用访问者对编译器已加载的包轮流进行访问
	for _, pkg := range cs.Loaded {
		cs.D("%v", pkg.Name())
		pkg.Accept(visitor)
	}
	return
}

// Compile 对指定包名进行编译
func (cs *CompileS) Compile(packageName string) (pkg *ast.Package, err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(gserrors.GSError); ok {
				err = e.(error)
			} else {
				err = gserrors.New(e.(error))
			}
		}
	}()
	defer gserrors.Ensure(func() bool {
		if err == nil {
			return pkg != nil
		}
		return true
	}, "if err == nil the return param pkg can not be nil")
	if loaded, ok := cs.Loaded[packageName]; ok {
		pkg = loaded
		return
	}
	// 循环应用检测 在当前loading的包中已存在同名包 则报错
	cs.circularRefCheck(packageName)
	// 在系统中查找对应的包路径
	cs.D("%s", packageName)
	fullPath := cs.searchPackage(packageName)
	// 生成一个抽象包节点
	pkg = ast.NewPackage(packageName)
	// 将包添加到loading列表
	cs.loading = append(cs.loading, pkg)
	// 遍历目标包目录下的每一个文件
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		// 系统遍历时报错则直接返回该错误
		if err != nil {
			return err
		}
		// 如果该文件是一个不同于fullPath的文件夹 则略过
		if info.IsDir() && path != fullPath {
			return filepath.SkipDir
		}
		// 如果不是gs文件 则忽略
		if filepath.Ext(path) != ".gs" {
			return nil
		}
		// 解析该gs文件 生成代码节点
		script, err := cs.parse(pkg, path)
		if err == nil { // 没有错误的话 则将路径保存为代码节点的额外信息
			setFilePath(script, path)
		}
		return err
	})
	// 如果有错误发生 则将刚才加入loading的包去掉
	if err != nil {
		cs.loading = cs.loading[:len(cs.loading)-1]
		return
	}
	cs.link(pkg)
	// 加载完成 将loading中的pkg移到Loaded字典
	cs.loading = cs.loading[:len(cs.loading)-1]
	cs.Loaded[packageName] = pkg

	return
}

// Type 在当前编译器已加载的指定名字包中查找指定名字的类型表达式
func (cs *CompileS) Type(packageName string, typeName string) (ast.Expr, error) {
	if pkg, ok := cs.Loaded[packageName]; ok {
		if target, ok := pkg.Types[typeName]; ok {
			return target, nil
		}
		return nil, gserrors.Newf(ErrCompileS, "can not found type(%s) in package(%s)", typeName, packageName)
	}
	return nil, gserrors.Newf(ErrCompileS, "can not found package(%s)", packageName)
}
