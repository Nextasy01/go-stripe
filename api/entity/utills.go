package entity

import (
	"reflect"
)

func ConvertList(src, dest interface{}) error {

	srcValue := reflect.ValueOf(src).Elem()
	destValue := reflect.ValueOf(dest).Elem()

	destValue.Set(reflect.MakeSlice(destValue.Type(), srcValue.Len(), srcValue.Cap()))

	for i := 0; i < srcValue.Len(); i++ {
		srcElem := srcValue.Index(i)
		destElem := destValue.Index(i)

		for j := 0; j < srcElem.NumField(); j++ {
			destField := destElem.FieldByName(srcElem.Type().Field(j).Name)
			srcField := srcElem.Field(j)

			if destField.IsValid() && destField.CanSet() && srcField.IsValid() {
				destField.Set(srcField)
			}
		}
	}

	return nil
}

func Convert(src, dest interface{}) error {
	destValue := reflect.ValueOf(dest).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		destField := destValue.FieldByName(srcValue.Type().Field(i).Name)
		srcField := srcValue.Field(i)

		if destField.IsValid() && destField.CanSet() && srcField.IsValid() {
			destField.Set(srcField)
		}
	}

	return nil
}
