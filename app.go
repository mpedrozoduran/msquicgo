package main

/*
#cgo CFLAGS: -I/home/miguel/samples/msquic/src/inc
#cgo LDFLAGS: -L/home/miguel/samples/msquic/build/bin/Release -l:libmsquic.so
#include "msquic.h"
#include "msquic_posix.h"
#include "stdlib.h"
const QUIC_REGISTRATION_CONFIG RegConfig = { "quicsample", QUIC_EXECUTION_PROFILE_LOW_LATENCY };
HQUIC Registration;
HQUIC Listener = nullptr;
const QUIC_API_TABLE* MsQuic;

const QUIC_BUFFER Alpn = { sizeof("sample") - 1, (uint8_t*)"sample" };

_IRQL_requires_max_(PASSIVE_LEVEL)
_Function_class_(QUIC_LISTENER_CALLBACK)
QUIC_STATUS
QUIC_API
ServerListenerCallback(
    _In_ HQUIC,
_In_opt_ void*,
_Inout_ QUIC_LISTENER_EVENT* Event
)
{
QUIC_STATUS Status = QUIC_STATUS_NOT_SUPPORTED;
switch (Event->Type) {
case QUIC_LISTENER_EVENT_NEW_CONNECTION:
MsQuic->SetCallbackHandler(Event->NEW_CONNECTION.Connection, (void*)ServerConnectionCallback, nullptr);
Status = MsQuic->ConnectionSetConfiguration(Event->NEW_CONNECTION.Connection, Configuration);
break;
default:
break;
}
return Status;
}

int call_msquic_reg_open(HQUIC* reg, QUIC_REGISTRATION_CONFIG* config) {
	return MsQuic->RegistrationOpen(config, reg);
}

int call_msquic_listeneropen(HQUIC* reg) {
	return MsQuic->ListenerOpen(reg, ServerListenerCallback, nullptr, &Listener)
}

int call_msquic_listenerstart(QUIC_ADDR* addr) {
	return MsQuic->ListenerStart(Listener, &Alpn, 1, &Address)
}
*/
import "C"
import (
	"log"
	"unsafe"
)

const (
	QuicAddressFamilyUnspec = 0
)

type Quic struct {
	APITable  *C.struct_QUIC_API_TABLE
	RegConfig C.struct_QUIC_REGISTRATION_CONFIG
	Config    *C.HQUIC
	Address   *C.QUIC_ADDR
}

func NewQuic() Quic {
	quic := Quic{RegConfig: C.RegConfig, Config: &C.Registration}
	C.QuicAddrSetFamily(&quic.Address, QuicAddressFamilyUnspec)
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
	status := C.MsQuicOpen(&m.APITable)
	if status > 0 {
		log.Fatalf("MsQuicOpen failed, status is: %v", status)
	}
}

func (m *Quic) registration() {
	cs := C.CString("myquicsample")
	C.call_msquic_reg_open(m.Config, &m.RegConfig)
	C.free(unsafe.Pointer(cs))
}

func (m *Quic) runServer() {
	C.call_msquic_listeneropen(&m.RegConfig)
	C.call_msquic_listenerstart(m.Address)
}

func (m *Quic) runClient() {

}

func (m *Quic) free() {
	C.MsQuicClose(m.APITable)
}
