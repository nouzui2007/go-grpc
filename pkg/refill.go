package pkg

import (
	"goexample/configs"
	"reflect"
	"time"
)

//Refill は i(元)をr(先)に詰め替えます。
//
//iとrに同じfield名があるものをコピーします。
//
//refill tag が"-"となっている場合はスルーします。
func Refill(i, r interface{}) interface{} {

	v := reflect.Indirect(reflect.ValueOf(i))
	result := reflect.Indirect(reflect.ValueOf(r))
	t := result.Type()

	for i := 0; i < result.NumField(); i++ {
		if t.Field(i).Tag.Get("refill") == "-" {
			continue
		}
		fName := t.Field(i).Name

		if result.Field(i).CanSet() && v.FieldByName(fName).IsValid() {
			result.Field(i).Set(v.FieldByName(fName))
		}
	}

	return r
}

//StructToMap は構造体をmapに変換します。(利用していないので削除)
func StructToMap(s interface{}) (m map[string]interface{}) {

	m = map[string]interface{}{}

	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fName := t.Field(i).Name
		m[fName] = v.FieldByName(fName)
	}
	return
}

//DereferenceIfPtr は利用していないので削除
func DereferenceIfPtr(value interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(value)).Interface()
}

//StrToSPtr はstringを stringポインタに変換します。
func StrToSPtr(s string) *string {
	return &s
}

//SPtrToStr はstringポインタを stringに変換します。
func SPtrToStr(sp *string) string {
	return reflect.Indirect(reflect.ValueOf(sp)).String()
}

//IntToIPtr はint32を int32ポインタに変換します。
func IntToIPtr(i int32) *int32 {
	return &i
}

//IPtrToInt はint32ポインタを int32に変換します。
func IPtrToInt(ip *int32) int32 {
	v := reflect.Indirect(reflect.ValueOf(ip)).Int()
	return int32(v)
}

//FltToFPtr はfloat32を float32ポインタに変換します。
func FltToFPtr(i float32) *float32 {
	return &i
}

//FPtrToFlt はfloat32ポインタを float32に変換します。
func FPtrToFlt(ip *float32) float32 {
	v := reflect.Indirect(reflect.ValueOf(ip)).Float()
	return float32(v)
}

//TimeToTPtr はtimeを timeポインタに変換します。
func TimeToTPtr(t time.Time) *time.Time {
	return &t
}

//TPtrToTime はtimeポインタを timeに変換します。
func TPtrToTime(t *time.Time) time.Time {
	return reflect.Indirect(reflect.ValueOf(t)).Interface().(time.Time)
}

//CheckTime は2006/01/02 15:04:05にパースします。
//
//todo(このメソッドいる？)
func CheckTime(s string) (t time.Time, err error) {
	t, err = time.Parse(configs.FormatDatetime, s)
	return
}
