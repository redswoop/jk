package main

/*
#include <libproc.h>
#include <sys/proc_info.h>
*/
import "C"
import "unsafe"

// getCwd returns the working directory of a process using libproc.
// Returns "" if the pid is invalid or the cwd cannot be determined.
func getCwd(pid int) string {
	var info C.struct_proc_vnodepathinfo
	size := C.int(unsafe.Sizeof(info))

	ret := C.proc_pidinfo(C.int(pid), C.PROC_PIDVNODEPATHINFO, 0, unsafe.Pointer(&info), size)
	if ret <= 0 {
		return ""
	}

	return C.GoString(&info.pvi_cdir.vip_path[0])
}
