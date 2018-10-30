package gomsr

import (
	"fmt"
	"testing"
)

func Test_ReadMSR(t *testing.T) {

	fd, err := ReadMSRWithLocation(0, 0, "/tmp/msr_test%d.txt")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	fmt.Printf("Got %d\n", fd)
}

// func Test_ReadMSRReal(t *testing.T) {
// 	fd, err := ReadMSR(0, 0x198)
// 	if err != nil {
// 		t.Fatalf("Error: %s", err)
// 	}

// 	fmt.Printf("Got 0x%x\n", fd)
// }
