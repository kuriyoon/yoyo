package yoyoSystem

import (
	"net"
)

type NetInfo struct {
	InfName string
	InfIp string
}

func GetNetworkInfo() (rNetInfo []NetInfo, rErr error)  {

	_iface , rErr := net.Interfaces()
	if rErr != nil {
		return
	}

	for _, iface := range _iface {

		iP := ""
		//Iface_Ip := ""
		addrs, rErr := iface.Addrs()
		if rErr != nil {
			continue
		}

		for _, addr := range addrs {
			iP = iP + " " + addr.String()
		}

		var tmpNetInfo NetInfo
		tmpNetInfo.InfName = iface.Name
		tmpNetInfo.InfIp = iP

		if tmpNetInfo.InfIp != "" {
			rNetInfo = append(rNetInfo, tmpNetInfo)
		}
	}
	return
}
