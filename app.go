package main

/*
#cgo CFLAGS: -I/home/miguel/samples/msquic/src/inc
#cgo LDFLAGS: -L/home/miguel/samples/msquic/build/bin/Release -l:libmsquic.so
#include "msquic.h"
*/
import "C"
import "log"

type Quic struct {
	APITable  *C.struct_QUIC_API_TABLE
	RegConfig C.struct_QUIC_REGISTRATION_CONFIG
	Config    C.struct_QUIC_HANDLE
}

func NewQuic() Quic {
	quic := Quic{}
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
	regConfig := C.struct_QUIC_REGISTRATION_CONFIG{"myquicsample", C.QUIC_EXECUTION_PROFILE_LOW_LATENCY}
	config := C.struct_QUIC_HANDLE
	m.Config = config
	m.RegConfig = regConfig
	m.APITable.RegistrationOpen(&regConfig, &config)
}

func (m *Quic) runServer() {

}

func (m *Quic) runClient() {

}

func (m *Quic) free() {
	C.MSQuicClose(&m.APITable)
}
