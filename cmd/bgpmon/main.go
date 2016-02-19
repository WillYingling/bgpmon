package main

import (
	"fmt"
	"os"

	pb "github.com/hamersaw/bgpmon/proto/bgpmond"

	cli "github.com/jawher/mow.cli"
	"google.golang.org/grpc"
)

/*
bgpmon open
	cassandra
	file

bgpmon close

bgpmon list
	sessions
	modules

bgpmon start

bgpmon status
	bgpmond
	module

bgpmon stop

bgpmon write
	as-location
	mrt-file
	ip-location
	prefix-location
*/

var ipAddress, port string

func main() {
	bgpmon := cli.App("bgpmon", "seamless command interfaces to bgpmond")
	bgpmon.Version("v version", "bgpmon 0.0.1")

	ipAddress = *bgpmon.StringOpt("i ipAddress", "", "ip address of bgpmond host")
	port = *bgpmon.StringOpt("p port", "12289", "port of bgpmond host")

	bgpmon.Command("close", "close a session on bgpmond host", Close)

	bgpmon.Command("list", "list running/open components of bgpmond host", func(cmd *cli.Cmd) {
		cmd.Command("sessions", "list open sessions", ListSessions)
		cmd.Command("modules", "list running modules", ListModules)
	})

	bgpmon.Command("open", "open a session on bgpmond host", func(cmd *cli.Cmd) {
		cmd.Command("cassandra", "open a session over a cassandra cluster", OpenCassandra)
		cmd.Command("file", "open a session to a file", OpenFile)
	})

	bgpmon.Command("start", "start a module on bgpmond host", func(cmd *cli.Cmd) {
		cmd.Command("gobgpd-link", "start a gobgpd link module", StartGoBGPDLinkModule)
		cmd.Command("prefix-hijack", "start prefix hijack module", StartPrefixHijackModule)
	})

	bgpmon.Command("status", "check status of components of bgpmond host", func(cmd *cli.Cmd) {
		cmd.Command("module", "check status of module", StatusModule)
		cmd.Command("session", "check status of session", StatusSession)
	})

	bgpmon.Command("stop", "stop a module on bgpmond host", StopModule)

	bgpmon.Command("write", "write data to bgpmond host", func(cmd *cli.Cmd) {
		cmd.Command("mrt-file", "write bgp messages from mrt file", WriteMRTFile)
	})

	bgpmon.Run(os.Args)
}

func getRPCClient() (pb.BgpmondClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", ipAddress, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewBgpmondClient(conn), nil
}