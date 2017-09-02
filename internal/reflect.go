package internal

import "reflect"

func ReflectName(a interface{}) string {
	iType := reflect.TypeOf(a).Elem()
	iName := iType.Name() + iType.PkgPath()
	return iName
}
