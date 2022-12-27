// @file 	position.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	position

package gslang

import (
	"fmt"
	"path/filepath"
)

// Position 源码文件中的具体位置
type Position struct {
	Filename string // 文件名
	Line     int    // 行号 从1开始
	Column   int    // 列号 从1开始
}

// ShortName 返回相对文件名
func (pos Position) ShortName() string {
	return filepath.Base(pos.Filename)
}

// String 节点位置的字符串显示
func (pos Position) String() string {
	return fmt.Sprintf("%s(%d:%d)", pos.Filename, pos.Line, pos.Column)
}

// Valid 检查节点位置是否有效
func (pos Position) Valid() bool {
	return pos.Line != 0
}
