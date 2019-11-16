package controllers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

func Add2Log(who string, what string, options ...string) {
	var opt, comma, whats string
	var mu sync.Mutex
	var err error
	path := "log"

	if len(options) > 0 {
		opt = " ("
		for _, op := range options {
			if len(opt) > 2 {
				comma = "; "
			}
			if len(op) > 0 {
				opt += comma + op
			}
		}
		opt += ")"
	}

	whats = strings.Replace(what, ":b:", "\033[97m", 1)
	whats = strings.Replace(whats, ":-:", "\033[0m", 1)
	what = strings.Replace(what, ":b:", "", 1)
	what = strings.Replace(what, ":-:", "", 1)

	if err = os.MkdirAll(path, 0777); err != nil {
		err = errors.New("Add2Log error os.MkdirAll: " + fmt.Sprint(err))
	}

	dateLog := time.Now().Format("2006_01_02")
	dateNow := time.Now().Format("2006.01.02")
	timeNow := time.Now().Format("15:04:05")

	mu.Lock()
	defer func() {
		mu.Unlock()
	}()

	if file, err := os.OpenFile(path+"/"+dateLog+"_log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); err == nil {
		if _, err = file.WriteString(dateNow + "	" + timeNow + "	" + who + " " + what + opt + "\n"); err == nil {
			if err = file.Close(); err != nil {
				err = errors.New("Add2Log error file.Close: " + fmt.Sprint(err))
			}
		} else {
			err = errors.New("Add2Log error file.WriteString: " + fmt.Sprint(err))
		}
	} else {
		err = errors.New("Add2Log error os.OpenFile: " + fmt.Sprint(err))
	}

	if _, err = os.Stderr.WriteString(dateNow + " " + timeNow + "	\033[91m" + who + "\033[0m " + whats + opt + "\n"); err != nil {
		err = errors.New("Add2Log error os.Stderr.WriteString: " + fmt.Sprint(err))
	}

	if err != nil {
		log.Print(err)
	}
}

func convert(b []byte) string {
	before := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{before.Data, before.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
