package gomonkey

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

func PtrOf(val []byte) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&val)).Data
}

func modifyBinary(target uintptr, bytes []byte) {
	targetPage := pageStart(target)
	endPage := pageStart(target + uintptr(len(bytes)) - 1)
	protectSize := syscall.Getpagesize()
	if endPage > targetPage {
		protectSize = int(endPage-targetPage) + syscall.Getpagesize()
	}
	res := write(target, PtrOf(bytes), len(bytes), targetPage, protectSize, syscall.PROT_READ|syscall.PROT_EXEC)
	if res != 0 {
		panic(fmt.Errorf("failed to write memory, code %v", res))
	}
}

//go:cgo_import_dynamic mach_task_self mach_task_self "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic mach_vm_protect mach_vm_protect "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic sys_icache_invalidate sys_icache_invalidate "/usr/lib/libSystem.B.dylib"
func write(target, data uintptr, len int, page uintptr, pageSize, oriProt int) int
