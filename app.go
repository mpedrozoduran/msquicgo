package main

/*
#cgo CFLAGS: -I/home/miguel/samples/msquic/src/inc
#cgo LDFLAGS: -L/home/miguel/samples/msquic/build/bin/Release -l:libmsquic.so
#include "msquic.h"
#include "stdlib.h"
const QUIC_REGISTRATION_CONFIG RegConfig = { "quicsample", QUIC_EXECUTION_PROFILE_LOW_LATENCY };
HQUIC Registration;
const QUIC_API_TABLE* MsQuic;

int call_msquic_reg_open(HQUIC* reg, QUIC_REGISTRATION_CONFIG* config) {
	return MsQuic->RegistrationOpen(config, reg);
}
*/
import "C"
import (
	"log"
	"unsafe"
)

type Quic struct {
	APITable  *C.struct_QUIC_API_TABLE
	RegConfig C.struct_QUIC_REGISTRATION_CONFIG
	Config    *C.HQUIC
}

func NewQuic() Quic {
	quic := Quic{RegConfig: C.RegConfig, Config: &C.Registration}
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

}

func (m *Quic) runClient() {

}

func (m *Quic) free() {
	C.MsQuicClose(m.APITable)
}
