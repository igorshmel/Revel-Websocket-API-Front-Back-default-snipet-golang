package api

import (
	"github.com/jinzhu/gorm"
	"stixi/back/bom"
)

func init() {
	request := "apidefault"
	bom.Register(request, (*ReqApiDefaultStruct)(nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type ReqApiDefaultStruct struct {
	D []byte
}

func (obj *ReqApiDefaultStruct) DoIt(DB *gorm.DB) (bom.Reply, error) {

	var err error
	out := RplApiDefaultStruct{}

	return &out, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type RplApiDefaultStruct struct {
	Rpl []byte
	Err error
}
