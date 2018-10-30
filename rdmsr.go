package gomsr

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"syscall"
)

const defaultFmtStr = "/dev/cpu/%d/msr"

func openMSR(cpu int, fmtString string) (int, error) {

	cpuDir := fmt.Sprintf(fmtString, cpu)
	return syscall.Open(cpuDir, 0, 777)
}

func readMSR(cpu int, msr int64, fmtStr string) (uint64, error) {
	fd, err := openMSR(cpu, fmtStr)
	if err != nil {
		return 0, err
	}

	regBuf := make([]byte, 8)

	rc, err := syscall.Pread(fd, regBuf, msr)

	if err != nil {
		return 0, err
	}

	if rc != 8 {
		return 0, fmt.Errorf("Read wrong count of bytes: %d", rc)
	}
	fmt.Println(hex.Dump(regBuf))

	//I'm gonna go ahead and assume an x86 processor will be little endian
	msrOut := binary.LittleEndian.Uint64(regBuf)

	return msrOut, nil
}

//MSRDev represents a handler for frequent read/write operations
//for one-off MSR read/writes, gomsr provides {Read,Write}MSR*() functions
type MSRDev struct {
	fd int
}

//MSR provides an interface for reoccouring access to a given CPU's MSR interface
func MSR(cpu int) (MSRDev, error) {
	cpuDir := fmt.Sprintf(defaultFmtStr, cpu)
	fd, err := syscall.Open(cpuDir, 0, 777)
	if err != nil {
		return MSRDev{}, err
	}
	return MSRDev{fd: fd}, nil
}

//MSRWithLocation is the same as MSR(), but takes a custom location, for use with testing or 3rd party utilities like llnl/msr-safe
//It takes a string that has a `%d` format specifier for the cpu. For example: /dev/cpu/%d/msr_safe
func MSRWithLocation(cpu int, fmtString string) (MSRDev, error) {
	cpuDir := fmt.Sprintf(fmtString, cpu)
	fd, err := syscall.Open(cpuDir, 0, 777)
	if err != nil {
		return MSRDev{}, err
	}
	return MSRDev{fd: fd}, nil
}

//Read reads a given MSR on the CPU and returns the uint64
func (d MSRDev) Read(msr int64) (uint64, error) {
	regBuf := make([]byte, 8)

	rc, err := syscall.Pread(d.fd, regBuf, msr)

	if err != nil {
		return 0, err
	}

	if rc != 8 {
		return 0, fmt.Errorf("Read wrong count of bytes: %d", rc)
	}

	//I'm gonna go ahead and assume an x86 processor will be little endian
	msrOut := binary.LittleEndian.Uint64(regBuf)

	return msrOut, nil
}

//Close closes the connection to the MSR
func (d MSRDev) Close() error {
	return syscall.Close(d.fd)
}

//ReadMSRWithLocation is like ReadMSR(), but takes a custom location, for use with testing or 3rd party utilities like llnl/msr-safe
//It takes a string that has a `%d` format specifier for the cpu. For example: /dev/cpu/%d/msr_safe
func ReadMSRWithLocation(cpu int, msr int64, fmtStr string) (uint64, error) {

	m, err := MSRWithLocation(cpu, fmtStr)
	if err != nil {
		return 0, err
	}

	msrD, err := m.Read(msr)
	if err != nil {
		return 0, err
	}

	return msrD, m.Close()

}

//ReadMSR reads the MSR on the given CPU as a one-time operation
func ReadMSR(cpu int, msr int64) (uint64, error) {
	m, err := MSR(cpu)
	if err != nil {
		return 0, err
	}

	msrD, err := m.Read(msr)
	if err != nil {
		return 0, err
	}

	return msrD, m.Close()

}
