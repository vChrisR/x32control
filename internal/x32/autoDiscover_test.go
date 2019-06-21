package x32_test

import (
	"fmt"
	"testing"

	"github.com/vchrisr/x32control/internal/x32"
)

func TestAutoDiscover(t *testing.T) {
	ip, err := x32.AutoDiscover(3)
	fmt.Println(ip)
	fmt.Println(err)
}
