package yoyoSystem

import (
	"os"
)

func init(){
}

// HostName Return
func GetHostname() (rHostname string, rErr error)  {
	rHostname, rErr = os.Hostname()
	return
}

func GetPid() (rPid int){
	rPid = os.Getpid()
	return
}