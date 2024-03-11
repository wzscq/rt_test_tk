package device

import (
	"fmt"
	"testing"
)

func TestGetImsi(t *testing.T) {
	deviceClient := DeviceClient{
		ServerUrl: "http://182.42.81.6:5000/",
	}

	tc, err := deviceClient.GetImsi()
	if err != nil {
		t.Error("GetImsi error")
		return
	}

	fmt.Println(tc)
}

func TestGetDialQuery(t *testing.T) {
	deviceClient := DeviceClient{
		ServerUrl: "http://182.42.81.6:5000/",
	}

	tc, err := deviceClient.GetDialQuery()
	if err != nil {
		t.Error("GetDialQuery error")
		return
	}

	fmt.Println(tc)
}

func TestDialTrigger(t *testing.T) {
	deviceClient := DeviceClient{
		ServerUrl: "http://182.42.81.6:5000/",
	}

	tc, err := deviceClient.DialTrigger()
	if err != nil {
		t.Error("DialTrigger error")
		return
	}

	fmt.Println(tc)
}
