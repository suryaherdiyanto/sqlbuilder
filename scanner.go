package sqlbuilder

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

func ScanStruct(d interface{}, rows *sql.Rows) error {
	columns, err := rows.Columns()

	ref := reflect.TypeOf(d)
	if ref.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("The destination must be pointer, passed: %v", ref.Kind()))
	}

	ref = ref.Elem()

	refKind := ref.Kind()
	if refKind != reflect.Struct {
		return errors.New(fmt.Sprintf("model.ScanRow only accepts struct kind: %v", refKind))
	}

	if refKind == reflect.Struct {
		var dRefs []interface{}
		for _, column := range columns {
			for i := 0; i < ref.NumField(); i++ {
				f := ref.Field(i)

				if f.Tag.Get("db") == column {
					dRefs = append(dRefs, reflect.ValueOf(d).Elem().Field(i).Addr().Interface())
				}
			}
		}
		if err = rows.Scan(dRefs...); err != nil {
			return err
		}
	}

	return nil

}

func ScanAll(d interface{}, rows *sql.Rows) error {
	ref := reflect.TypeOf(d)
	val := reflect.ValueOf(d)

	if ref.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("The destination must be a pointer, passed: %v", ref.Kind()))
	}
	ref = ref.Elem()

	if val.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("The destination must be a pointer, passed: %v", val.Kind()))
	}
	val = val.Elem()

	if ref.Kind() != reflect.Slice {
		return errors.New(fmt.Sprintf("model.ScanAll only accepts slice kind: %v", ref.Kind()))
	}

	base := ref.Elem()
	val.SetLen(0)

	for rows.Next() {
		if base.Kind() == reflect.Struct {
			v := reflect.New(base)
			if err := ScanStruct(v.Interface(), rows); err != nil {
				return err
			}
			val.Set(reflect.Append(val, reflect.Indirect(v)))
		}

		if base.Kind() == reflect.Map {
			var m = make(map[string]interface{})
			if err := ScanMap(&m, rows); err != nil {
				return err
			}
			val.Set(reflect.Append(val, reflect.ValueOf(m)))
		}
	}

	return nil
}

func ScanMap(d interface{}, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	ref := reflect.TypeOf(d)

	if ref.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("The destination must be pointer, passed: %v", ref.Kind()))
	}

	ref = ref.Elem()

	refKind := ref.Kind()

	if refKind != reflect.Map {
		return errors.New(fmt.Sprintf("model.ScanMap only accepts map kind: %v", refKind))
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}
	if err := rows.Scan(values...); err != nil {
		return err
	}

	for i, column := range columns {
		value := reflect.ValueOf(values[i]).Elem()
		switch value.Elem().Kind() {
		case reflect.Int64:
			reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(int(value.Interface().(int64))))
		case reflect.Float64:
			reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(float32(value.Interface().(float64))))
		case reflect.Bool:
			reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(value.Interface().(bool)))
		case reflect.Invalid:
			reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(""))
		default:
			if value.Elem().Kind() != reflect.Slice {
				if value.Elem().Type().Name() == "Time" {
					reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(value.Interface().(time.Time)))
					continue
				}
			}
			reflect.ValueOf(d).Elem().SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(string(value.Interface().([]byte))))
		}
	}

	return nil
}
