// @file 	table.go
// @author 	caibo
// @email 	caibo923@gmail.com
// @desc 	table

package ast

// Field 域 节点
type Field struct {
	BaseExpr        //
	ID       uint16 // ID
	Type     Expr   // 类型表达式
}

// Table 表 节点
type Table struct {
	BaseExpr          //
	Fields   []*Field // 表的域列表
}

// NewTable 在代码节点内新建表
func (node *Script) NewTable(name string) (expr *Table) {
	expr = &Table{}
	// 设置表节点为给定的名字 设置所属代码节点
	expr.Init(name, node)
	return expr
}

// Field 在表内查找给定名字的域  返回该域和是否找到
func (expr *Table) Field(name string) (*Field, bool) {
	for _, field := range expr.Fields {
		if field.Name() == name {
			return field, true
		}
	}
	return nil, false
}

// NewField 在表内新建域
func (expr *Table) NewField(name string) (*Field, bool) {
	// 如果已存在同名域则直接返回
	for _, field := range expr.Fields {
		if field.Name() == name {
			return field, false
		}
	}
	// 新建域 ID为表的当前域列表长度
	field := &Field{
		ID: uint16(len(expr.Fields)),
	}
	// 设置名字 设置所属代码为 所属表的所属代码节点
	field.Init(name, expr.Script())
	// 设置父节点为此表节点
	field.SetParent(expr)
	// 将域添加到表的域列表
	expr.Fields = append(expr.Fields, field)
	return field, true
}
