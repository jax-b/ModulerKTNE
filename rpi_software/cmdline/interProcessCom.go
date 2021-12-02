package cmdline

import (
	"fmt"
	"os"
	"strings"

	ipc "github.com/james-barrow/golang-ipc"
)

type InterProcessCom struct {
	clcon *ipc.Client
}

func NewIPC() *InterProcessCom {
	cc, err := ipc.StartClient("ktne-ipc", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &InterProcessCom{
		clcon: cc,
	}
}
func (sipc *InterProcessCom) Close() error {
	sipc.clcon.Close()
	return nil
}
func (sipc *InterProcessCom) SendCustom(msg string) (bool, error) {
	err := sipc.clcon.Write(3, []byte(msg))
	fmt.Println("Sent:", msg)
	if err != nil {
		return false, err
	}
	rmsg, err := sipc.clcon.Read()
	if err != nil {
		return false, err
	}
	fmt.Print("Received ")
	if rmsg.MsgType == 1 {
		fmt.Print("Status: ")
		fmt.Println(string(rmsg.Data))
		if strings.Contains(string(rmsg.Data), ",ok") {
			return true, nil
		} else {
			return false, nil
		}
	} else if rmsg.MsgType == 2 {
		fmt.Print("Data: ")
		fmt.Println(strings.Split(string(rmsg.Data), ":")[1])
		return true, nil
	} else {
		fmt.Println("Unknown")
		fmt.Println(string(rmsg.Data))
		return false, nil
	}
}
