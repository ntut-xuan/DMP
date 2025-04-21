package cardset_test

import (
	"clyde1811/dmp/cardset"
	"fmt"
	"strconv"
	"testing"
)

func TestParseCardByString(t *testing.T) {
	c, err := cardset.ParseCardByString("C_9_0x35ab959a9c849c5a")

	if err != nil {
		fmt.Errorf("Failed: Parse card.")
	}

	if c.Suite != "C" {
		fmt.Errorf("Failed: Have wrong suite.")
	}

	if c.Number != "9" {
		fmt.Errorf("Failed: Have wrong number.")
	}

	randomNumber := c.RandomNumber
	excpetedRandomNumber, err := strconv.ParseInt("35ab959a9c849c5a", 16, 64)

	if randomNumber != excpetedRandomNumber {
		fmt.Errorf("Failed: Have wrong number.")
	}
}
