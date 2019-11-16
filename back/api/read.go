package api

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"stixi/back/bom"
	"stixi/back/hlp"
)

func init() {
	request := "read"
	bom.Register(request, (*ReqRead)(nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Txt struct {
	Title    string `gorm:"Column:title"`
	Body     string `gorm:"Column:body"`
	AuthorID int16  `gorm:"Column:author"`
	NickName string `gorm:"Column:nickname"`
}

type ReqRead struct {
	D []byte
}

func (obj *ReqRead) DoIt(DB *gorm.DB) (bom.Reply, error) {

	var txt Txt
	var err error
	var dat map[string]interface{}

	out := RplRead{}

	json.Unmarshal(obj.D, &dat)

	title := dat["title"].(string)

///////////////////////////////////////////
// Заглушка для логики чтения переменных //
	txt.Title = title
///////////////////////////////////////////

	out.Rpl = hlp.JsonMarshal(txt)

	return &out, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type RplRead struct {
	Rpl []byte
	Err error
}
