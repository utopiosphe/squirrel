package squirrel

import (
	"testing"
)

func TestReduceColumns(t *testing.T) {
	// 创建一个包含多个列的查询
	query := Select("id", "name", "email", "age", "created_at").
		From("users")

	// 验证初始列
	columns := query.GetSelectColumns()
	expected := []string{"id", "name", "email", "age", "created_at"}
	if len(columns) != len(expected) {
		t.Errorf("Expected %d columns, got %d", len(expected), len(columns))
	}

	// 删除指定的列
	reducedQuery := query.ReduceColumns("email", "age")

	// 验证删除后的列
	reducedColumns := reducedQuery.GetSelectColumns()
	expectedReduced := []string{"id", "name", "created_at"}
	if len(reducedColumns) != len(expectedReduced) {
		t.Errorf("Expected %d columns after reduction, got %d", len(expectedReduced), len(reducedColumns))
	}

	// 验证 SQL 生成
	sql, _, err := reducedQuery.ToSql()
	if err != nil {
		t.Errorf("Error generating SQL: %v", err)
	}
	expectedSQL := "SELECT id, name, created_at FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestReduceColumnsEmpty(t *testing.T) {
	// 测试空参数的情况
	query := Select("id", "name").From("users")
	reducedQuery := query.ReduceColumns()
	
	// 应该返回原始查询
	sql1, _, _ := query.ToSql()
	sql2, _, _ := reducedQuery.ToSql()
	if sql1 != sql2 {
		t.Errorf("Empty ReduceColumns should return original query")
	}
}

func TestReduceColumnsNonExistent(t *testing.T) {
	// 测试删除不存在的列
	query := Select("id", "name").From("users")
	reducedQuery := query.ReduceColumns("non_existent_column")
	
	// 应该返回原始查询
	sql1, _, _ := query.ToSql()
	sql2, _, _ := reducedQuery.ToSql()
	if sql1 != sql2 {
		t.Errorf("Reducing non-existent columns should return original query")
	}
} 