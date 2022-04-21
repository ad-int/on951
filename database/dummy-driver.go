package database

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
)

type DummyDialector struct{}

func (dialector DummyDialector) Name() string {
	return "kibiras"
}

func (dialector DummyDialector) Initialize(db *gorm.DB) error {
	return errors.New(DbConnectionError)
}
func (dialector DummyDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}
func (dialector DummyDialector) DataTypeOf(field *schema.Field) string {
	return reflect.TypeOf(field).String()
}

func (dialector DummyDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return nil
}

func (dialector DummyDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {

}

func (dialector DummyDialector) QuoteTo(writer clause.Writer, s string) {
}

func (dialector DummyDialector) Explain(sql string, vars ...interface{}) string {
	return sql
}
