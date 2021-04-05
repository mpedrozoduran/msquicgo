package main

/*
#cgo CFLAGS: -I/home/miguel/samples/msquic/src/inc
#cgo LDFLAGS: -L/home/miguel/samples/msquic/build/bin/Release -l:libmsquic.so
#include "msquic.h"
#include "msquic_posix.h"
#include "stdlib.h"
const QUIC_REGISTRATION_CONFIG RegConfig = { "quicsample", QUIC_EXECUTION_PROFILE_LOW_LATENCY };
HQUIC Registration;
const QUIC_API_TABLE* MsQuic;

static int call_msquic_reg_open(HQUIC* reg, QUIC_REGISTRATION_CONFIG* config) {
	return MsQuic->RegistrationOpen(config, reg);
}

static int call_msquic_listeneropen(HQUIC reg, void* callback, HQUIC Listener) {
	return MsQuic->ListenerOpen(reg, callback, ((void*)0), &Listener);
}

static int call_msquic_listenerstart(QUIC_ADDR* addr, HQUIC Listener) {
	const QUIC_BUFFER Alpn = { sizeof("sample") - 1, (uint8_t*)"sample" };
	return MsQuic->ListenerStart(Listener, &Alpn, 1, &addr);
}
*/
import "C"
import (
	"fmt"
	"github.com/mattn/go-pointer"
	"log"
	"unsafe"
)

const (
	QuicAddressFamilyUnspec       = 0
	QuicUdpPort             int16 = 4567
)

type Quic struct {
	APITable  *C.struct_QUIC_API_TABLE
	RegConfig C.struct_QUIC_REGISTRATION_CONFIG
	Config    *C.HQUIC
	Address   *C.QUIC_ADDR
}

//export ServerListenerCallback
func ServerListenerCallback(hquic *C.HQUIC, voidP unsafe.Pointer, listener *C.QUIC_LISTENER_EVENT) int {
	fmt.Println("Call ServerListenerCallback")
	return 0
}

//export ServerConnectionCallback
func ServerConnectionCallback(hquic *C.HQUIC, voidP unsafe.Pointer, event *C.QUIC_CONNECTION_EVENT) int {
	fmt.Println("Call ServerConnectionCallback")
	return 0
}

func NewQuic() Quic {
	quic := Quic{RegConfig: C.RegConfig, Config: &C.Registration}
	C.QuicAddrSetFamily(pointer.Save(&quic.Address), QuicAddressFamilyUnspec)
	C.QuicAddrSetPort(pointer.Save(&quic.Address), QuicUdpPort)
	quic.open()
	quic.registration()
	return quic
}

func main() {
	quic := NewQuic()
	quic.runServer()
	quic.runClient()
}

func (m *Quic) open() {
	status := C.MsQuicOpen(pointer.Save(&m.APITable))
	if status > 0 {
		log.Fatalf("MsQuicOpen failed, status is: %v", status)
	}
}

func (m *Quic) registration() {
	cs := C.CString("myquicsample")
	C.call_msquic_reg_open(pointer.Save(m.Config), pointer.Save(&m.RegConfig))
	C.free(unsafe.Pointer(cs))
}

func (m *Quic) runServer() {
	status := C.call_msquic_listeneropen(pointer.Save(&m.RegConfig), ServerListenerCallback, nil)
	if status > 0 {
		log.Fatalf("MsQuic->ListenerOpen failed, status is: %v", status)
	}
	status = C.call_msquic_listenerstart(pointer.Save(m.Address), nil)
	if status > 0 {
		log.Fatalf("MsQuic->ListenerStart failed, status is: %v", status)
	}
}

func (m *Quic) runClient() {

}

func (m *Quic) free() {
	C.MsQuicClose(pointer.Save(m.APITable))
}
