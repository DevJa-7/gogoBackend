package twillo

import (
	"fmt"
	"testing"
)

func TestSendVerifySMS(t *testing.T) {
	err := SendVerifySMS("+50671735146", "123456")
	if err != nil {
		fmt.Println(err)
	}
}
