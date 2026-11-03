//go:build llgo && linux && arm64 && !android

package net

import (
	"github.com/goplus/lib/c"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

type nativeNetworkChar = uint8

func TestGeneratedNetworkLayouts(t *testing.T) {
	// Expected values are compiled from the original Bootlin glibc 2.24 C headers.
	var networkFunctionSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{
		{"InetAddr", c.Func(InetAddr)},
		{"InAddr.InetLnaof", c.Func((InAddr).InetLnaof)},
		{"InetMakeaddr", c.Func(InetMakeaddr)},
		{"InAddr.InetNetof", c.Func((InAddr).InetNetof)},
		{"InetNetwork", c.Func(InetNetwork)},
		{"InAddr.InetNtoa", c.Func((InAddr).InetNtoa)},
		{"InetPton", c.Func(InetPton)},
		{"InetNtop", c.Func(InetNtop)},
		{"InetAton", c.Func(InetAton)},
		{"InetNeta", c.Func(InetNeta)},
		{"InetNetNtop", c.Func(InetNetNtop)},
		{"InetNetPton", c.Func(InetNetPton)},
		{"InetNsapAddr", c.Func(InetNsapAddr)},
		{"InetNsapNtoa", c.Func(InetNsapNtoa)},
		{"IfNametoindex", c.Func(IfNametoindex)},
		{"IfIndextoname", c.Func(IfIndextoname)},
		{"IfNameindex", c.Func(IfNameindex)},
		{"*IfNameindexT.IfFreenameindex", c.Func((*IfNameindexT).IfFreenameindex)},
		{"Herror", c.Func(Herror)},
		{"Hstrerror", c.Func(Hstrerror)},
		{"Sethostent", c.Func(Sethostent)},
		{"Endhostent", c.Func(Endhostent)},
		{"Gethostent", c.Func(Gethostent)},
		{"Gethostbyaddr", c.Func(Gethostbyaddr)},
		{"GetHostByName", c.Func(GetHostByName)},
		{"Gethostbyname2", c.Func(Gethostbyname2)},
		{"*Hostent.GethostentR", c.Func((*Hostent).GethostentR)},
		{"GethostbyaddrR", c.Func(GethostbyaddrR)},
		{"GethostbynameR", c.Func(GethostbynameR)},
		{"Gethostbyname2R", c.Func(Gethostbyname2R)},
		{"Setnetent", c.Func(Setnetent)},
		{"Endnetent", c.Func(Endnetent)},
		{"Getnetent", c.Func(Getnetent)},
		{"Getnetbyaddr", c.Func(Getnetbyaddr)},
		{"Getnetbyname", c.Func(Getnetbyname)},
		{"*Netent.GetnetentR", c.Func((*Netent).GetnetentR)},
		{"GetnetbyaddrR", c.Func(GetnetbyaddrR)},
		{"GetnetbynameR", c.Func(GetnetbynameR)},
		{"Setservent", c.Func(Setservent)},
		{"Endservent", c.Func(Endservent)},
		{"Getservent", c.Func(Getservent)},
		{"Getservbyname", c.Func(Getservbyname)},
		{"Getservbyport", c.Func(Getservbyport)},
		{"*Servent.GetserventR", c.Func((*Servent).GetserventR)},
		{"GetservbynameR", c.Func(GetservbynameR)},
		{"GetservbyportR", c.Func(GetservbyportR)},
		{"Setprotoent", c.Func(Setprotoent)},
		{"Endprotoent", c.Func(Endprotoent)},
		{"Getprotoent", c.Func(Getprotoent)},
		{"Getprotobyname", c.Func(Getprotobyname)},
		{"Getprotobynumber", c.Func(Getprotobynumber)},
		{"*Protoent.GetprotoentR", c.Func((*Protoent).GetprotoentR)},
		{"GetprotobynameR", c.Func(GetprotobynameR)},
		{"GetprotobynumberR", c.Func(GetprotobynumberR)},
		{"Setnetgrent", c.Func(Setnetgrent)},
		{"Endnetgrent", c.Func(Endnetgrent)},
		{"Getnetgrent", c.Func(Getnetgrent)},
		{"Innetgr", c.Func(Innetgr)},
		{"GetnetgrentR", c.Func(GetnetgrentR)},
		{"Rcmd", c.Func(Rcmd)},
		{"RcmdAf", c.Func(RcmdAf)},
		{"Rexec", c.Func(Rexec)},
		{"RexecAf", c.Func(RexecAf)},
		{"Ruserok", c.Func(Ruserok)},
		{"RuserokAf", c.Func(RuserokAf)},
		{"Iruserok", c.Func(Iruserok)},
		{"IruserokAf", c.Func(IruserokAf)},
		{"Rresvport", c.Func(Rresvport)},
		{"RresvportAf", c.Func(RresvportAf)},
		{"Getaddrinfo", c.Func(Getaddrinfo)},
		{"Freeaddrinfo", c.Func(Freeaddrinfo)},
		{"GaiStrerror", c.Func(GaiStrerror)},
		{"*SockAddr.Getnameinfo", c.Func((*SockAddr).Getnameinfo)},
		{"Ntohl", c.Func(Ntohl)},
		{"Ntohs", c.Func(Ntohs)},
		{"Htonl", c.Func(Htonl)},
		{"Htons", c.Func(Htons)},
		{"Bindresvport", c.Func(Bindresvport)},
		{"Setrpcent", c.Func(Setrpcent)},
		{"Endrpcent", c.Func(Endrpcent)},
		{"Getrpcbyname", c.Func(Getrpcbyname)},
		{"Getrpcbynumber", c.Func(Getrpcbynumber)},
		{"Getrpcent", c.Func(Getrpcent)},
		{"GetrpcbynameR", c.Func(GetrpcbynameR)},
		{"GetrpcbynumberR", c.Func(GetrpcbynumberR)},
		{"*Rpcent.GetrpcentR", c.Func((*Rpcent).GetrpcentR)},
		{"Socket", c.Func(Socket)},
		{"Socketpair", c.Func(Socketpair)},
		{"Bind", c.Func(Bind)},
		{"Getsockname", c.Func(Getsockname)},
		{"Connect", c.Func(Connect)},
		{"Getpeername", c.Func(Getpeername)},
		{"Send", c.Func(Send)},
		{"Recv", c.Func(Recv)},
		{"Sendto", c.Func(Sendto)},
		{"Recvfrom", c.Func(Recvfrom)},
		{"Sendmsg", c.Func(Sendmsg)},
		{"Recvmsg", c.Func(Recvmsg)},
		{"Getsockopt", c.Func(Getsockopt)},
		{"SetSockOpt", c.Func(SetSockOpt)},
		{"Listen", c.Func(Listen)},
		{"Accept", c.Func(Accept)},
		{"Shutdown", c.Func(Shutdown)},
		{"Sockatmark", c.Func(Sockatmark)},
		{"Isfdtype", c.Func(Isfdtype)},
	}
	var networkFunctionPrototypes = []struct {
		name      string
		signature reflect.Type
		inputs    []networkABISlot
		output    *networkABISlot
		variadic  bool
	}{
		{"InetAddr", reflect.TypeOf(InetAddr), []networkABISlot{{8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"InAddr.InetLnaof", reflect.TypeOf((InAddr).InetLnaof), []networkABISlot{{4, 4, 12}}, &networkABISlot{4, 4, 1}, false},
		{"InetMakeaddr", reflect.TypeOf(InetMakeaddr), []networkABISlot{{4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 12}, false},
		{"InAddr.InetNetof", reflect.TypeOf((InAddr).InetNetof), []networkABISlot{{4, 4, 12}}, &networkABISlot{4, 4, 1}, false},
		{"InetNetwork", reflect.TypeOf(InetNetwork), []networkABISlot{{8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"InAddr.InetNtoa", reflect.TypeOf((InAddr).InetNtoa), []networkABISlot{{4, 4, 12}}, &networkABISlot{8, 8, 5}, false},
		{"InetPton", reflect.TypeOf(InetPton), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"InetNtop", reflect.TypeOf(InetNtop), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"InetAton", reflect.TypeOf(InetAton), []networkABISlot{{8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"InetNeta", reflect.TypeOf(InetNeta), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 1}}, &networkABISlot{8, 8, 5}, false},
		{"InetNetNtop", reflect.TypeOf(InetNetNtop), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 1}}, &networkABISlot{8, 8, 5}, false},
		{"InetNetPton", reflect.TypeOf(InetNetPton), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}}, &networkABISlot{4, 4, 1}, false},
		{"InetNsapAddr", reflect.TypeOf(InetNsapAddr), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"InetNsapNtoa", reflect.TypeOf(InetNsapNtoa), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"IfNametoindex", reflect.TypeOf(IfNametoindex), []networkABISlot{{8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"IfIndextoname", reflect.TypeOf(IfIndextoname), []networkABISlot{{4, 4, 1}, {8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"IfNameindex", reflect.TypeOf(IfNameindex), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"*IfNameindexT.IfFreenameindex", reflect.TypeOf((*IfNameindexT).IfFreenameindex), []networkABISlot{{8, 8, 5}}, nil, false},
		{"Herror", reflect.TypeOf(Herror), []networkABISlot{{8, 8, 5}}, nil, false},
		{"Hstrerror", reflect.TypeOf(Hstrerror), []networkABISlot{{4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"Sethostent", reflect.TypeOf(Sethostent), []networkABISlot{{4, 4, 1}}, nil, false},
		{"Endhostent", reflect.TypeOf(Endhostent), []networkABISlot{}, nil, false},
		{"Gethostent", reflect.TypeOf(Gethostent), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"Gethostbyaddr", reflect.TypeOf(Gethostbyaddr), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"GetHostByName", reflect.TypeOf(GetHostByName), []networkABISlot{{8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"Gethostbyname2", reflect.TypeOf(Gethostbyname2), []networkABISlot{{8, 8, 5}, {4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"*Hostent.GethostentR", reflect.TypeOf((*Hostent).GethostentR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GethostbyaddrR", reflect.TypeOf(GethostbyaddrR), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GethostbynameR", reflect.TypeOf(GethostbynameR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Gethostbyname2R", reflect.TypeOf(Gethostbyname2R), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Setnetent", reflect.TypeOf(Setnetent), []networkABISlot{{4, 4, 1}}, nil, false},
		{"Endnetent", reflect.TypeOf(Endnetent), []networkABISlot{}, nil, false},
		{"Getnetent", reflect.TypeOf(Getnetent), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"Getnetbyaddr", reflect.TypeOf(Getnetbyaddr), []networkABISlot{{4, 4, 1}, {4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"Getnetbyname", reflect.TypeOf(Getnetbyname), []networkABISlot{{8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"*Netent.GetnetentR", reflect.TypeOf((*Netent).GetnetentR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetnetbyaddrR", reflect.TypeOf(GetnetbyaddrR), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetnetbynameR", reflect.TypeOf(GetnetbynameR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Setservent", reflect.TypeOf(Setservent), []networkABISlot{{4, 4, 1}}, nil, false},
		{"Endservent", reflect.TypeOf(Endservent), []networkABISlot{}, nil, false},
		{"Getservent", reflect.TypeOf(Getservent), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"Getservbyname", reflect.TypeOf(Getservbyname), []networkABISlot{{8, 8, 5}, {8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"Getservbyport", reflect.TypeOf(Getservbyport), []networkABISlot{{4, 4, 1}, {8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"*Servent.GetserventR", reflect.TypeOf((*Servent).GetserventR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetservbynameR", reflect.TypeOf(GetservbynameR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetservbyportR", reflect.TypeOf(GetservbyportR), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Setprotoent", reflect.TypeOf(Setprotoent), []networkABISlot{{4, 4, 1}}, nil, false},
		{"Endprotoent", reflect.TypeOf(Endprotoent), []networkABISlot{}, nil, false},
		{"Getprotoent", reflect.TypeOf(Getprotoent), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"Getprotobyname", reflect.TypeOf(Getprotobyname), []networkABISlot{{8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"Getprotobynumber", reflect.TypeOf(Getprotobynumber), []networkABISlot{{4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"*Protoent.GetprotoentR", reflect.TypeOf((*Protoent).GetprotoentR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetprotobynameR", reflect.TypeOf(GetprotobynameR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetprotobynumberR", reflect.TypeOf(GetprotobynumberR), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Setnetgrent", reflect.TypeOf(Setnetgrent), []networkABISlot{{8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Endnetgrent", reflect.TypeOf(Endnetgrent), []networkABISlot{}, nil, false},
		{"Getnetgrent", reflect.TypeOf(Getnetgrent), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Innetgr", reflect.TypeOf(Innetgr), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetnetgrentR", reflect.TypeOf(GetnetgrentR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Rcmd", reflect.TypeOf(Rcmd), []networkABISlot{{8, 8, 5}, {2, 2, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"RcmdAf", reflect.TypeOf(RcmdAf), []networkABISlot{{8, 8, 5}, {2, 2, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {2, 2, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Rexec", reflect.TypeOf(Rexec), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"RexecAf", reflect.TypeOf(RexecAf), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {2, 2, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Ruserok", reflect.TypeOf(Ruserok), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"RuserokAf", reflect.TypeOf(RuserokAf), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {2, 2, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Iruserok", reflect.TypeOf(Iruserok), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"IruserokAf", reflect.TypeOf(IruserokAf), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {2, 2, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Rresvport", reflect.TypeOf(Rresvport), []networkABISlot{{8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"RresvportAf", reflect.TypeOf(RresvportAf), []networkABISlot{{8, 8, 5}, {2, 2, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Getaddrinfo", reflect.TypeOf(Getaddrinfo), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Freeaddrinfo", reflect.TypeOf(Freeaddrinfo), []networkABISlot{{8, 8, 5}}, nil, false},
		{"GaiStrerror", reflect.TypeOf(GaiStrerror), []networkABISlot{{4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"*SockAddr.Getnameinfo", reflect.TypeOf((*SockAddr).Getnameinfo), []networkABISlot{{8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {4, 4, 1}, {8, 8, 5}, {4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Ntohl", reflect.TypeOf(Ntohl), []networkABISlot{{4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Ntohs", reflect.TypeOf(Ntohs), []networkABISlot{{2, 2, 1}}, &networkABISlot{2, 2, 1}, false},
		{"Htonl", reflect.TypeOf(Htonl), []networkABISlot{{4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Htons", reflect.TypeOf(Htons), []networkABISlot{{2, 2, 1}}, &networkABISlot{2, 2, 1}, false},
		{"Bindresvport", reflect.TypeOf(Bindresvport), []networkABISlot{{4, 4, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Setrpcent", reflect.TypeOf(Setrpcent), []networkABISlot{{4, 4, 1}}, nil, false},
		{"Endrpcent", reflect.TypeOf(Endrpcent), []networkABISlot{}, nil, false},
		{"Getrpcbyname", reflect.TypeOf(Getrpcbyname), []networkABISlot{{8, 8, 5}}, &networkABISlot{8, 8, 5}, false},
		{"Getrpcbynumber", reflect.TypeOf(Getrpcbynumber), []networkABISlot{{4, 4, 1}}, &networkABISlot{8, 8, 5}, false},
		{"Getrpcent", reflect.TypeOf(Getrpcent), []networkABISlot{}, &networkABISlot{8, 8, 5}, false},
		{"GetrpcbynameR", reflect.TypeOf(GetrpcbynameR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"GetrpcbynumberR", reflect.TypeOf(GetrpcbynumberR), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"*Rpcent.GetrpcentR", reflect.TypeOf((*Rpcent).GetrpcentR), []networkABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Socket", reflect.TypeOf(Socket), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Socketpair", reflect.TypeOf(Socketpair), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {4, 4, 1}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Bind", reflect.TypeOf(Bind), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Getsockname", reflect.TypeOf(Getsockname), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Connect", reflect.TypeOf(Connect), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Getpeername", reflect.TypeOf(Getpeername), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Send", reflect.TypeOf(Send), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 1}, {4, 4, 1}}, &networkABISlot{8, 8, 1}, false},
		{"Recv", reflect.TypeOf(Recv), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 1}, {4, 4, 1}}, &networkABISlot{8, 8, 1}, false},
		{"Sendto", reflect.TypeOf(Sendto), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 1}, {4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{8, 8, 1}, false},
		{"Recvfrom", reflect.TypeOf(Recvfrom), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 1}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{8, 8, 1}, false},
		{"Sendmsg", reflect.TypeOf(Sendmsg), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{8, 8, 1}, false},
		{"Recvmsg", reflect.TypeOf(Recvmsg), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{8, 8, 1}, false},
		{"Getsockopt", reflect.TypeOf(Getsockopt), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"SetSockOpt", reflect.TypeOf(SetSockOpt), []networkABISlot{{4, 4, 1}, {4, 4, 1}, {4, 4, 1}, {8, 8, 5}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Listen", reflect.TypeOf(Listen), []networkABISlot{{4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Accept", reflect.TypeOf(Accept), []networkABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &networkABISlot{4, 4, 1}, false},
		{"Shutdown", reflect.TypeOf(Shutdown), []networkABISlot{{4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Sockatmark", reflect.TypeOf(Sockatmark), []networkABISlot{{4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
		{"Isfdtype", reflect.TypeOf(Isfdtype), []networkABISlot{{4, 4, 1}, {4, 4, 1}}, &networkABISlot{4, 4, 1}, false},
	}
	var networkTypeLayouts = []struct {
		symbol, property string
		got, want        uintptr
	}{
		{"IpOpts", "size", unsafe.Sizeof(*new(IpOpts)), 44},
		{"IpOpts", "alignment", unsafe.Alignof(*new(IpOpts)), 4},
		{"IpOpts", "IpDst offset", unsafe.Offsetof((*new(IpOpts)).IpDst), 0},
		{"IpOpts", "IpDst size", unsafe.Sizeof((*new(IpOpts)).IpDst), 4},
		{"IpOpts", "IpDst alignment", unsafe.Alignof((*new(IpOpts)).IpDst), 4},
		{"IpOpts", "IpOpts offset", unsafe.Offsetof((*new(IpOpts)).IpOpts), 4},
		{"IpOpts", "IpOpts size", unsafe.Sizeof((*new(IpOpts)).IpOpts), 40},
		{"IpOpts", "IpOpts alignment", unsafe.Alignof((*new(IpOpts)).IpOpts), 1},
		{"IpMreqn", "size", unsafe.Sizeof(*new(IpMreqn)), 12},
		{"IpMreqn", "alignment", unsafe.Alignof(*new(IpMreqn)), 4},
		{"IpMreqn", "ImrMultiaddr offset", unsafe.Offsetof((*new(IpMreqn)).ImrMultiaddr), 0},
		{"IpMreqn", "ImrMultiaddr size", unsafe.Sizeof((*new(IpMreqn)).ImrMultiaddr), 4},
		{"IpMreqn", "ImrMultiaddr alignment", unsafe.Alignof((*new(IpMreqn)).ImrMultiaddr), 4},
		{"IpMreqn", "ImrAddress offset", unsafe.Offsetof((*new(IpMreqn)).ImrAddress), 4},
		{"IpMreqn", "ImrAddress size", unsafe.Sizeof((*new(IpMreqn)).ImrAddress), 4},
		{"IpMreqn", "ImrAddress alignment", unsafe.Alignof((*new(IpMreqn)).ImrAddress), 4},
		{"IpMreqn", "ImrIfindex offset", unsafe.Offsetof((*new(IpMreqn)).ImrIfindex), 8},
		{"IpMreqn", "ImrIfindex size", unsafe.Sizeof((*new(IpMreqn)).ImrIfindex), 4},
		{"IpMreqn", "ImrIfindex alignment", unsafe.Alignof((*new(IpMreqn)).ImrIfindex), 4},
		{"InPktinfo", "size", unsafe.Sizeof(*new(InPktinfo)), 12},
		{"InPktinfo", "alignment", unsafe.Alignof(*new(InPktinfo)), 4},
		{"InPktinfo", "IpiIfindex offset", unsafe.Offsetof((*new(InPktinfo)).IpiIfindex), 0},
		{"InPktinfo", "IpiIfindex size", unsafe.Sizeof((*new(InPktinfo)).IpiIfindex), 4},
		{"InPktinfo", "IpiIfindex alignment", unsafe.Alignof((*new(InPktinfo)).IpiIfindex), 4},
		{"InPktinfo", "IpiSpecDst offset", unsafe.Offsetof((*new(InPktinfo)).IpiSpecDst), 4},
		{"InPktinfo", "IpiSpecDst size", unsafe.Sizeof((*new(InPktinfo)).IpiSpecDst), 4},
		{"InPktinfo", "IpiSpecDst alignment", unsafe.Alignof((*new(InPktinfo)).IpiSpecDst), 4},
		{"InPktinfo", "IpiAddr offset", unsafe.Offsetof((*new(InPktinfo)).IpiAddr), 8},
		{"InPktinfo", "IpiAddr size", unsafe.Sizeof((*new(InPktinfo)).IpiAddr), 4},
		{"InPktinfo", "IpiAddr alignment", unsafe.Alignof((*new(InPktinfo)).IpiAddr), 4},
		{"Netent", "size", unsafe.Sizeof(*new(Netent)), 24},
		{"Netent", "alignment", unsafe.Alignof(*new(Netent)), 8},
		{"Netent", "NName offset", unsafe.Offsetof((*new(Netent)).NName), 0},
		{"Netent", "NName size", unsafe.Sizeof((*new(Netent)).NName), 8},
		{"Netent", "NName alignment", unsafe.Alignof((*new(Netent)).NName), 8},
		{"Netent", "NAliases offset", unsafe.Offsetof((*new(Netent)).NAliases), 8},
		{"Netent", "NAliases size", unsafe.Sizeof((*new(Netent)).NAliases), 8},
		{"Netent", "NAliases alignment", unsafe.Alignof((*new(Netent)).NAliases), 8},
		{"Netent", "NAddrtype offset", unsafe.Offsetof((*new(Netent)).NAddrtype), 16},
		{"Netent", "NAddrtype size", unsafe.Sizeof((*new(Netent)).NAddrtype), 4},
		{"Netent", "NAddrtype alignment", unsafe.Alignof((*new(Netent)).NAddrtype), 4},
		{"Netent", "NNet offset", unsafe.Offsetof((*new(Netent)).NNet), 20},
		{"Netent", "NNet size", unsafe.Sizeof((*new(Netent)).NNet), 4},
		{"Netent", "NNet alignment", unsafe.Alignof((*new(Netent)).NNet), 4},
		{"SaFamilyT", "size", unsafe.Sizeof(*new(SaFamilyT)), 2},
		{"SaFamilyT", "alignment", unsafe.Alignof(*new(SaFamilyT)), 2},
		{"SockAddr", "size", unsafe.Sizeof(*new(SockAddr)), 16},
		{"SockAddr", "alignment", unsafe.Alignof(*new(SockAddr)), 2},
		{"SockAddr", "Family offset", unsafe.Offsetof((*new(SockAddr)).Family), 0},
		{"SockAddr", "Family size", unsafe.Sizeof((*new(SockAddr)).Family), 2},
		{"SockAddr", "Family alignment", unsafe.Alignof((*new(SockAddr)).Family), 2},
		{"SockAddr", "Data offset", unsafe.Offsetof((*new(SockAddr)).Data), 2},
		{"SockAddr", "Data size", unsafe.Sizeof((*new(SockAddr)).Data), 14},
		{"SockAddr", "Data alignment", unsafe.Alignof((*new(SockAddr)).Data), 1},
		{"SockaddrStorage", "size", unsafe.Sizeof(*new(SockaddrStorage)), 128},
		{"SockaddrStorage", "alignment", unsafe.Alignof(*new(SockaddrStorage)), 8},
		{"SockaddrStorage", "Family offset", unsafe.Offsetof((*new(SockaddrStorage)).Family), 0},
		{"SockaddrStorage", "Family size", unsafe.Sizeof((*new(SockaddrStorage)).Family), 2},
		{"SockaddrStorage", "Family alignment", unsafe.Alignof((*new(SockaddrStorage)).Family), 2},
		{"SockaddrStorage", "pad offset", unsafe.Offsetof((*new(SockaddrStorage)).pad), 2},
		{"SockaddrStorage", "pad size", unsafe.Sizeof((*new(SockaddrStorage)).pad), 118},
		{"SockaddrStorage", "pad alignment", unsafe.Alignof((*new(SockaddrStorage)).pad), 1},
		{"SockaddrStorage", "align offset", unsafe.Offsetof((*new(SockaddrStorage)).align), 120},
		{"SockaddrStorage", "align size", unsafe.Sizeof((*new(SockaddrStorage)).align), 8},
		{"SockaddrStorage", "align alignment", unsafe.Alignof((*new(SockaddrStorage)).align), 8},
		{"Msghdr", "size", unsafe.Sizeof(*new(Msghdr)), 56},
		{"Msghdr", "alignment", unsafe.Alignof(*new(Msghdr)), 8},
		{"Msghdr", "MsgName offset", unsafe.Offsetof((*new(Msghdr)).MsgName), 0},
		{"Msghdr", "MsgName size", unsafe.Sizeof((*new(Msghdr)).MsgName), 8},
		{"Msghdr", "MsgName alignment", unsafe.Alignof((*new(Msghdr)).MsgName), 8},
		{"Msghdr", "MsgNamelen offset", unsafe.Offsetof((*new(Msghdr)).MsgNamelen), 8},
		{"Msghdr", "MsgNamelen size", unsafe.Sizeof((*new(Msghdr)).MsgNamelen), 4},
		{"Msghdr", "MsgNamelen alignment", unsafe.Alignof((*new(Msghdr)).MsgNamelen), 4},
		{"Msghdr", "MsgIov offset", unsafe.Offsetof((*new(Msghdr)).MsgIov), 16},
		{"Msghdr", "MsgIov size", unsafe.Sizeof((*new(Msghdr)).MsgIov), 8},
		{"Msghdr", "MsgIov alignment", unsafe.Alignof((*new(Msghdr)).MsgIov), 8},
		{"Msghdr", "MsgIovlen offset", unsafe.Offsetof((*new(Msghdr)).MsgIovlen), 24},
		{"Msghdr", "MsgIovlen size", unsafe.Sizeof((*new(Msghdr)).MsgIovlen), 8},
		{"Msghdr", "MsgIovlen alignment", unsafe.Alignof((*new(Msghdr)).MsgIovlen), 8},
		{"Msghdr", "MsgControl offset", unsafe.Offsetof((*new(Msghdr)).MsgControl), 32},
		{"Msghdr", "MsgControl size", unsafe.Sizeof((*new(Msghdr)).MsgControl), 8},
		{"Msghdr", "MsgControl alignment", unsafe.Alignof((*new(Msghdr)).MsgControl), 8},
		{"Msghdr", "MsgControllen offset", unsafe.Offsetof((*new(Msghdr)).MsgControllen), 40},
		{"Msghdr", "MsgControllen size", unsafe.Sizeof((*new(Msghdr)).MsgControllen), 8},
		{"Msghdr", "MsgControllen alignment", unsafe.Alignof((*new(Msghdr)).MsgControllen), 8},
		{"Msghdr", "MsgFlags offset", unsafe.Offsetof((*new(Msghdr)).MsgFlags), 48},
		{"Msghdr", "MsgFlags size", unsafe.Sizeof((*new(Msghdr)).MsgFlags), 4},
		{"Msghdr", "MsgFlags alignment", unsafe.Alignof((*new(Msghdr)).MsgFlags), 4},
		{"Cmsghdr", "size", unsafe.Sizeof(*new(Cmsghdr)), 16},
		{"Cmsghdr", "alignment", unsafe.Alignof(*new(Cmsghdr)), 8},
		{"Cmsghdr", "CmsgLen offset", unsafe.Offsetof((*new(Cmsghdr)).CmsgLen), 0},
		{"Cmsghdr", "CmsgLen size", unsafe.Sizeof((*new(Cmsghdr)).CmsgLen), 8},
		{"Cmsghdr", "CmsgLen alignment", unsafe.Alignof((*new(Cmsghdr)).CmsgLen), 8},
		{"Cmsghdr", "CmsgLevel offset", unsafe.Offsetof((*new(Cmsghdr)).CmsgLevel), 8},
		{"Cmsghdr", "CmsgLevel size", unsafe.Sizeof((*new(Cmsghdr)).CmsgLevel), 4},
		{"Cmsghdr", "CmsgLevel alignment", unsafe.Alignof((*new(Cmsghdr)).CmsgLevel), 4},
		{"Cmsghdr", "CmsgType offset", unsafe.Offsetof((*new(Cmsghdr)).CmsgType), 12},
		{"Cmsghdr", "CmsgType size", unsafe.Sizeof((*new(Cmsghdr)).CmsgType), 4},
		{"Cmsghdr", "CmsgType alignment", unsafe.Alignof((*new(Cmsghdr)).CmsgType), 4},
		{"Linger", "size", unsafe.Sizeof(*new(Linger)), 8},
		{"Linger", "alignment", unsafe.Alignof(*new(Linger)), 4},
		{"Linger", "LOnoff offset", unsafe.Offsetof((*new(Linger)).LOnoff), 0},
		{"Linger", "LOnoff size", unsafe.Sizeof((*new(Linger)).LOnoff), 4},
		{"Linger", "LOnoff alignment", unsafe.Alignof((*new(Linger)).LOnoff), 4},
		{"Linger", "LLinger offset", unsafe.Offsetof((*new(Linger)).LLinger), 4},
		{"Linger", "LLinger size", unsafe.Sizeof((*new(Linger)).LLinger), 4},
		{"Linger", "LLinger alignment", unsafe.Alignof((*new(Linger)).LLinger), 4},
		{"SocklenT", "size", unsafe.Sizeof(*new(SocklenT)), 4},
		{"SocklenT", "alignment", unsafe.Alignof(*new(SocklenT)), 4},
		{"X__SocketType", "size", unsafe.Sizeof(*new(X__SocketType)), 4},
		{"X__SocketType", "alignment", unsafe.Alignof(*new(X__SocketType)), 4},
		{"X__FsidT", "size", unsafe.Sizeof(*new(X__FsidT)), 8},
		{"X__FsidT", "alignment", unsafe.Alignof(*new(X__FsidT)), 4},
		{"X__FsidT", "X__Val offset", unsafe.Offsetof((*new(X__FsidT)).X__Val), 0},
		{"X__FsidT", "X__Val size", unsafe.Sizeof((*new(X__FsidT)).X__Val), 8},
		{"X__FsidT", "X__Val alignment", unsafe.Alignof((*new(X__FsidT)).X__Val), 4},
		{"X__UChar", "size", unsafe.Sizeof(*new(X__UChar)), 1},
		{"X__UChar", "alignment", unsafe.Alignof(*new(X__UChar)), 1},
		{"X__UShort", "size", unsafe.Sizeof(*new(X__UShort)), 2},
		{"X__UShort", "alignment", unsafe.Alignof(*new(X__UShort)), 2},
		{"X__UInt", "size", unsafe.Sizeof(*new(X__UInt)), 4},
		{"X__UInt", "alignment", unsafe.Alignof(*new(X__UInt)), 4},
		{"X__ULong", "size", unsafe.Sizeof(*new(X__ULong)), 8},
		{"X__ULong", "alignment", unsafe.Alignof(*new(X__ULong)), 8},
		{"X__Int8T", "size", unsafe.Sizeof(*new(X__Int8T)), 1},
		{"X__Int8T", "alignment", unsafe.Alignof(*new(X__Int8T)), 1},
		{"X__Uint8T", "size", unsafe.Sizeof(*new(X__Uint8T)), 1},
		{"X__Uint8T", "alignment", unsafe.Alignof(*new(X__Uint8T)), 1},
		{"X__Int16T", "size", unsafe.Sizeof(*new(X__Int16T)), 2},
		{"X__Int16T", "alignment", unsafe.Alignof(*new(X__Int16T)), 2},
		{"X__Uint16T", "size", unsafe.Sizeof(*new(X__Uint16T)), 2},
		{"X__Uint16T", "alignment", unsafe.Alignof(*new(X__Uint16T)), 2},
		{"X__Int32T", "size", unsafe.Sizeof(*new(X__Int32T)), 4},
		{"X__Int32T", "alignment", unsafe.Alignof(*new(X__Int32T)), 4},
		{"X__Uint32T", "size", unsafe.Sizeof(*new(X__Uint32T)), 4},
		{"X__Uint32T", "alignment", unsafe.Alignof(*new(X__Uint32T)), 4},
		{"X__Int64T", "size", unsafe.Sizeof(*new(X__Int64T)), 8},
		{"X__Int64T", "alignment", unsafe.Alignof(*new(X__Int64T)), 8},
		{"X__Uint64T", "size", unsafe.Sizeof(*new(X__Uint64T)), 8},
		{"X__Uint64T", "alignment", unsafe.Alignof(*new(X__Uint64T)), 8},
		{"X__QuadT", "size", unsafe.Sizeof(*new(X__QuadT)), 8},
		{"X__QuadT", "alignment", unsafe.Alignof(*new(X__QuadT)), 8},
		{"X__UQuadT", "size", unsafe.Sizeof(*new(X__UQuadT)), 8},
		{"X__UQuadT", "alignment", unsafe.Alignof(*new(X__UQuadT)), 8},
		{"X__DevT", "size", unsafe.Sizeof(*new(X__DevT)), 8},
		{"X__DevT", "alignment", unsafe.Alignof(*new(X__DevT)), 8},
		{"X__UidT", "size", unsafe.Sizeof(*new(X__UidT)), 4},
		{"X__UidT", "alignment", unsafe.Alignof(*new(X__UidT)), 4},
		{"X__GidT", "size", unsafe.Sizeof(*new(X__GidT)), 4},
		{"X__GidT", "alignment", unsafe.Alignof(*new(X__GidT)), 4},
		{"X__InoT", "size", unsafe.Sizeof(*new(X__InoT)), 8},
		{"X__InoT", "alignment", unsafe.Alignof(*new(X__InoT)), 8},
		{"X__Ino64T", "size", unsafe.Sizeof(*new(X__Ino64T)), 8},
		{"X__Ino64T", "alignment", unsafe.Alignof(*new(X__Ino64T)), 8},
		{"X__ModeT", "size", unsafe.Sizeof(*new(X__ModeT)), 4},
		{"X__ModeT", "alignment", unsafe.Alignof(*new(X__ModeT)), 4},
		{"X__NlinkT", "size", unsafe.Sizeof(*new(X__NlinkT)), 4},
		{"X__NlinkT", "alignment", unsafe.Alignof(*new(X__NlinkT)), 4},
		{"X__OffT", "size", unsafe.Sizeof(*new(X__OffT)), 8},
		{"X__OffT", "alignment", unsafe.Alignof(*new(X__OffT)), 8},
		{"X__Off64T", "size", unsafe.Sizeof(*new(X__Off64T)), 8},
		{"X__Off64T", "alignment", unsafe.Alignof(*new(X__Off64T)), 8},
		{"X__PidT", "size", unsafe.Sizeof(*new(X__PidT)), 4},
		{"X__PidT", "alignment", unsafe.Alignof(*new(X__PidT)), 4},
		{"X__ClockT", "size", unsafe.Sizeof(*new(X__ClockT)), 8},
		{"X__ClockT", "alignment", unsafe.Alignof(*new(X__ClockT)), 8},
		{"X__RlimT", "size", unsafe.Sizeof(*new(X__RlimT)), 8},
		{"X__RlimT", "alignment", unsafe.Alignof(*new(X__RlimT)), 8},
		{"X__Rlim64T", "size", unsafe.Sizeof(*new(X__Rlim64T)), 8},
		{"X__Rlim64T", "alignment", unsafe.Alignof(*new(X__Rlim64T)), 8},
		{"X__IdT", "size", unsafe.Sizeof(*new(X__IdT)), 4},
		{"X__IdT", "alignment", unsafe.Alignof(*new(X__IdT)), 4},
		{"X__TimeT", "size", unsafe.Sizeof(*new(X__TimeT)), 8},
		{"X__TimeT", "alignment", unsafe.Alignof(*new(X__TimeT)), 8},
		{"X__UsecondsT", "size", unsafe.Sizeof(*new(X__UsecondsT)), 4},
		{"X__UsecondsT", "alignment", unsafe.Alignof(*new(X__UsecondsT)), 4},
		{"X__SusecondsT", "size", unsafe.Sizeof(*new(X__SusecondsT)), 8},
		{"X__SusecondsT", "alignment", unsafe.Alignof(*new(X__SusecondsT)), 8},
		{"X__DaddrT", "size", unsafe.Sizeof(*new(X__DaddrT)), 4},
		{"X__DaddrT", "alignment", unsafe.Alignof(*new(X__DaddrT)), 4},
		{"X__KeyT", "size", unsafe.Sizeof(*new(X__KeyT)), 4},
		{"X__KeyT", "alignment", unsafe.Alignof(*new(X__KeyT)), 4},
		{"X__ClockidT", "size", unsafe.Sizeof(*new(X__ClockidT)), 4},
		{"X__ClockidT", "alignment", unsafe.Alignof(*new(X__ClockidT)), 4},
		{"X__TimerT", "size", unsafe.Sizeof(*new(X__TimerT)), 8},
		{"X__TimerT", "alignment", unsafe.Alignof(*new(X__TimerT)), 8},
		{"X__BlksizeT", "size", unsafe.Sizeof(*new(X__BlksizeT)), 4},
		{"X__BlksizeT", "alignment", unsafe.Alignof(*new(X__BlksizeT)), 4},
		{"X__BlkcntT", "size", unsafe.Sizeof(*new(X__BlkcntT)), 8},
		{"X__BlkcntT", "alignment", unsafe.Alignof(*new(X__BlkcntT)), 8},
		{"X__Blkcnt64T", "size", unsafe.Sizeof(*new(X__Blkcnt64T)), 8},
		{"X__Blkcnt64T", "alignment", unsafe.Alignof(*new(X__Blkcnt64T)), 8},
		{"X__FsblkcntT", "size", unsafe.Sizeof(*new(X__FsblkcntT)), 8},
		{"X__FsblkcntT", "alignment", unsafe.Alignof(*new(X__FsblkcntT)), 8},
		{"X__Fsblkcnt64T", "size", unsafe.Sizeof(*new(X__Fsblkcnt64T)), 8},
		{"X__Fsblkcnt64T", "alignment", unsafe.Alignof(*new(X__Fsblkcnt64T)), 8},
		{"X__FsfilcntT", "size", unsafe.Sizeof(*new(X__FsfilcntT)), 8},
		{"X__FsfilcntT", "alignment", unsafe.Alignof(*new(X__FsfilcntT)), 8},
		{"X__Fsfilcnt64T", "size", unsafe.Sizeof(*new(X__Fsfilcnt64T)), 8},
		{"X__Fsfilcnt64T", "alignment", unsafe.Alignof(*new(X__Fsfilcnt64T)), 8},
		{"X__FswordT", "size", unsafe.Sizeof(*new(X__FswordT)), 8},
		{"X__FswordT", "alignment", unsafe.Alignof(*new(X__FswordT)), 8},
		{"X__SsizeT", "size", unsafe.Sizeof(*new(X__SsizeT)), 8},
		{"X__SsizeT", "alignment", unsafe.Alignof(*new(X__SsizeT)), 8},
		{"X__SyscallSlongT", "size", unsafe.Sizeof(*new(X__SyscallSlongT)), 8},
		{"X__SyscallSlongT", "alignment", unsafe.Alignof(*new(X__SyscallSlongT)), 8},
		{"X__SyscallUlongT", "size", unsafe.Sizeof(*new(X__SyscallUlongT)), 8},
		{"X__SyscallUlongT", "alignment", unsafe.Alignof(*new(X__SyscallUlongT)), 8},
		{"X__LoffT", "size", unsafe.Sizeof(*new(X__LoffT)), 8},
		{"X__LoffT", "alignment", unsafe.Alignof(*new(X__LoffT)), 8},
		{"X__QaddrT", "size", unsafe.Sizeof(*new(X__QaddrT)), 8},
		{"X__QaddrT", "alignment", unsafe.Alignof(*new(X__QaddrT)), 8},
		{"X__CaddrT", "size", unsafe.Sizeof(*new(X__CaddrT)), 8},
		{"X__CaddrT", "alignment", unsafe.Alignof(*new(X__CaddrT)), 8},
		{"X__IntptrT", "size", unsafe.Sizeof(*new(X__IntptrT)), 8},
		{"X__IntptrT", "alignment", unsafe.Alignof(*new(X__IntptrT)), 8},
		{"X__SocklenT", "size", unsafe.Sizeof(*new(X__SocklenT)), 4},
		{"X__SocklenT", "alignment", unsafe.Alignof(*new(X__SocklenT)), 4},
		{"Iovec", "size", unsafe.Sizeof(*new(Iovec)), 16},
		{"Iovec", "alignment", unsafe.Alignof(*new(Iovec)), 8},
		{"Iovec", "IovBase offset", unsafe.Offsetof((*new(Iovec)).IovBase), 0},
		{"Iovec", "IovBase size", unsafe.Sizeof((*new(Iovec)).IovBase), 8},
		{"Iovec", "IovBase alignment", unsafe.Alignof((*new(Iovec)).IovBase), 8},
		{"Iovec", "IovLen offset", unsafe.Offsetof((*new(Iovec)).IovLen), 8},
		{"Iovec", "IovLen size", unsafe.Sizeof((*new(Iovec)).IovLen), 8},
		{"Iovec", "IovLen alignment", unsafe.Alignof((*new(Iovec)).IovLen), 8},
		{"SizeT", "size", unsafe.Sizeof(*new(SizeT)), 8},
		{"SizeT", "alignment", unsafe.Alignof(*new(SizeT)), 8},
		{"IfNameindexT", "size", unsafe.Sizeof(*new(IfNameindexT)), 16},
		{"IfNameindexT", "alignment", unsafe.Alignof(*new(IfNameindexT)), 8},
		{"IfNameindexT", "IfIndex offset", unsafe.Offsetof((*new(IfNameindexT)).IfIndex), 0},
		{"IfNameindexT", "IfIndex size", unsafe.Sizeof((*new(IfNameindexT)).IfIndex), 4},
		{"IfNameindexT", "IfIndex alignment", unsafe.Alignof((*new(IfNameindexT)).IfIndex), 4},
		{"IfNameindexT", "IfName offset", unsafe.Offsetof((*new(IfNameindexT)).IfName), 8},
		{"IfNameindexT", "IfName size", unsafe.Sizeof((*new(IfNameindexT)).IfName), 8},
		{"IfNameindexT", "IfName alignment", unsafe.Alignof((*new(IfNameindexT)).IfName), 8},
		{"Ifaddr", "size", unsafe.Sizeof(*new(Ifaddr)), 48},
		{"Ifaddr", "alignment", unsafe.Alignof(*new(Ifaddr)), 8},
		{"Ifaddr", "IfaAddr offset", unsafe.Offsetof((*new(Ifaddr)).IfaAddr), 0},
		{"Ifaddr", "IfaAddr size", unsafe.Sizeof((*new(Ifaddr)).IfaAddr), 16},
		{"Ifaddr", "IfaAddr alignment", unsafe.Alignof((*new(Ifaddr)).IfaAddr), 2},
		{"Ifaddr", "IfaIfu offset", unsafe.Offsetof((*new(Ifaddr)).IfaIfu), 16},
		{"Ifaddr", "IfaIfu size", unsafe.Sizeof((*new(Ifaddr)).IfaIfu), 16},
		{"Ifaddr", "IfaIfu alignment", unsafe.Alignof((*new(Ifaddr)).IfaIfu), 2},
		{"Ifaddr", "IfaIfp offset", unsafe.Offsetof((*new(Ifaddr)).IfaIfp), 32},
		{"Ifaddr", "IfaIfp size", unsafe.Sizeof((*new(Ifaddr)).IfaIfp), 8},
		{"Ifaddr", "IfaIfp alignment", unsafe.Alignof((*new(Ifaddr)).IfaIfp), 8},
		{"Ifaddr", "IfaNext offset", unsafe.Offsetof((*new(Ifaddr)).IfaNext), 40},
		{"Ifaddr", "IfaNext size", unsafe.Sizeof((*new(Ifaddr)).IfaNext), 8},
		{"Ifaddr", "IfaNext alignment", unsafe.Alignof((*new(Ifaddr)).IfaNext), 8},
		{"Ifmap", "size", unsafe.Sizeof(*new(Ifmap)), 24},
		{"Ifmap", "alignment", unsafe.Alignof(*new(Ifmap)), 8},
		{"Ifmap", "MemStart offset", unsafe.Offsetof((*new(Ifmap)).MemStart), 0},
		{"Ifmap", "MemStart size", unsafe.Sizeof((*new(Ifmap)).MemStart), 8},
		{"Ifmap", "MemStart alignment", unsafe.Alignof((*new(Ifmap)).MemStart), 8},
		{"Ifmap", "MemEnd offset", unsafe.Offsetof((*new(Ifmap)).MemEnd), 8},
		{"Ifmap", "MemEnd size", unsafe.Sizeof((*new(Ifmap)).MemEnd), 8},
		{"Ifmap", "MemEnd alignment", unsafe.Alignof((*new(Ifmap)).MemEnd), 8},
		{"Ifmap", "BaseAddr offset", unsafe.Offsetof((*new(Ifmap)).BaseAddr), 16},
		{"Ifmap", "BaseAddr size", unsafe.Sizeof((*new(Ifmap)).BaseAddr), 2},
		{"Ifmap", "BaseAddr alignment", unsafe.Alignof((*new(Ifmap)).BaseAddr), 2},
		{"Ifmap", "Irq offset", unsafe.Offsetof((*new(Ifmap)).Irq), 18},
		{"Ifmap", "Irq size", unsafe.Sizeof((*new(Ifmap)).Irq), 1},
		{"Ifmap", "Irq alignment", unsafe.Alignof((*new(Ifmap)).Irq), 1},
		{"Ifmap", "Dma offset", unsafe.Offsetof((*new(Ifmap)).Dma), 19},
		{"Ifmap", "Dma size", unsafe.Sizeof((*new(Ifmap)).Dma), 1},
		{"Ifmap", "Dma alignment", unsafe.Alignof((*new(Ifmap)).Dma), 1},
		{"Ifmap", "Port offset", unsafe.Offsetof((*new(Ifmap)).Port), 20},
		{"Ifmap", "Port size", unsafe.Sizeof((*new(Ifmap)).Port), 1},
		{"Ifmap", "Port alignment", unsafe.Alignof((*new(Ifmap)).Port), 1},
		{"Ifreq", "size", unsafe.Sizeof(*new(Ifreq)), 40},
		{"Ifreq", "alignment", unsafe.Alignof(*new(Ifreq)), 8},
		{"Ifreq", "IfrIfrn offset", unsafe.Offsetof((*new(Ifreq)).IfrIfrn), 0},
		{"Ifreq", "IfrIfrn size", unsafe.Sizeof((*new(Ifreq)).IfrIfrn), 16},
		{"Ifreq", "IfrIfrn alignment", unsafe.Alignof((*new(Ifreq)).IfrIfrn), 1},
		{"Ifreq", "IfrIfru offset", unsafe.Offsetof((*new(Ifreq)).IfrIfru), 16},
		{"Ifreq", "IfrIfru size", unsafe.Sizeof((*new(Ifreq)).IfrIfru), 24},
		{"Ifreq", "IfrIfru alignment", unsafe.Alignof((*new(Ifreq)).IfrIfru), 8},
		{"Ifconf", "size", unsafe.Sizeof(*new(Ifconf)), 16},
		{"Ifconf", "alignment", unsafe.Alignof(*new(Ifconf)), 8},
		{"Ifconf", "IfcLen offset", unsafe.Offsetof((*new(Ifconf)).IfcLen), 0},
		{"Ifconf", "IfcLen size", unsafe.Sizeof((*new(Ifconf)).IfcLen), 4},
		{"Ifconf", "IfcLen alignment", unsafe.Alignof((*new(Ifconf)).IfcLen), 4},
		{"Ifconf", "IfcIfcu offset", unsafe.Offsetof((*new(Ifconf)).IfcIfcu), 8},
		{"Ifconf", "IfcIfcu size", unsafe.Sizeof((*new(Ifconf)).IfcIfcu), 8},
		{"Ifconf", "IfcIfcu alignment", unsafe.Alignof((*new(Ifconf)).IfcIfcu), 8},
		{"Iface", "pointer size", unsafe.Sizeof((*Iface)(nil)), 8},
		{"Iface", "pointer alignment", unsafe.Alignof((*Iface)(nil)), 8},
		{"Hostent", "size", unsafe.Sizeof(*new(Hostent)), 32},
		{"Hostent", "alignment", unsafe.Alignof(*new(Hostent)), 8},
		{"Hostent", "Name offset", unsafe.Offsetof((*new(Hostent)).Name), 0},
		{"Hostent", "Name size", unsafe.Sizeof((*new(Hostent)).Name), 8},
		{"Hostent", "Name alignment", unsafe.Alignof((*new(Hostent)).Name), 8},
		{"Hostent", "Aliases offset", unsafe.Offsetof((*new(Hostent)).Aliases), 8},
		{"Hostent", "Aliases size", unsafe.Sizeof((*new(Hostent)).Aliases), 8},
		{"Hostent", "Aliases alignment", unsafe.Alignof((*new(Hostent)).Aliases), 8},
		{"Hostent", "AddrType offset", unsafe.Offsetof((*new(Hostent)).AddrType), 16},
		{"Hostent", "AddrType size", unsafe.Sizeof((*new(Hostent)).AddrType), 4},
		{"Hostent", "AddrType alignment", unsafe.Alignof((*new(Hostent)).AddrType), 4},
		{"Hostent", "Length offset", unsafe.Offsetof((*new(Hostent)).Length), 20},
		{"Hostent", "Length size", unsafe.Sizeof((*new(Hostent)).Length), 4},
		{"Hostent", "Length alignment", unsafe.Alignof((*new(Hostent)).Length), 4},
		{"Hostent", "AddrList offset", unsafe.Offsetof((*new(Hostent)).AddrList), 24},
		{"Hostent", "AddrList size", unsafe.Sizeof((*new(Hostent)).AddrList), 8},
		{"Hostent", "AddrList alignment", unsafe.Alignof((*new(Hostent)).AddrList), 8},
		{"Servent", "size", unsafe.Sizeof(*new(Servent)), 32},
		{"Servent", "alignment", unsafe.Alignof(*new(Servent)), 8},
		{"Servent", "SName offset", unsafe.Offsetof((*new(Servent)).SName), 0},
		{"Servent", "SName size", unsafe.Sizeof((*new(Servent)).SName), 8},
		{"Servent", "SName alignment", unsafe.Alignof((*new(Servent)).SName), 8},
		{"Servent", "SAliases offset", unsafe.Offsetof((*new(Servent)).SAliases), 8},
		{"Servent", "SAliases size", unsafe.Sizeof((*new(Servent)).SAliases), 8},
		{"Servent", "SAliases alignment", unsafe.Alignof((*new(Servent)).SAliases), 8},
		{"Servent", "SPort offset", unsafe.Offsetof((*new(Servent)).SPort), 16},
		{"Servent", "SPort size", unsafe.Sizeof((*new(Servent)).SPort), 4},
		{"Servent", "SPort alignment", unsafe.Alignof((*new(Servent)).SPort), 4},
		{"Servent", "SProto offset", unsafe.Offsetof((*new(Servent)).SProto), 24},
		{"Servent", "SProto size", unsafe.Sizeof((*new(Servent)).SProto), 8},
		{"Servent", "SProto alignment", unsafe.Alignof((*new(Servent)).SProto), 8},
		{"Protoent", "size", unsafe.Sizeof(*new(Protoent)), 24},
		{"Protoent", "alignment", unsafe.Alignof(*new(Protoent)), 8},
		{"Protoent", "PName offset", unsafe.Offsetof((*new(Protoent)).PName), 0},
		{"Protoent", "PName size", unsafe.Sizeof((*new(Protoent)).PName), 8},
		{"Protoent", "PName alignment", unsafe.Alignof((*new(Protoent)).PName), 8},
		{"Protoent", "PAliases offset", unsafe.Offsetof((*new(Protoent)).PAliases), 8},
		{"Protoent", "PAliases size", unsafe.Sizeof((*new(Protoent)).PAliases), 8},
		{"Protoent", "PAliases alignment", unsafe.Alignof((*new(Protoent)).PAliases), 8},
		{"Protoent", "PProto offset", unsafe.Offsetof((*new(Protoent)).PProto), 16},
		{"Protoent", "PProto size", unsafe.Sizeof((*new(Protoent)).PProto), 4},
		{"Protoent", "PProto alignment", unsafe.Alignof((*new(Protoent)).PProto), 4},
		{"AddrInfo", "size", unsafe.Sizeof(*new(AddrInfo)), 48},
		{"AddrInfo", "alignment", unsafe.Alignof(*new(AddrInfo)), 8},
		{"AddrInfo", "Flags offset", unsafe.Offsetof((*new(AddrInfo)).Flags), 0},
		{"AddrInfo", "Flags size", unsafe.Sizeof((*new(AddrInfo)).Flags), 4},
		{"AddrInfo", "Flags alignment", unsafe.Alignof((*new(AddrInfo)).Flags), 4},
		{"AddrInfo", "Family offset", unsafe.Offsetof((*new(AddrInfo)).Family), 4},
		{"AddrInfo", "Family size", unsafe.Sizeof((*new(AddrInfo)).Family), 4},
		{"AddrInfo", "Family alignment", unsafe.Alignof((*new(AddrInfo)).Family), 4},
		{"AddrInfo", "SockType offset", unsafe.Offsetof((*new(AddrInfo)).SockType), 8},
		{"AddrInfo", "SockType size", unsafe.Sizeof((*new(AddrInfo)).SockType), 4},
		{"AddrInfo", "SockType alignment", unsafe.Alignof((*new(AddrInfo)).SockType), 4},
		{"AddrInfo", "Protocol offset", unsafe.Offsetof((*new(AddrInfo)).Protocol), 12},
		{"AddrInfo", "Protocol size", unsafe.Sizeof((*new(AddrInfo)).Protocol), 4},
		{"AddrInfo", "Protocol alignment", unsafe.Alignof((*new(AddrInfo)).Protocol), 4},
		{"AddrInfo", "AddrLen offset", unsafe.Offsetof((*new(AddrInfo)).AddrLen), 16},
		{"AddrInfo", "AddrLen size", unsafe.Sizeof((*new(AddrInfo)).AddrLen), 4},
		{"AddrInfo", "AddrLen alignment", unsafe.Alignof((*new(AddrInfo)).AddrLen), 4},
		{"AddrInfo", "Addr offset", unsafe.Offsetof((*new(AddrInfo)).Addr), 24},
		{"AddrInfo", "Addr size", unsafe.Sizeof((*new(AddrInfo)).Addr), 8},
		{"AddrInfo", "Addr alignment", unsafe.Alignof((*new(AddrInfo)).Addr), 8},
		{"AddrInfo", "CanOnName offset", unsafe.Offsetof((*new(AddrInfo)).CanOnName), 32},
		{"AddrInfo", "CanOnName size", unsafe.Sizeof((*new(AddrInfo)).CanOnName), 8},
		{"AddrInfo", "CanOnName alignment", unsafe.Alignof((*new(AddrInfo)).CanOnName), 8},
		{"AddrInfo", "Next offset", unsafe.Offsetof((*new(AddrInfo)).Next), 40},
		{"AddrInfo", "Next size", unsafe.Sizeof((*new(AddrInfo)).Next), 8},
		{"AddrInfo", "Next alignment", unsafe.Alignof((*new(AddrInfo)).Next), 8},
		{"InAddr", "size", unsafe.Sizeof(*new(InAddr)), 4},
		{"InAddr", "alignment", unsafe.Alignof(*new(InAddr)), 4},
		{"InAddr", "Addr offset", unsafe.Offsetof((*new(InAddr)).Addr), 0},
		{"InAddr", "Addr size", unsafe.Sizeof((*new(InAddr)).Addr), 4},
		{"InAddr", "Addr alignment", unsafe.Alignof((*new(InAddr)).Addr), 4},
		{"In6Addr", "size", unsafe.Sizeof(*new(In6Addr)), 16},
		{"In6Addr", "alignment", unsafe.Alignof(*new(In6Addr)), 4},
		{"In6Addr", "U6Addr offset", unsafe.Offsetof((*new(In6Addr)).U6Addr), 0},
		{"In6Addr", "U6Addr size", unsafe.Sizeof((*new(In6Addr)).U6Addr), 16},
		{"In6Addr", "U6Addr alignment", unsafe.Alignof((*new(In6Addr)).U6Addr), 1},
		{"SockaddrIn", "size", unsafe.Sizeof(*new(SockaddrIn)), 16},
		{"SockaddrIn", "alignment", unsafe.Alignof(*new(SockaddrIn)), 4},
		{"SockaddrIn", "Family offset", unsafe.Offsetof((*new(SockaddrIn)).Family), 0},
		{"SockaddrIn", "Family size", unsafe.Sizeof((*new(SockaddrIn)).Family), 2},
		{"SockaddrIn", "Family alignment", unsafe.Alignof((*new(SockaddrIn)).Family), 2},
		{"SockaddrIn", "Port offset", unsafe.Offsetof((*new(SockaddrIn)).Port), 2},
		{"SockaddrIn", "Port size", unsafe.Sizeof((*new(SockaddrIn)).Port), 2},
		{"SockaddrIn", "Port alignment", unsafe.Alignof((*new(SockaddrIn)).Port), 2},
		{"SockaddrIn", "Addr offset", unsafe.Offsetof((*new(SockaddrIn)).Addr), 4},
		{"SockaddrIn", "Addr size", unsafe.Sizeof((*new(SockaddrIn)).Addr), 4},
		{"SockaddrIn", "Addr alignment", unsafe.Alignof((*new(SockaddrIn)).Addr), 4},
		{"SockaddrIn", "Zero offset", unsafe.Offsetof((*new(SockaddrIn)).Zero), 8},
		{"SockaddrIn", "Zero size", unsafe.Sizeof((*new(SockaddrIn)).Zero), 8},
		{"SockaddrIn", "Zero alignment", unsafe.Alignof((*new(SockaddrIn)).Zero), 1},
		{"SockaddrIn6", "size", unsafe.Sizeof(*new(SockaddrIn6)), 28},
		{"SockaddrIn6", "alignment", unsafe.Alignof(*new(SockaddrIn6)), 4},
		{"SockaddrIn6", "Family offset", unsafe.Offsetof((*new(SockaddrIn6)).Family), 0},
		{"SockaddrIn6", "Family size", unsafe.Sizeof((*new(SockaddrIn6)).Family), 2},
		{"SockaddrIn6", "Family alignment", unsafe.Alignof((*new(SockaddrIn6)).Family), 2},
		{"SockaddrIn6", "Port offset", unsafe.Offsetof((*new(SockaddrIn6)).Port), 2},
		{"SockaddrIn6", "Port size", unsafe.Sizeof((*new(SockaddrIn6)).Port), 2},
		{"SockaddrIn6", "Port alignment", unsafe.Alignof((*new(SockaddrIn6)).Port), 2},
		{"SockaddrIn6", "Flowinfo offset", unsafe.Offsetof((*new(SockaddrIn6)).Flowinfo), 4},
		{"SockaddrIn6", "Flowinfo size", unsafe.Sizeof((*new(SockaddrIn6)).Flowinfo), 4},
		{"SockaddrIn6", "Flowinfo alignment", unsafe.Alignof((*new(SockaddrIn6)).Flowinfo), 4},
		{"SockaddrIn6", "Addr offset", unsafe.Offsetof((*new(SockaddrIn6)).Addr), 8},
		{"SockaddrIn6", "Addr size", unsafe.Sizeof((*new(SockaddrIn6)).Addr), 16},
		{"SockaddrIn6", "Addr alignment", unsafe.Alignof((*new(SockaddrIn6)).Addr), 4},
		{"SockaddrIn6", "ScopeId offset", unsafe.Offsetof((*new(SockaddrIn6)).ScopeId), 24},
		{"SockaddrIn6", "ScopeId size", unsafe.Sizeof((*new(SockaddrIn6)).ScopeId), 4},
		{"SockaddrIn6", "ScopeId alignment", unsafe.Alignof((*new(SockaddrIn6)).ScopeId), 4},
		{"IpMreq", "size", unsafe.Sizeof(*new(IpMreq)), 8},
		{"IpMreq", "alignment", unsafe.Alignof(*new(IpMreq)), 4},
		{"IpMreq", "ImrMultiaddr offset", unsafe.Offsetof((*new(IpMreq)).ImrMultiaddr), 0},
		{"IpMreq", "ImrMultiaddr size", unsafe.Sizeof((*new(IpMreq)).ImrMultiaddr), 4},
		{"IpMreq", "ImrMultiaddr alignment", unsafe.Alignof((*new(IpMreq)).ImrMultiaddr), 4},
		{"IpMreq", "ImrInterface offset", unsafe.Offsetof((*new(IpMreq)).ImrInterface), 4},
		{"IpMreq", "ImrInterface size", unsafe.Sizeof((*new(IpMreq)).ImrInterface), 4},
		{"IpMreq", "ImrInterface alignment", unsafe.Alignof((*new(IpMreq)).ImrInterface), 4},
		{"IpMreqSource", "size", unsafe.Sizeof(*new(IpMreqSource)), 12},
		{"IpMreqSource", "alignment", unsafe.Alignof(*new(IpMreqSource)), 4},
		{"IpMreqSource", "ImrMultiaddr offset", unsafe.Offsetof((*new(IpMreqSource)).ImrMultiaddr), 0},
		{"IpMreqSource", "ImrMultiaddr size", unsafe.Sizeof((*new(IpMreqSource)).ImrMultiaddr), 4},
		{"IpMreqSource", "ImrMultiaddr alignment", unsafe.Alignof((*new(IpMreqSource)).ImrMultiaddr), 4},
		{"IpMreqSource", "ImrInterface offset", unsafe.Offsetof((*new(IpMreqSource)).ImrInterface), 4},
		{"IpMreqSource", "ImrInterface size", unsafe.Sizeof((*new(IpMreqSource)).ImrInterface), 4},
		{"IpMreqSource", "ImrInterface alignment", unsafe.Alignof((*new(IpMreqSource)).ImrInterface), 4},
		{"IpMreqSource", "ImrSourceaddr offset", unsafe.Offsetof((*new(IpMreqSource)).ImrSourceaddr), 8},
		{"IpMreqSource", "ImrSourceaddr size", unsafe.Sizeof((*new(IpMreqSource)).ImrSourceaddr), 4},
		{"IpMreqSource", "ImrSourceaddr alignment", unsafe.Alignof((*new(IpMreqSource)).ImrSourceaddr), 4},
		{"Ipv6Mreq", "size", unsafe.Sizeof(*new(Ipv6Mreq)), 20},
		{"Ipv6Mreq", "alignment", unsafe.Alignof(*new(Ipv6Mreq)), 4},
		{"Ipv6Mreq", "Ipv6mrMultiaddr offset", unsafe.Offsetof((*new(Ipv6Mreq)).Ipv6mrMultiaddr), 0},
		{"Ipv6Mreq", "Ipv6mrMultiaddr size", unsafe.Sizeof((*new(Ipv6Mreq)).Ipv6mrMultiaddr), 16},
		{"Ipv6Mreq", "Ipv6mrMultiaddr alignment", unsafe.Alignof((*new(Ipv6Mreq)).Ipv6mrMultiaddr), 4},
		{"Ipv6Mreq", "Ipv6mrInterface offset", unsafe.Offsetof((*new(Ipv6Mreq)).Ipv6mrInterface), 16},
		{"Ipv6Mreq", "Ipv6mrInterface size", unsafe.Sizeof((*new(Ipv6Mreq)).Ipv6mrInterface), 4},
		{"Ipv6Mreq", "Ipv6mrInterface alignment", unsafe.Alignof((*new(Ipv6Mreq)).Ipv6mrInterface), 4},
		{"GroupReq", "size", unsafe.Sizeof(*new(GroupReq)), 136},
		{"GroupReq", "alignment", unsafe.Alignof(*new(GroupReq)), 8},
		{"GroupReq", "GrInterface offset", unsafe.Offsetof((*new(GroupReq)).GrInterface), 0},
		{"GroupReq", "GrInterface size", unsafe.Sizeof((*new(GroupReq)).GrInterface), 4},
		{"GroupReq", "GrInterface alignment", unsafe.Alignof((*new(GroupReq)).GrInterface), 4},
		{"GroupReq", "GrGroup offset", unsafe.Offsetof((*new(GroupReq)).GrGroup), 8},
		{"GroupReq", "GrGroup size", unsafe.Sizeof((*new(GroupReq)).GrGroup), 128},
		{"GroupReq", "GrGroup alignment", unsafe.Alignof((*new(GroupReq)).GrGroup), 8},
		{"GroupSourceReq", "size", unsafe.Sizeof(*new(GroupSourceReq)), 264},
		{"GroupSourceReq", "alignment", unsafe.Alignof(*new(GroupSourceReq)), 8},
		{"GroupSourceReq", "GsrInterface offset", unsafe.Offsetof((*new(GroupSourceReq)).GsrInterface), 0},
		{"GroupSourceReq", "GsrInterface size", unsafe.Sizeof((*new(GroupSourceReq)).GsrInterface), 4},
		{"GroupSourceReq", "GsrInterface alignment", unsafe.Alignof((*new(GroupSourceReq)).GsrInterface), 4},
		{"GroupSourceReq", "GsrGroup offset", unsafe.Offsetof((*new(GroupSourceReq)).GsrGroup), 8},
		{"GroupSourceReq", "GsrGroup size", unsafe.Sizeof((*new(GroupSourceReq)).GsrGroup), 128},
		{"GroupSourceReq", "GsrGroup alignment", unsafe.Alignof((*new(GroupSourceReq)).GsrGroup), 8},
		{"GroupSourceReq", "GsrSource offset", unsafe.Offsetof((*new(GroupSourceReq)).GsrSource), 136},
		{"GroupSourceReq", "GsrSource size", unsafe.Sizeof((*new(GroupSourceReq)).GsrSource), 128},
		{"GroupSourceReq", "GsrSource alignment", unsafe.Alignof((*new(GroupSourceReq)).GsrSource), 8},
		{"IpMsfilter", "size", unsafe.Sizeof(*new(IpMsfilter)), 20},
		{"IpMsfilter", "alignment", unsafe.Alignof(*new(IpMsfilter)), 4},
		{"IpMsfilter", "ImsfMultiaddr offset", unsafe.Offsetof((*new(IpMsfilter)).ImsfMultiaddr), 0},
		{"IpMsfilter", "ImsfMultiaddr size", unsafe.Sizeof((*new(IpMsfilter)).ImsfMultiaddr), 4},
		{"IpMsfilter", "ImsfMultiaddr alignment", unsafe.Alignof((*new(IpMsfilter)).ImsfMultiaddr), 4},
		{"IpMsfilter", "ImsfInterface offset", unsafe.Offsetof((*new(IpMsfilter)).ImsfInterface), 4},
		{"IpMsfilter", "ImsfInterface size", unsafe.Sizeof((*new(IpMsfilter)).ImsfInterface), 4},
		{"IpMsfilter", "ImsfInterface alignment", unsafe.Alignof((*new(IpMsfilter)).ImsfInterface), 4},
		{"IpMsfilter", "ImsfFmode offset", unsafe.Offsetof((*new(IpMsfilter)).ImsfFmode), 8},
		{"IpMsfilter", "ImsfFmode size", unsafe.Sizeof((*new(IpMsfilter)).ImsfFmode), 4},
		{"IpMsfilter", "ImsfFmode alignment", unsafe.Alignof((*new(IpMsfilter)).ImsfFmode), 4},
		{"IpMsfilter", "ImsfNumsrc offset", unsafe.Offsetof((*new(IpMsfilter)).ImsfNumsrc), 12},
		{"IpMsfilter", "ImsfNumsrc size", unsafe.Sizeof((*new(IpMsfilter)).ImsfNumsrc), 4},
		{"IpMsfilter", "ImsfNumsrc alignment", unsafe.Alignof((*new(IpMsfilter)).ImsfNumsrc), 4},
		{"IpMsfilter", "ImsfSlist offset", unsafe.Offsetof((*new(IpMsfilter)).ImsfSlist), 16},
		{"IpMsfilter", "ImsfSlist size", unsafe.Sizeof((*new(IpMsfilter)).ImsfSlist), 4},
		{"IpMsfilter", "ImsfSlist alignment", unsafe.Alignof((*new(IpMsfilter)).ImsfSlist), 4},
		{"GroupFilter", "size", unsafe.Sizeof(*new(GroupFilter)), 272},
		{"GroupFilter", "alignment", unsafe.Alignof(*new(GroupFilter)), 8},
		{"GroupFilter", "GfInterface offset", unsafe.Offsetof((*new(GroupFilter)).GfInterface), 0},
		{"GroupFilter", "GfInterface size", unsafe.Sizeof((*new(GroupFilter)).GfInterface), 4},
		{"GroupFilter", "GfInterface alignment", unsafe.Alignof((*new(GroupFilter)).GfInterface), 4},
		{"GroupFilter", "GfGroup offset", unsafe.Offsetof((*new(GroupFilter)).GfGroup), 8},
		{"GroupFilter", "GfGroup size", unsafe.Sizeof((*new(GroupFilter)).GfGroup), 128},
		{"GroupFilter", "GfGroup alignment", unsafe.Alignof((*new(GroupFilter)).GfGroup), 8},
		{"GroupFilter", "GfFmode offset", unsafe.Offsetof((*new(GroupFilter)).GfFmode), 136},
		{"GroupFilter", "GfFmode size", unsafe.Sizeof((*new(GroupFilter)).GfFmode), 4},
		{"GroupFilter", "GfFmode alignment", unsafe.Alignof((*new(GroupFilter)).GfFmode), 4},
		{"GroupFilter", "GfNumsrc offset", unsafe.Offsetof((*new(GroupFilter)).GfNumsrc), 140},
		{"GroupFilter", "GfNumsrc size", unsafe.Sizeof((*new(GroupFilter)).GfNumsrc), 4},
		{"GroupFilter", "GfNumsrc alignment", unsafe.Alignof((*new(GroupFilter)).GfNumsrc), 4},
		{"GroupFilter", "GfSlist offset", unsafe.Offsetof((*new(GroupFilter)).GfSlist), 144},
		{"GroupFilter", "GfSlist size", unsafe.Sizeof((*new(GroupFilter)).GfSlist), 128},
		{"GroupFilter", "GfSlist alignment", unsafe.Alignof((*new(GroupFilter)).GfSlist), 8},
		{"InAddrT", "size", unsafe.Sizeof(*new(InAddrT)), 4},
		{"InAddrT", "alignment", unsafe.Alignof(*new(InAddrT)), 4},
		{"InPortT", "size", unsafe.Sizeof(*new(InPortT)), 2},
		{"InPortT", "alignment", unsafe.Alignof(*new(InPortT)), 2},
		{"Rpcent", "size", unsafe.Sizeof(*new(Rpcent)), 24},
		{"Rpcent", "alignment", unsafe.Alignof(*new(Rpcent)), 8},
		{"Rpcent", "RName offset", unsafe.Offsetof((*new(Rpcent)).RName), 0},
		{"Rpcent", "RName size", unsafe.Sizeof((*new(Rpcent)).RName), 8},
		{"Rpcent", "RName alignment", unsafe.Alignof((*new(Rpcent)).RName), 8},
		{"Rpcent", "RAliases offset", unsafe.Offsetof((*new(Rpcent)).RAliases), 8},
		{"Rpcent", "RAliases size", unsafe.Sizeof((*new(Rpcent)).RAliases), 8},
		{"Rpcent", "RAliases alignment", unsafe.Alignof((*new(Rpcent)).RAliases), 8},
		{"Rpcent", "RNumber offset", unsafe.Offsetof((*new(Rpcent)).RNumber), 16},
		{"Rpcent", "RNumber size", unsafe.Sizeof((*new(Rpcent)).RNumber), 4},
		{"Rpcent", "RNumber alignment", unsafe.Alignof((*new(Rpcent)).RNumber), 4},
		{"Uint16T", "size", unsafe.Sizeof(*new(Uint16T)), 2},
		{"Uint16T", "alignment", unsafe.Alignof(*new(Uint16T)), 2},
		{"Uint32T", "size", unsafe.Sizeof(*new(Uint32T)), 4},
		{"Uint32T", "alignment", unsafe.Alignof(*new(Uint32T)), 4},
		{"Osockaddr", "size", unsafe.Sizeof(*new(Osockaddr)), 16},
		{"Osockaddr", "alignment", unsafe.Alignof(*new(Osockaddr)), 2},
		{"Osockaddr", "SaFamily offset", unsafe.Offsetof((*new(Osockaddr)).SaFamily), 0},
		{"Osockaddr", "SaFamily size", unsafe.Sizeof((*new(Osockaddr)).SaFamily), 2},
		{"Osockaddr", "SaFamily alignment", unsafe.Alignof((*new(Osockaddr)).SaFamily), 2},
		{"Osockaddr", "SaData offset", unsafe.Offsetof((*new(Osockaddr)).SaData), 2},
		{"Osockaddr", "SaData size", unsafe.Sizeof((*new(Osockaddr)).SaData), 14},
		{"Osockaddr", "SaData alignment", unsafe.Alignof((*new(Osockaddr)).SaData), 1},
		{"SsizeT", "size", unsafe.Sizeof(*new(SsizeT)), 8},
		{"SsizeT", "alignment", unsafe.Alignof(*new(SsizeT)), 8},
		{"SockaddrUn", "size", unsafe.Sizeof(*new(SockaddrUn)), 110},
		{"SockaddrUn", "alignment", unsafe.Alignof(*new(SockaddrUn)), 2},
		{"SockaddrUn", "SunFamily offset", unsafe.Offsetof((*new(SockaddrUn)).SunFamily), 0},
		{"SockaddrUn", "SunFamily size", unsafe.Sizeof((*new(SockaddrUn)).SunFamily), 2},
		{"SockaddrUn", "SunFamily alignment", unsafe.Alignof((*new(SockaddrUn)).SunFamily), 2},
		{"SockaddrUn", "SunPath offset", unsafe.Offsetof((*new(SockaddrUn)).SunPath), 2},
		{"SockaddrUn", "SunPath size", unsafe.Sizeof((*new(SockaddrUn)).SunPath), 108},
		{"SockaddrUn", "SunPath alignment", unsafe.Alignof((*new(SockaddrUn)).SunPath), 1},
		{"In6addrAny", "global size", unsafe.Sizeof(In6addrAny), 16},
		{"In6addrAny", "global alignment", unsafe.Alignof(In6addrAny), 4},
		{"In6addrLoopback", "global size", unsafe.Sizeof(In6addrLoopback), 16},
		{"In6addrLoopback", "global alignment", unsafe.Alignof(In6addrLoopback), 4},
	}
	var networkIntegerConstants = []struct {
		name      string
		got, want int64
	}{
		{"X_ARPA_INET_H", int64(X_ARPA_INET_H), 1},
		{"SOL_SOCKET", int64(SOL_SOCKET), 1},
		{"SO_DEBUG", int64(SO_DEBUG), 1},
		{"SO_REUSEADDR", int64(SO_REUSEADDR), 2},
		{"SO_TYPE", int64(SO_TYPE), 3},
		{"SO_ERROR", int64(SO_ERROR), 4},
		{"SO_DONTROUTE", int64(SO_DONTROUTE), 5},
		{"SO_BROADCAST", int64(SO_BROADCAST), 6},
		{"SO_SNDBUF", int64(SO_SNDBUF), 7},
		{"SO_RCVBUF", int64(SO_RCVBUF), 8},
		{"SO_SNDBUFFORCE", int64(SO_SNDBUFFORCE), 32},
		{"SO_RCVBUFFORCE", int64(SO_RCVBUFFORCE), 33},
		{"SO_KEEPALIVE", int64(SO_KEEPALIVE), 9},
		{"SO_OOBINLINE", int64(SO_OOBINLINE), 10},
		{"SO_NO_CHECK", int64(SO_NO_CHECK), 11},
		{"SO_PRIORITY", int64(SO_PRIORITY), 12},
		{"SO_LINGER", int64(SO_LINGER), 13},
		{"SO_BSDCOMPAT", int64(SO_BSDCOMPAT), 14},
		{"SO_REUSEPORT", int64(SO_REUSEPORT), 15},
		{"SO_PASSCRED", int64(SO_PASSCRED), 16},
		{"SO_PEERCRED", int64(SO_PEERCRED), 17},
		{"SO_RCVLOWAT", int64(SO_RCVLOWAT), 18},
		{"SO_SNDLOWAT", int64(SO_SNDLOWAT), 19},
		{"SO_RCVTIMEO", int64(SO_RCVTIMEO), 20},
		{"SO_SNDTIMEO", int64(SO_SNDTIMEO), 21},
		{"SO_SECURITY_AUTHENTICATION", int64(SO_SECURITY_AUTHENTICATION), 22},
		{"SO_SECURITY_ENCRYPTION_TRANSPORT", int64(SO_SECURITY_ENCRYPTION_TRANSPORT), 23},
		{"SO_SECURITY_ENCRYPTION_NETWORK", int64(SO_SECURITY_ENCRYPTION_NETWORK), 24},
		{"SO_BINDTODEVICE", int64(SO_BINDTODEVICE), 25},
		{"SO_ATTACH_FILTER", int64(SO_ATTACH_FILTER), 26},
		{"SO_DETACH_FILTER", int64(SO_DETACH_FILTER), 27},
		{"SO_GET_FILTER", int64(SO_GET_FILTER), 26},
		{"SO_PEERNAME", int64(SO_PEERNAME), 28},
		{"SO_TIMESTAMP", int64(SO_TIMESTAMP), 29},
		{"SCM_TIMESTAMP", int64(SCM_TIMESTAMP), 29},
		{"SO_ACCEPTCONN", int64(SO_ACCEPTCONN), 30},
		{"SO_PEERSEC", int64(SO_PEERSEC), 31},
		{"SO_PASSSEC", int64(SO_PASSSEC), 34},
		{"SO_TIMESTAMPNS", int64(SO_TIMESTAMPNS), 35},
		{"SCM_TIMESTAMPNS", int64(SCM_TIMESTAMPNS), 35},
		{"SO_MARK", int64(SO_MARK), 36},
		{"SO_TIMESTAMPING", int64(SO_TIMESTAMPING), 37},
		{"SCM_TIMESTAMPING", int64(SCM_TIMESTAMPING), 37},
		{"SO_PROTOCOL", int64(SO_PROTOCOL), 38},
		{"SO_DOMAIN", int64(SO_DOMAIN), 39},
		{"SO_RXQ_OVFL", int64(SO_RXQ_OVFL), 40},
		{"SO_WIFI_STATUS", int64(SO_WIFI_STATUS), 41},
		{"SCM_WIFI_STATUS", int64(SCM_WIFI_STATUS), 41},
		{"SO_PEEK_OFF", int64(SO_PEEK_OFF), 42},
		{"SO_NOFCS", int64(SO_NOFCS), 43},
		{"SO_LOCK_FILTER", int64(SO_LOCK_FILTER), 44},
		{"SO_SELECT_ERR_QUEUE", int64(SO_SELECT_ERR_QUEUE), 45},
		{"FIOSETOWN", int64(FIOSETOWN), 35073},
		{"SIOCSPGRP", int64(SIOCSPGRP), 35074},
		{"FIOGETOWN", int64(FIOGETOWN), 35075},
		{"SIOCGPGRP", int64(SIOCGPGRP), 35076},
		{"SIOCATMARK", int64(SIOCATMARK), 35077},
		{"SIOCGSTAMP", int64(SIOCGSTAMP), 35078},
		{"SIOCGSTAMPNS", int64(SIOCGSTAMPNS), 35079},
		{"X__USE_KERNEL_IPV6_DEFS", int64(X__USE_KERNEL_IPV6_DEFS), 0},
		{"IP_OPTIONS", int64(IP_OPTIONS), 4},
		{"IP_HDRINCL", int64(IP_HDRINCL), 3},
		{"IP_TOS", int64(IP_TOS), 1},
		{"IP_TTL", int64(IP_TTL), 2},
		{"IP_RECVOPTS", int64(IP_RECVOPTS), 6},
		{"IP_RETOPTS", int64(IP_RETOPTS), 7},
		{"IP_MULTICAST_IF", int64(IP_MULTICAST_IF), 32},
		{"IP_MULTICAST_TTL", int64(IP_MULTICAST_TTL), 33},
		{"IP_MULTICAST_LOOP", int64(IP_MULTICAST_LOOP), 34},
		{"IP_ADD_MEMBERSHIP", int64(IP_ADD_MEMBERSHIP), 35},
		{"IP_DROP_MEMBERSHIP", int64(IP_DROP_MEMBERSHIP), 36},
		{"IP_UNBLOCK_SOURCE", int64(IP_UNBLOCK_SOURCE), 37},
		{"IP_BLOCK_SOURCE", int64(IP_BLOCK_SOURCE), 38},
		{"IP_ADD_SOURCE_MEMBERSHIP", int64(IP_ADD_SOURCE_MEMBERSHIP), 39},
		{"IP_DROP_SOURCE_MEMBERSHIP", int64(IP_DROP_SOURCE_MEMBERSHIP), 40},
		{"IP_MSFILTER", int64(IP_MSFILTER), 41},
		{"MCAST_JOIN_GROUP", int64(MCAST_JOIN_GROUP), 42},
		{"MCAST_BLOCK_SOURCE", int64(MCAST_BLOCK_SOURCE), 43},
		{"MCAST_UNBLOCK_SOURCE", int64(MCAST_UNBLOCK_SOURCE), 44},
		{"MCAST_LEAVE_GROUP", int64(MCAST_LEAVE_GROUP), 45},
		{"MCAST_JOIN_SOURCE_GROUP", int64(MCAST_JOIN_SOURCE_GROUP), 46},
		{"MCAST_LEAVE_SOURCE_GROUP", int64(MCAST_LEAVE_SOURCE_GROUP), 47},
		{"MCAST_MSFILTER", int64(MCAST_MSFILTER), 48},
		{"IP_MULTICAST_ALL", int64(IP_MULTICAST_ALL), 49},
		{"IP_UNICAST_IF", int64(IP_UNICAST_IF), 50},
		{"MCAST_EXCLUDE", int64(MCAST_EXCLUDE), 0},
		{"MCAST_INCLUDE", int64(MCAST_INCLUDE), 1},
		{"IP_ROUTER_ALERT", int64(IP_ROUTER_ALERT), 5},
		{"IP_PKTINFO", int64(IP_PKTINFO), 8},
		{"IP_PKTOPTIONS", int64(IP_PKTOPTIONS), 9},
		{"IP_PMTUDISC", int64(IP_PMTUDISC), 10},
		{"IP_MTU_DISCOVER", int64(IP_MTU_DISCOVER), 10},
		{"IP_RECVERR", int64(IP_RECVERR), 11},
		{"IP_RECVTTL", int64(IP_RECVTTL), 12},
		{"IP_RECVTOS", int64(IP_RECVTOS), 13},
		{"IP_MTU", int64(IP_MTU), 14},
		{"IP_FREEBIND", int64(IP_FREEBIND), 15},
		{"IP_IPSEC_POLICY", int64(IP_IPSEC_POLICY), 16},
		{"IP_XFRM_POLICY", int64(IP_XFRM_POLICY), 17},
		{"IP_PASSSEC", int64(IP_PASSSEC), 18},
		{"IP_TRANSPARENT", int64(IP_TRANSPARENT), 19},
		{"IP_ORIGDSTADDR", int64(IP_ORIGDSTADDR), 20},
		{"IP_RECVORIGDSTADDR", int64(IP_RECVORIGDSTADDR), 20},
		{"IP_MINTTL", int64(IP_MINTTL), 21},
		{"IP_NODEFRAG", int64(IP_NODEFRAG), 22},
		{"IP_CHECKSUM", int64(IP_CHECKSUM), 23},
		{"IP_BIND_ADDRESS_NO_PORT", int64(IP_BIND_ADDRESS_NO_PORT), 24},
		{"IP_PMTUDISC_DONT", int64(IP_PMTUDISC_DONT), 0},
		{"IP_PMTUDISC_WANT", int64(IP_PMTUDISC_WANT), 1},
		{"IP_PMTUDISC_DO", int64(IP_PMTUDISC_DO), 2},
		{"IP_PMTUDISC_PROBE", int64(IP_PMTUDISC_PROBE), 3},
		{"IP_PMTUDISC_INTERFACE", int64(IP_PMTUDISC_INTERFACE), 4},
		{"IP_PMTUDISC_OMIT", int64(IP_PMTUDISC_OMIT), 5},
		{"SOL_IP", int64(SOL_IP), 0},
		{"IP_DEFAULT_MULTICAST_TTL", int64(IP_DEFAULT_MULTICAST_TTL), 1},
		{"IP_DEFAULT_MULTICAST_LOOP", int64(IP_DEFAULT_MULTICAST_LOOP), 1},
		{"IP_MAX_MEMBERSHIPS", int64(IP_MAX_MEMBERSHIPS), 20},
		{"IPV6_ADDRFORM", int64(IPV6_ADDRFORM), 1},
		{"IPV6_2292PKTINFO", int64(IPV6_2292PKTINFO), 2},
		{"IPV6_2292HOPOPTS", int64(IPV6_2292HOPOPTS), 3},
		{"IPV6_2292DSTOPTS", int64(IPV6_2292DSTOPTS), 4},
		{"IPV6_2292RTHDR", int64(IPV6_2292RTHDR), 5},
		{"IPV6_2292PKTOPTIONS", int64(IPV6_2292PKTOPTIONS), 6},
		{"IPV6_CHECKSUM", int64(IPV6_CHECKSUM), 7},
		{"IPV6_2292HOPLIMIT", int64(IPV6_2292HOPLIMIT), 8},
		{"IPV6_NEXTHOP", int64(IPV6_NEXTHOP), 9},
		{"IPV6_AUTHHDR", int64(IPV6_AUTHHDR), 10},
		{"IPV6_UNICAST_HOPS", int64(IPV6_UNICAST_HOPS), 16},
		{"IPV6_MULTICAST_IF", int64(IPV6_MULTICAST_IF), 17},
		{"IPV6_MULTICAST_HOPS", int64(IPV6_MULTICAST_HOPS), 18},
		{"IPV6_MULTICAST_LOOP", int64(IPV6_MULTICAST_LOOP), 19},
		{"IPV6_JOIN_GROUP", int64(IPV6_JOIN_GROUP), 20},
		{"IPV6_LEAVE_GROUP", int64(IPV6_LEAVE_GROUP), 21},
		{"IPV6_ROUTER_ALERT", int64(IPV6_ROUTER_ALERT), 22},
		{"IPV6_MTU_DISCOVER", int64(IPV6_MTU_DISCOVER), 23},
		{"IPV6_MTU", int64(IPV6_MTU), 24},
		{"IPV6_RECVERR", int64(IPV6_RECVERR), 25},
		{"IPV6_V6ONLY", int64(IPV6_V6ONLY), 26},
		{"IPV6_JOIN_ANYCAST", int64(IPV6_JOIN_ANYCAST), 27},
		{"IPV6_LEAVE_ANYCAST", int64(IPV6_LEAVE_ANYCAST), 28},
		{"IPV6_IPSEC_POLICY", int64(IPV6_IPSEC_POLICY), 34},
		{"IPV6_XFRM_POLICY", int64(IPV6_XFRM_POLICY), 35},
		{"IPV6_HDRINCL", int64(IPV6_HDRINCL), 36},
		{"IPV6_RECVPKTINFO", int64(IPV6_RECVPKTINFO), 49},
		{"IPV6_PKTINFO", int64(IPV6_PKTINFO), 50},
		{"IPV6_RECVHOPLIMIT", int64(IPV6_RECVHOPLIMIT), 51},
		{"IPV6_HOPLIMIT", int64(IPV6_HOPLIMIT), 52},
		{"IPV6_RECVHOPOPTS", int64(IPV6_RECVHOPOPTS), 53},
		{"IPV6_HOPOPTS", int64(IPV6_HOPOPTS), 54},
		{"IPV6_RTHDRDSTOPTS", int64(IPV6_RTHDRDSTOPTS), 55},
		{"IPV6_RECVRTHDR", int64(IPV6_RECVRTHDR), 56},
		{"IPV6_RTHDR", int64(IPV6_RTHDR), 57},
		{"IPV6_RECVDSTOPTS", int64(IPV6_RECVDSTOPTS), 58},
		{"IPV6_DSTOPTS", int64(IPV6_DSTOPTS), 59},
		{"IPV6_RECVPATHMTU", int64(IPV6_RECVPATHMTU), 60},
		{"IPV6_PATHMTU", int64(IPV6_PATHMTU), 61},
		{"IPV6_DONTFRAG", int64(IPV6_DONTFRAG), 62},
		{"IPV6_RECVTCLASS", int64(IPV6_RECVTCLASS), 66},
		{"IPV6_TCLASS", int64(IPV6_TCLASS), 67},
		{"IPV6_ADD_MEMBERSHIP", int64(IPV6_ADD_MEMBERSHIP), 20},
		{"IPV6_DROP_MEMBERSHIP", int64(IPV6_DROP_MEMBERSHIP), 21},
		{"IPV6_RXHOPOPTS", int64(IPV6_RXHOPOPTS), 54},
		{"IPV6_RXDSTOPTS", int64(IPV6_RXDSTOPTS), 59},
		{"IPV6_PMTUDISC_DONT", int64(IPV6_PMTUDISC_DONT), 0},
		{"IPV6_PMTUDISC_WANT", int64(IPV6_PMTUDISC_WANT), 1},
		{"IPV6_PMTUDISC_DO", int64(IPV6_PMTUDISC_DO), 2},
		{"IPV6_PMTUDISC_PROBE", int64(IPV6_PMTUDISC_PROBE), 3},
		{"IPV6_PMTUDISC_INTERFACE", int64(IPV6_PMTUDISC_INTERFACE), 4},
		{"IPV6_PMTUDISC_OMIT", int64(IPV6_PMTUDISC_OMIT), 5},
		{"SOL_IPV6", int64(SOL_IPV6), 41},
		{"SOL_ICMPV6", int64(SOL_ICMPV6), 58},
		{"IPV6_RTHDR_LOOSE", int64(IPV6_RTHDR_LOOSE), 0},
		{"IPV6_RTHDR_STRICT", int64(IPV6_RTHDR_STRICT), 1},
		{"IPV6_RTHDR_TYPE_0", int64(IPV6_RTHDR_TYPE_0), 0},
		{"X_BITS_SOCKADDR_H", int64(X_BITS_SOCKADDR_H), 1},
		{"X_SS_SIZE", int64(X_SS_SIZE), 128},
		{"PF_UNSPEC", int64(PF_UNSPEC), 0},
		{"PF_LOCAL", int64(PF_LOCAL), 1},
		{"PF_UNIX", int64(PF_UNIX), 1},
		{"PF_FILE", int64(PF_FILE), 1},
		{"PF_INET", int64(PF_INET), 2},
		{"PF_AX25", int64(PF_AX25), 3},
		{"PF_IPX", int64(PF_IPX), 4},
		{"PF_APPLETALK", int64(PF_APPLETALK), 5},
		{"PF_NETROM", int64(PF_NETROM), 6},
		{"PF_BRIDGE", int64(PF_BRIDGE), 7},
		{"PF_ATMPVC", int64(PF_ATMPVC), 8},
		{"PF_X25", int64(PF_X25), 9},
		{"PF_INET6", int64(PF_INET6), 10},
		{"PF_ROSE", int64(PF_ROSE), 11},
		{"PF_DECnet", int64(PF_DECnet), 12},
		{"PF_NETBEUI", int64(PF_NETBEUI), 13},
		{"PF_SECURITY", int64(PF_SECURITY), 14},
		{"PF_KEY", int64(PF_KEY), 15},
		{"PF_NETLINK", int64(PF_NETLINK), 16},
		{"PF_ROUTE", int64(PF_ROUTE), 16},
		{"PF_PACKET", int64(PF_PACKET), 17},
		{"PF_ASH", int64(PF_ASH), 18},
		{"PF_ECONET", int64(PF_ECONET), 19},
		{"PF_ATMSVC", int64(PF_ATMSVC), 20},
		{"PF_RDS", int64(PF_RDS), 21},
		{"PF_SNA", int64(PF_SNA), 22},
		{"PF_IRDA", int64(PF_IRDA), 23},
		{"PF_PPPOX", int64(PF_PPPOX), 24},
		{"PF_WANPIPE", int64(PF_WANPIPE), 25},
		{"PF_LLC", int64(PF_LLC), 26},
		{"PF_IB", int64(PF_IB), 27},
		{"PF_MPLS", int64(PF_MPLS), 28},
		{"PF_CAN", int64(PF_CAN), 29},
		{"PF_TIPC", int64(PF_TIPC), 30},
		{"PF_BLUETOOTH", int64(PF_BLUETOOTH), 31},
		{"PF_IUCV", int64(PF_IUCV), 32},
		{"PF_RXRPC", int64(PF_RXRPC), 33},
		{"PF_ISDN", int64(PF_ISDN), 34},
		{"PF_PHONET", int64(PF_PHONET), 35},
		{"PF_IEEE802154", int64(PF_IEEE802154), 36},
		{"PF_CAIF", int64(PF_CAIF), 37},
		{"PF_ALG", int64(PF_ALG), 38},
		{"PF_NFC", int64(PF_NFC), 39},
		{"PF_VSOCK", int64(PF_VSOCK), 40},
		{"PF_KCM", int64(PF_KCM), 41},
		{"PF_MAX", int64(PF_MAX), 42},
		{"AF_UNSPEC", int64(AF_UNSPEC), 0},
		{"AF_LOCAL", int64(AF_LOCAL), 1},
		{"AF_UNIX", int64(AF_UNIX), 1},
		{"AF_FILE", int64(AF_FILE), 1},
		{"AF_INET", int64(AF_INET), 2},
		{"AF_AX25", int64(AF_AX25), 3},
		{"AF_IPX", int64(AF_IPX), 4},
		{"AF_APPLETALK", int64(AF_APPLETALK), 5},
		{"AF_NETROM", int64(AF_NETROM), 6},
		{"AF_BRIDGE", int64(AF_BRIDGE), 7},
		{"AF_ATMPVC", int64(AF_ATMPVC), 8},
		{"AF_X25", int64(AF_X25), 9},
		{"AF_INET6", int64(AF_INET6), 10},
		{"AF_ROSE", int64(AF_ROSE), 11},
		{"AF_DECnet", int64(AF_DECnet), 12},
		{"AF_NETBEUI", int64(AF_NETBEUI), 13},
		{"AF_SECURITY", int64(AF_SECURITY), 14},
		{"AF_KEY", int64(AF_KEY), 15},
		{"AF_NETLINK", int64(AF_NETLINK), 16},
		{"AF_ROUTE", int64(AF_ROUTE), 16},
		{"AF_PACKET", int64(AF_PACKET), 17},
		{"AF_ASH", int64(AF_ASH), 18},
		{"AF_ECONET", int64(AF_ECONET), 19},
		{"AF_ATMSVC", int64(AF_ATMSVC), 20},
		{"AF_RDS", int64(AF_RDS), 21},
		{"AF_SNA", int64(AF_SNA), 22},
		{"AF_IRDA", int64(AF_IRDA), 23},
		{"AF_PPPOX", int64(AF_PPPOX), 24},
		{"AF_WANPIPE", int64(AF_WANPIPE), 25},
		{"AF_LLC", int64(AF_LLC), 26},
		{"AF_IB", int64(AF_IB), 27},
		{"AF_MPLS", int64(AF_MPLS), 28},
		{"AF_CAN", int64(AF_CAN), 29},
		{"AF_TIPC", int64(AF_TIPC), 30},
		{"AF_BLUETOOTH", int64(AF_BLUETOOTH), 31},
		{"AF_IUCV", int64(AF_IUCV), 32},
		{"AF_RXRPC", int64(AF_RXRPC), 33},
		{"AF_ISDN", int64(AF_ISDN), 34},
		{"AF_PHONET", int64(AF_PHONET), 35},
		{"AF_IEEE802154", int64(AF_IEEE802154), 36},
		{"AF_CAIF", int64(AF_CAIF), 37},
		{"AF_ALG", int64(AF_ALG), 38},
		{"AF_NFC", int64(AF_NFC), 39},
		{"AF_VSOCK", int64(AF_VSOCK), 40},
		{"AF_KCM", int64(AF_KCM), 41},
		{"AF_MAX", int64(AF_MAX), 42},
		{"SOL_RAW", int64(SOL_RAW), 255},
		{"SOL_DECNET", int64(SOL_DECNET), 261},
		{"SOL_X25", int64(SOL_X25), 262},
		{"SOL_PACKET", int64(SOL_PACKET), 263},
		{"SOL_ATM", int64(SOL_ATM), 264},
		{"SOL_AAL", int64(SOL_AAL), 265},
		{"SOL_IRDA", int64(SOL_IRDA), 266},
		{"SOL_NETBEUI", int64(SOL_NETBEUI), 267},
		{"SOL_LLC", int64(SOL_LLC), 268},
		{"SOL_DCCP", int64(SOL_DCCP), 269},
		{"SOL_NETLINK", int64(SOL_NETLINK), 270},
		{"SOL_TIPC", int64(SOL_TIPC), 271},
		{"SOL_RXRPC", int64(SOL_RXRPC), 272},
		{"SOL_PPPOL2TP", int64(SOL_PPPOL2TP), 273},
		{"SOL_BLUETOOTH", int64(SOL_BLUETOOTH), 274},
		{"SOL_PNPIPE", int64(SOL_PNPIPE), 275},
		{"SOL_RDS", int64(SOL_RDS), 276},
		{"SOL_IUCV", int64(SOL_IUCV), 277},
		{"SOL_CAIF", int64(SOL_CAIF), 278},
		{"SOL_ALG", int64(SOL_ALG), 279},
		{"SOL_NFC", int64(SOL_NFC), 280},
		{"SOL_KCM", int64(SOL_KCM), 281},
		{"SOMAXCONN", int64(SOMAXCONN), 128},
		{"MSG_OOB", int64(MSG_OOB), 1},
		{"MSG_PEEK", int64(MSG_PEEK), 2},
		{"MSG_DONTROUTE", int64(MSG_DONTROUTE), 4},
		{"MSG_CTRUNC", int64(MSG_CTRUNC), 8},
		{"MSG_PROXY", int64(MSG_PROXY), 16},
		{"MSG_TRUNC", int64(MSG_TRUNC), 32},
		{"MSG_DONTWAIT", int64(MSG_DONTWAIT), 64},
		{"MSG_EOR", int64(MSG_EOR), 128},
		{"MSG_WAITALL", int64(MSG_WAITALL), 256},
		{"MSG_FIN", int64(MSG_FIN), 512},
		{"MSG_SYN", int64(MSG_SYN), 1024},
		{"MSG_CONFIRM", int64(MSG_CONFIRM), 2048},
		{"MSG_RST", int64(MSG_RST), 4096},
		{"MSG_ERRQUEUE", int64(MSG_ERRQUEUE), 8192},
		{"MSG_NOSIGNAL", int64(MSG_NOSIGNAL), 16384},
		{"MSG_MORE", int64(MSG_MORE), 32768},
		{"MSG_WAITFORONE", int64(MSG_WAITFORONE), 65536},
		{"MSG_BATCH", int64(MSG_BATCH), 262144},
		{"MSG_FASTOPEN", int64(MSG_FASTOPEN), 536870912},
		{"MSG_CMSG_CLOEXEC", int64(MSG_CMSG_CLOEXEC), 1073741824},
		{"SCM_RIGHTS", int64(SCM_RIGHTS), 1},
		{"SOCK_DCCP", int64(SOCK_DCCP), 6},
		{"SOCK_PACKET", int64(SOCK_PACKET), 10},
		{"SOCK_CLOEXEC", int64(SOCK_CLOEXEC), 524288},
		{"SOCK_NONBLOCK", int64(SOCK_NONBLOCK), 2048},
		{"X_NET_IF_H", int64(X_NET_IF_H), 1},
		{"IF_NAMESIZE", int64(IF_NAMESIZE), 16},
		{"IFHWADDRLEN", int64(IFHWADDRLEN), 6},
		{"IFNAMSIZ", int64(IFNAMSIZ), 16},
		{"IFF_UP", int64(IFF_UP), 1},
		{"IFF_BROADCAST", int64(IFF_BROADCAST), 2},
		{"IFF_DEBUG", int64(IFF_DEBUG), 4},
		{"IFF_LOOPBACK", int64(IFF_LOOPBACK), 8},
		{"IFF_POINTOPOINT", int64(IFF_POINTOPOINT), 16},
		{"IFF_NOTRAILERS", int64(IFF_NOTRAILERS), 32},
		{"IFF_RUNNING", int64(IFF_RUNNING), 64},
		{"IFF_NOARP", int64(IFF_NOARP), 128},
		{"IFF_PROMISC", int64(IFF_PROMISC), 256},
		{"IFF_ALLMULTI", int64(IFF_ALLMULTI), 512},
		{"IFF_MASTER", int64(IFF_MASTER), 1024},
		{"IFF_SLAVE", int64(IFF_SLAVE), 2048},
		{"IFF_MULTICAST", int64(IFF_MULTICAST), 4096},
		{"IFF_PORTSEL", int64(IFF_PORTSEL), 8192},
		{"IFF_AUTOMEDIA", int64(IFF_AUTOMEDIA), 16384},
		{"IFF_DYNAMIC", int64(IFF_DYNAMIC), 32768},
		{"X_NETDB_H", int64(X_NETDB_H), 1},
		{"HOST_NOT_FOUND", int64(HOST_NOT_FOUND), 1},
		{"TRY_AGAIN", int64(TRY_AGAIN), 2},
		{"NO_RECOVERY", int64(NO_RECOVERY), 3},
		{"NO_DATA", int64(NO_DATA), 4},
		{"NETDB_INTERNAL", int64(NETDB_INTERNAL), -1},
		{"NETDB_SUCCESS", int64(NETDB_SUCCESS), 0},
		{"NO_ADDRESS", int64(NO_ADDRESS), 4},
		{"AI_PASSIVE", int64(AI_PASSIVE), 1},
		{"AI_CANONNAME", int64(AI_CANONNAME), 2},
		{"AI_NUMERICHOST", int64(AI_NUMERICHOST), 4},
		{"AI_V4MAPPED", int64(AI_V4MAPPED), 8},
		{"AI_ALL", int64(AI_ALL), 16},
		{"AI_ADDRCONFIG", int64(AI_ADDRCONFIG), 32},
		{"AI_NUMERICSERV", int64(AI_NUMERICSERV), 1024},
		{"EAI_BADFLAGS", int64(EAI_BADFLAGS), -1},
		{"EAI_NONAME", int64(EAI_NONAME), -2},
		{"EAI_AGAIN", int64(EAI_AGAIN), -3},
		{"EAI_FAIL", int64(EAI_FAIL), -4},
		{"EAI_FAMILY", int64(EAI_FAMILY), -6},
		{"EAI_SOCKTYPE", int64(EAI_SOCKTYPE), -7},
		{"EAI_SERVICE", int64(EAI_SERVICE), -8},
		{"EAI_MEMORY", int64(EAI_MEMORY), -10},
		{"EAI_SYSTEM", int64(EAI_SYSTEM), -11},
		{"EAI_OVERFLOW", int64(EAI_OVERFLOW), -12},
		{"NI_MAXHOST", int64(NI_MAXHOST), 1025},
		{"NI_MAXSERV", int64(NI_MAXSERV), 32},
		{"NI_NUMERICHOST", int64(NI_NUMERICHOST), 1},
		{"NI_NUMERICSERV", int64(NI_NUMERICSERV), 2},
		{"NI_NOFQDN", int64(NI_NOFQDN), 4},
		{"NI_NAMEREQD", int64(NI_NAMEREQD), 8},
		{"NI_DGRAM", int64(NI_DGRAM), 16},
		{"X_NETINET_IN_H", int64(X_NETINET_IN_H), 1},
		{"IN_CLASSA_NET", int64(IN_CLASSA_NET), 4278190080},
		{"IN_CLASSA_NSHIFT", int64(IN_CLASSA_NSHIFT), 24},
		{"IN_CLASSA_HOST", int64(IN_CLASSA_HOST), 16777215},
		{"IN_CLASSA_MAX", int64(IN_CLASSA_MAX), 128},
		{"IN_CLASSB_NET", int64(IN_CLASSB_NET), 4294901760},
		{"IN_CLASSB_NSHIFT", int64(IN_CLASSB_NSHIFT), 16},
		{"IN_CLASSB_HOST", int64(IN_CLASSB_HOST), 65535},
		{"IN_CLASSB_MAX", int64(IN_CLASSB_MAX), 65536},
		{"IN_CLASSC_NET", int64(IN_CLASSC_NET), 4294967040},
		{"IN_CLASSC_NSHIFT", int64(IN_CLASSC_NSHIFT), 8},
		{"IN_CLASSC_HOST", int64(IN_CLASSC_HOST), 255},
		{"IN_LOOPBACKNET", int64(IN_LOOPBACKNET), 127},
		{"INET_ADDRSTRLEN", int64(INET_ADDRSTRLEN), 16},
		{"INET6_ADDRSTRLEN", int64(INET6_ADDRSTRLEN), 46},
		{"IPPROTO_IP", int64(IPPROTO_IP), 0},
		{"IPPROTO_ICMP", int64(IPPROTO_ICMP), 1},
		{"IPPROTO_IGMP", int64(IPPROTO_IGMP), 2},
		{"IPPROTO_IPIP", int64(IPPROTO_IPIP), 4},
		{"IPPROTO_TCP", int64(IPPROTO_TCP), 6},
		{"IPPROTO_EGP", int64(IPPROTO_EGP), 8},
		{"IPPROTO_PUP", int64(IPPROTO_PUP), 12},
		{"IPPROTO_UDP", int64(IPPROTO_UDP), 17},
		{"IPPROTO_IDP", int64(IPPROTO_IDP), 22},
		{"IPPROTO_TP", int64(IPPROTO_TP), 29},
		{"IPPROTO_DCCP", int64(IPPROTO_DCCP), 33},
		{"IPPROTO_IPV6", int64(IPPROTO_IPV6), 41},
		{"IPPROTO_RSVP", int64(IPPROTO_RSVP), 46},
		{"IPPROTO_GRE", int64(IPPROTO_GRE), 47},
		{"IPPROTO_ESP", int64(IPPROTO_ESP), 50},
		{"IPPROTO_AH", int64(IPPROTO_AH), 51},
		{"IPPROTO_MTP", int64(IPPROTO_MTP), 92},
		{"IPPROTO_BEETPH", int64(IPPROTO_BEETPH), 94},
		{"IPPROTO_ENCAP", int64(IPPROTO_ENCAP), 98},
		{"IPPROTO_PIM", int64(IPPROTO_PIM), 103},
		{"IPPROTO_COMP", int64(IPPROTO_COMP), 108},
		{"IPPROTO_SCTP", int64(IPPROTO_SCTP), 132},
		{"IPPROTO_UDPLITE", int64(IPPROTO_UDPLITE), 136},
		{"IPPROTO_MPLS", int64(IPPROTO_MPLS), 137},
		{"IPPROTO_RAW", int64(IPPROTO_RAW), 255},
		{"IPPROTO_MAX", int64(IPPROTO_MAX), 256},
		{"IPPROTO_HOPOPTS", int64(IPPROTO_HOPOPTS), 0},
		{"IPPROTO_ROUTING", int64(IPPROTO_ROUTING), 43},
		{"IPPROTO_FRAGMENT", int64(IPPROTO_FRAGMENT), 44},
		{"IPPROTO_ICMPV6", int64(IPPROTO_ICMPV6), 58},
		{"IPPROTO_NONE", int64(IPPROTO_NONE), 59},
		{"IPPROTO_DSTOPTS", int64(IPPROTO_DSTOPTS), 60},
		{"IPPROTO_MH", int64(IPPROTO_MH), 135},
		{"IPPORT_ECHO", int64(IPPORT_ECHO), 7},
		{"IPPORT_DISCARD", int64(IPPORT_DISCARD), 9},
		{"IPPORT_SYSTAT", int64(IPPORT_SYSTAT), 11},
		{"IPPORT_DAYTIME", int64(IPPORT_DAYTIME), 13},
		{"IPPORT_NETSTAT", int64(IPPORT_NETSTAT), 15},
		{"IPPORT_FTP", int64(IPPORT_FTP), 21},
		{"IPPORT_TELNET", int64(IPPORT_TELNET), 23},
		{"IPPORT_SMTP", int64(IPPORT_SMTP), 25},
		{"IPPORT_TIMESERVER", int64(IPPORT_TIMESERVER), 37},
		{"IPPORT_NAMESERVER", int64(IPPORT_NAMESERVER), 42},
		{"IPPORT_WHOIS", int64(IPPORT_WHOIS), 43},
		{"IPPORT_MTP", int64(IPPORT_MTP), 57},
		{"IPPORT_TFTP", int64(IPPORT_TFTP), 69},
		{"IPPORT_RJE", int64(IPPORT_RJE), 77},
		{"IPPORT_FINGER", int64(IPPORT_FINGER), 79},
		{"IPPORT_TTYLINK", int64(IPPORT_TTYLINK), 87},
		{"IPPORT_SUPDUP", int64(IPPORT_SUPDUP), 95},
		{"IPPORT_EXECSERVER", int64(IPPORT_EXECSERVER), 512},
		{"IPPORT_LOGINSERVER", int64(IPPORT_LOGINSERVER), 513},
		{"IPPORT_CMDSERVER", int64(IPPORT_CMDSERVER), 514},
		{"IPPORT_EFSSERVER", int64(IPPORT_EFSSERVER), 520},
		{"IPPORT_BIFFUDP", int64(IPPORT_BIFFUDP), 512},
		{"IPPORT_WHOSERVER", int64(IPPORT_WHOSERVER), 513},
		{"IPPORT_ROUTESERVER", int64(IPPORT_ROUTESERVER), 520},
		{"IPPORT_RESERVED", int64(IPPORT_RESERVED), 1024},
		{"IPPORT_USERRESERVED", int64(IPPORT_USERRESERVED), 5000},
		{"X_RPC_NETDB_H", int64(X_RPC_NETDB_H), 1},
		{"X_SYS_SOCKET_H", int64(X_SYS_SOCKET_H), 1},
		{"SHUT_RD", int64(SHUT_RD), 0},
		{"SHUT_WR", int64(SHUT_WR), 1},
		{"SHUT_RDWR", int64(SHUT_RDWR), 2},
		{"X_SYS_UN_H", int64(X_SYS_UN_H), 1},
		{"ALIGNSIZE", int64(ALIGNSIZE), 8},
		{"MAXSIZE", int64(MAXSIZE), 128},
		{"PAD1_SIZE", int64(PAD1_SIZE), 118},
		{"PAD2_SIZE", int64(PAD2_SIZE), 0},
		{"SOCK_STREAM", int64(SOCK_STREAM), 1},
		{"SOCK_DGRAM", int64(SOCK_DGRAM), 2},
		{"SOCK_RAW", int64(SOCK_RAW), 3},
		{"SOCK_RDM", int64(SOCK_RDM), 4},
		{"SOCK_SEQPACKET", int64(SOCK_SEQPACKET), 5},
	}
	var networkStringConstants = []struct{ name, got, want string }{
		{"X_PATH_HEQUIV", X_PATH_HEQUIV, "/etc/hosts.equiv"},
		{"X_PATH_HOSTS", X_PATH_HOSTS, "/etc/hosts"},
		{"X_PATH_NETWORKS", X_PATH_NETWORKS, "/etc/networks"},
		{"X_PATH_NSSWITCH_CONF", X_PATH_NSSWITCH_CONF, "/etc/nsswitch.conf"},
		{"X_PATH_PROTOCOLS", X_PATH_PROTOCOLS, "/etc/protocols"},
		{"X_PATH_SERVICES", X_PATH_SERVICES, "/etc/services"},
	}
	var networkVariableSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{
		{"In6addrAny", unsafe.Pointer(&In6addrAny)},
		{"In6addrLoopback", unsafe.Pointer(&In6addrLoopback)},
	}
	var networkAccessorChecks = []struct {
		name                string
		check               func() (uintptr, uintptr, uintptr)
		offset, size, align uintptr
	}{
		{"*_llcppg_anon_1.XGof_ref_ifu_broadaddr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_1
			view := value.XGof_ref_ifu_broadaddr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_1.XGof_ref_ifu_dstaddr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_1
			view := value.XGof_ref_ifu_dstaddr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_2.XGof_ref_ifrn_name", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_2
			view := value.XGof_ref_ifrn_name()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 1},
		{"*_llcppg_anon_3.XGof_ref_ifru_addr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_addr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_dstaddr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_dstaddr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_broadaddr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_broadaddr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_netmask", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_netmask()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_hwaddr", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_hwaddr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_flags", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_flags()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 2, 2},
		{"*_llcppg_anon_3.XGof_ref_ifru_ivalue", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_ivalue()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 4, 4},
		{"*_llcppg_anon_3.XGof_ref_ifru_mtu", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_mtu()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 4, 4},
		{"*_llcppg_anon_3.XGof_ref_ifru_map", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_map()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 24, 8},
		{"*_llcppg_anon_3.XGof_ref_ifru_slave", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_slave()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 1},
		{"*_llcppg_anon_3.XGof_ref_ifru_newname", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_newname()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 1},
		{"*_llcppg_anon_3.XGof_ref_ifru_data", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_3
			view := value.XGof_ref_ifru_data()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*_llcppg_anon_4.XGof_ref_ifcu_buf", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_4
			view := value.XGof_ref_ifcu_buf()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*_llcppg_anon_4.XGof_ref_ifcu_req", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_4
			view := value.XGof_ref_ifcu_req()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
	}

	t.Run("function symbols", func(t *testing.T) {
		for _, symbol := range networkFunctionSymbols {
			t.Run(symbol.name, func(t *testing.T) {
				if symbol.address == nil {
					t.Fatal("C function address is nil")
				}
				t.Logf("C function address: %p", symbol.address)
			})
		}
	})
	t.Run("native prototypes", func(t *testing.T) {
		check := func(name string, got reflect.Type, want networkABISlot) {
			if got.Size() != want.size || uintptr(got.Align()) != want.align {
				t.Errorf("%s size/alignment = %d/%d, C = %d/%d", name, got.Size(), got.Align(), want.size, want.align)
			}
			class := 0
			switch got.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				class = 1
			case reflect.Ptr, reflect.UnsafePointer, reflect.Func:
				class = 5
			case reflect.Float32, reflect.Float64:
				class = 8
			case reflect.Struct:
				class = 12
			}
			if class != want.class {
				t.Errorf("%s ABI class = %d (%v), C = %d", name, class, got.Kind(), want.class)
			}
		}
		for _, tc := range networkFunctionPrototypes {
			t.Run(tc.name, func(t *testing.T) {
				inputs := len(tc.inputs)
				if tc.variadic {
					inputs++
				}
				if tc.signature.Kind() != reflect.Func || tc.signature.NumIn() != inputs || tc.signature.IsVariadic() != tc.variadic {
					t.Fatalf("Go signature = %v, C fixed parameters = %d, variadic = %v", tc.signature, len(tc.inputs), tc.variadic)
				}
				for i, slot := range tc.inputs {
					check(tc.name+" parameter", tc.signature.In(i), slot)
				}
				outputs := 0
				if tc.output != nil {
					outputs = 1
				}
				if tc.signature.NumOut() != outputs {
					t.Fatalf("Go result count = %d, C = %d", tc.signature.NumOut(), outputs)
				}
				if tc.output != nil {
					check(tc.name+" result", tc.signature.Out(0), *tc.output)
				}
			})
		}
	})
	t.Run("native types and globals", func(t *testing.T) {
		for _, tc := range networkTypeLayouts {
			t.Run(tc.symbol+"/"+tc.property, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go layout = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
		for _, symbol := range networkVariableSymbols {
			if symbol.address == nil {
				t.Errorf("%s C variable address is nil", symbol.name)
			}
			t.Logf("%s C variable address: %p", symbol.name, symbol.address)
		}
	})
	t.Run("native constants", func(t *testing.T) {
		for _, tc := range networkIntegerConstants {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go constant = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
		for _, tc := range networkStringConstants {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go constant = %q, C headers = %q", tc.got, tc.want)
				}
			})
		}
	})
	t.Run("union accessors", func(t *testing.T) {
		for _, tc := range networkAccessorChecks {
			t.Run(tc.name, func(t *testing.T) {
				offset, size, align := tc.check()
				if offset != tc.offset || size != tc.size || align != tc.align {
					t.Fatalf("Go field offset/size/alignment = %d/%d/%d, C = %d/%d/%d", offset, size, align, tc.offset, tc.size, tc.align)
				}
			})
		}
	})
	t.Run("symbol inventory", func(t *testing.T) {
		functions, prototypes, types, constants, variables := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
		for _, tc := range networkFunctionPrototypes {
			prototypes[tc.name] = true
		}
		for _, symbol := range networkFunctionSymbols {
			functions[symbol.name] = true
			if !prototypes[symbol.name] {
				t.Errorf("%s has no native prototype test", symbol.name)
			}
		}
		for _, tc := range networkAccessorChecks {
			functions[tc.name] = true
		}
		for _, tc := range networkTypeLayouts {
			types[tc.symbol] = true
		}
		for _, tc := range networkIntegerConstants {
			constants[tc.name] = true
		}
		for _, tc := range networkStringConstants {
			constants[tc.name] = true
		}
		for _, symbol := range networkVariableSymbols {
			variables[symbol.name] = true
		}
		constants["XGoPackage"], constants["LLGoPackage"] = true, true
		wantLink := "link: -lc -lresolv"
		if runtime.GOOS == "darwin" {
			wantLink = "link: -lSystem"
		}
		if !XGoPackage || LLGoPackage != wantLink {
			t.Errorf("package metadata = %v/%q, want true/%q", XGoPackage, LLGoPackage, wantLink)
		}
		files := []string{"libc_" + runtime.GOOS + "_" + runtime.GOARCH + ".go", "legacy_native.go"}
		if runtime.GOOS == "darwin" {
			files = append(files, "legacy_darwin_arm64.go")
		}
		for _, filename := range files {
			source, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			for _, declaration := range source.Decls {
				switch decl := declaration.(type) {
				case *ast.FuncDecl:
					name := decl.Name.Name
					if !ast.IsExported(name) {
						continue
					}
					if decl.Recv != nil {
						receiver, prefix := decl.Recv.List[0].Type, ""
						if pointer, ok := receiver.(*ast.StarExpr); ok {
							receiver, prefix = pointer.X, "*"
						}
						name = prefix + receiver.(*ast.Ident).Name + "." + name
					}
					if !functions[name] {
						t.Errorf("%s: function %s has no test", filename, name)
					}
					if strings.HasPrefix(decl.Name.Name, "Native") {
						t.Errorf("%s retains duplicate %s", filename, name)
					}
					if decl.Doc != nil {
						for _, comment := range decl.Doc.List {
							for _, private := range []string{"C.pfctlinput", "C.__h_errno_location", "C.__cmsg_nxthdr"} {
								if strings.Contains(comment.Text, private) {
									t.Errorf("%s still exposes private C entry point %s", filename, private)
								}
							}
						}
					}
				case *ast.GenDecl:
					for _, specification := range decl.Specs {
						switch spec := specification.(type) {
						case *ast.TypeSpec:
							if ast.IsExported(spec.Name.Name) && !types[spec.Name.Name] {
								t.Errorf("%s: type %s has no ABI test", filename, spec.Name.Name)
							}
						case *ast.ValueSpec:
							for _, name := range spec.Names {
								if !ast.IsExported(name.Name) {
									continue
								}
								if decl.Tok == token.CONST && !constants[name.Name] {
									t.Errorf("%s: constant %s has no test", filename, name.Name)
								}
								if decl.Tok == token.VAR && !variables[name.Name] {
									t.Errorf("%s: variable %s has no test", filename, name.Name)
								}
							}
						}
					}
				}
			}
		}
		t.Logf("coverage: %d C functions, %d Go union accessors, %d types/global layouts, %d constants, %d globals", len(networkFunctionSymbols), len(networkAccessorChecks), len(networkTypeLayouts), len(constants), len(networkVariableSymbols))
	})

}
