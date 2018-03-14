/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/
package logging

import (
	"log"
	
	"os"
	"fmt"
	
)


func Plog(msg string) {
	
	s:="[INFO][PLUGINLOG] "+msg
	log.Printf(s)
	
	
}

func Plogf(msg string, a ...interface{}) {
	
	s:=fmt.Sprintf("[INFO][PLUGINLOG] "+msg, a)
	log.Printf(s)
	
	
}
func PlogWarn(msg string) {
	
	
	log.Printf("[WARN][PLUGINLOG]  "+msg)
	
	
}


func PlogErrorf(msg string, e error) {
	
	s:=fmt.Sprintf("[ERROR][PLUGINLOG] "+msg, e.Error())
	log.Printf(s)
	
	
}


//need not be called for terraform inited operation / only called from acceptance test
// still under validation
func Init(){
	log.SetOutput(os.Stdout)
	Plog("__INIT__LOGGING__")

}
