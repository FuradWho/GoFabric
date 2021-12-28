// Code generaTed by fileb0x at "2019-04-14 04:38:21.489607614 +0600 +06 m=+4.245926076" from config file "b0x.yml" DO NOT EDIT.

package swaggerFiles

import (
	"log"
	"os"
)

// FileSwaggerUICSSMap is "/swagger-ui.css.map"
var FileSwaggerUICSSMap = []byte("\x7b\x22\x76\x65\x72\x73\x69\x6f\x6e\x22\x3a\x33\x2c\x22\x73\x6f\x75\x72\x63\x65\x73\x22\x3a\x5b\x5d\x2c\x22\x6e\x61\x6d\x65\x73\x22\x3a\x5b\x5d\x2c\x22\x6d\x61\x70\x70\x69\x6e\x67\x73\x22\x3a\x22\x22\x2c\x22\x66\x69\x6c\x65\x22\x3a\x22\x73\x77\x61\x67\x67\x65\x72\x2d\x75\x69\x2e\x63\x73\x73\x22\x2c\x22\x73\x6f\x75\x72\x63\x65\x52\x6f\x6f\x74\x22\x3a\x22\x22\x7d")

func init() {

	f, err := FS.OpenFile(CTX, "/swagger-ui.css.map", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(FileSwaggerUICSSMap)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}