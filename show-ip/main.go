package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// func convertToIP(ip uint32) string {
// 	b := 4
// 	var octet [4]uint32
// 	for i := 0; i < b; i++ {
// 		octet[i] = ip >> (i * 8)
// 	}
// 	return fmt.Sprintf("%d.%d.%d.%d", octet[3], octet[2], octet[1], octet[0])
// }

func PrintInfo(header []byte) {
	sip := header[0:4]
	dip := header[5:9]
	// proto := header[12:14]
	source := fmt.Sprintf("%d.%d.%d.%d", sip[0], sip[1], sip[2], sip[3])
	destination := fmt.Sprintf("%d.%d.%d.%d", dip[0], dip[1], dip[2], dip[3])
	fmt.Println("Source= ", source, " Destination= ", destination)
}

func main() {
	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	var objs ipObjects
	if err := loadIpObjects(&objs, nil); err != nil {
		log.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	ifname := "wlp2s0" // Change this to an interface on your machine.
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatalf("Getting interface %s: %s", ifname, err)
	}

	link, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.GetIps,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("Attaching XDP:", err)
	}
	defer link.Close()

	log.Printf("Reading IPs through %s..", ifname)

	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			// var ip uint32
			var val []byte
			err := objs.Ips.Lookup(uint32(0), &val)
			if err != nil {
				log.Fatal("Map lookup: ", err)
			}
			PrintInfo(val)
		case <-stop:
			log.Print("Received signal, exiting..")
			return
		}
	}
}
