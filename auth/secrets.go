package auth

import "os"

var ExpectedPassword string
var HMACSecret []byte

func init() {
	var ok bool
	ExpectedPassword, ok = os.LookupEnv("PASSWORD")
	if !ok {
		ExpectedPassword = "test"
	}

	stringSecret, ok := os.LookupEnv("SECRET")
	if ok {
		HMACSecret = []byte(stringSecret)
	} else
	{
		HMACSecret = []byte("test")
	}
}
