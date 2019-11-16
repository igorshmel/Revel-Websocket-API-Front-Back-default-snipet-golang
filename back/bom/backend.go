package bom

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"

	"stixi/back/hlp"
)

/////////////////////////////////////////////////////////////////////////
type requestBackend struct {
	request RequestMessage
	errno   int
}

func Backend() *requestBackend {
	out := requestBackend{}
	return &out
}

func (obj *requestBackend) DoIt(req []byte, rpl interface{}, query string, DB *gorm.DB) {

	if len(query) < 50 {
		if fmt.Sprint(reflect.TypeOf(rpl).Kind()) == "ptr" {

			if err := json.Unmarshal(req, &obj.request.Message); err == nil {
				if registered, ok := regMess[query]; ok {
					if message := registered.CreateRequest(); message != nil {
						if err := json.Unmarshal(obj.request.Message, message); err == nil {
							if reply, err := message.DoIt(DB); err == nil {
								if reply != nil {
									if err := json.Unmarshal(hlp.AnyToByte(reply), rpl); err != nil {
										fmt.Println("err: unable unmarshal rpl")
									}
								} else {
									fmt.Println("err: reply is nul")
								}
							} else {
								fmt.Println("err: DoIt: " + fmt.Sprint(err))
							}
						} else {
							fmt.Println("err: unable unmarshal message")
						}
					} else {
						fmt.Println("err: unable to create interface for request")
					}
				} else {
					fmt.Println("err: unable registeredMessages query")
				}
			} else {
				fmt.Println("err: unable unmarshal req")
			}
		} else {
			fmt.Println("err: rpl is not ptr")
		}
	} else {
		fmt.Println("err: query is too long")
	}
}
