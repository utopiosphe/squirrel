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

func TestReduceColumnsInOneLine(t *testing.T) {
	// 测试从逗号分隔的列中删除单个列
	query := Select("id, name").From("users")
	reducedQuery := query.ReduceColumns("id")
	sql, _, err := reducedQuery.ToSql()
	if err != nil {
		t.Errorf("Error generating SQL: %v", err)
	}
	expectedSQL := "SELECT name FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	// 测试从逗号分隔的列中删除多个列
	query2 := Select("id, name, email").From("users")
	reducedQuery2 := query2.ReduceColumns("id", "email")
	sql2, _, err := reducedQuery2.ToSql()
	if err != nil {
		t.Errorf("Error generating SQL: %v", err)
	}
	expectedSQL2 := "SELECT name FROM users"
	if sql2 != expectedSQL2 {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL2, sql2)
	}

	// 测试删除所有列的情况
	query3 := Select("id, name").From("users")
	reducedQuery3 := query3.ReduceColumns("id", "name")
	_, _, err = reducedQuery3.ToSql()
	if err == nil {
		t.Errorf("Expected error when no columns remain, but got none")
	}

	// 测试混合情况：逗号分隔的列和独立列
	query4 := Select("id, name", "email", "age, created_at").From("users")
	reducedQuery4 := query4.ReduceColumns("id", "email", "age")
	sql4, _, err := reducedQuery4.ToSql()
	if err != nil {
		t.Errorf("Error generating SQL: %v", err)
	}
	expectedSQL4 := "SELECT name, created_at FROM users"
	if sql4 != expectedSQL4 {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL4, sql4)
	}
}
