package bom

import (
	"fmt"
	"reflect"
	"strings"
)

//////////////////////////////////////////////////////////////////////////
var regMess map[string]requestRefType
var regReqType map[string]string
var reqInterfaceType = reflect.TypeOf((*Request)(nil)).Elem()

func init() {
	regReqType = make(map[string]string, 10)
	regMess = make(map[string]requestRefType, 10)
}

///////////////////////////////////////////////////////////////////////////
type requestRefType struct {
	request reflect.Type
}

func (obj *requestRefType) CreateRequest() Request {
	return reflect.New(obj.request).Interface().(Request)
}

//////////////////////////////////////////////////////////////////////////
func Register(req string, reqType Request) {

	refTypeOf := reflect.TypeOf(reqType)
	reqRefTypeStruct := requestRefType{}

	if refTypeOf.Kind() != reflect.Ptr || !refTypeOf.Implements(reqInterfaceType) || refTypeOf.Elem().Kind() != reflect.Struct {
		fmt.Println("err: Unable to register Request.")
	}

	reqRefTypeStruct.request = refTypeOf.Elem()

	regMess[strings.ToLower(req)] = reqRefTypeStruct
	regReqType[refTypeOf.String()] = req
}
