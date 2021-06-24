package timeHelper

import (
	"fmt"
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	fmt.Println(time.Now().Unix())

	fmt.Println(time.Date(1990, 1, 1, 8, 0, 0, 0, time.Now().Location()).Unix())
	fmt.Println(time.Date(1990, 1, 1, 20, 0, 0, 0, time.Now().Location()).Unix())
}
