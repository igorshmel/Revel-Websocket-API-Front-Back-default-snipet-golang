package api

import (
	"encoding/json"
	"stixi/back/bom"
	"stixi/back/hlp"

	"github.com/jinzhu/gorm"
)

func init() {
	request := "write"
	bom.Register(request, (*ReqWrite)(nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type TxtList struct {
	Id    int32  `gorm:"Column:id"`
	Title string `gorm:"Column:title"`
}

type ReqWrite struct {
	D []byte
}

func (obj *ReqWrite) DoIt(DB *gorm.DB) (bom.Reply, error) {

	var err error
	var txtList TxtList
	var dat map[string]interface{}

	out := RplWrite{}

	json.Unmarshal(obj.D, &dat)

	title := dat["title"].(string)

///////////////////////////////////////////
// Заглушка для логики записи переменных //
	txtList.Title = title
///////////////////////////////////////////

	out.Rpl = hlp.JsonMarshal(txtList)

	return &out, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type RplWrite struct {
	Rpl []byte
	Err error
}
