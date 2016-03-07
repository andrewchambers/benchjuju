// mksyscall_windows.pl -l32 symlink/symlink_windows.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package symlink

import "unsafe"
import "syscall"

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCreateSymbolicLinkW       = modkernel32.NewProc("CreateSymbolicLinkW")
	procGetFinalPathNameByHandleW = modkernel32.NewProc("GetFinalPathNameByHandleW")
)

func createSymbolicLink(symlinkname *uint16, targetname *uint16, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procCreateSymbolicLinkW.Addr(), 3, uintptr(unsafe.Pointer(symlinkname)), uintptr(unsafe.Pointer(targetname)), uintptr(flags))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getFinalPathNameByHandle(handle syscall.Handle, buf *uint16, buflen uint32, flags uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetFinalPathNameByHandleW.Addr(), 4, uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(buflen), uintptr(flags), 0, 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
