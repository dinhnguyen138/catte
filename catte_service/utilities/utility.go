package utilities

import "github.com/sparrc/go-ping"

func CheckPing(ipAddress string) bool {
	pinger, err := ping.NewPinger(ipAddress)
	if err != nil {
		panic(err)
	}

	pinger.Count = 3
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	if stats.PacketsSent-stats.PacketsRecv == 0 {
		return true
	} else {
		return false
	}
}
