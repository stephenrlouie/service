package main

import (
	"fmt"

	"github.com/stephenrlouie/service"
	"github.com/stephenrlouie/service/examples/sleep"
)

func main() {
	sg := service.New()

	// sg.Add(&hello.Hello{
	// 	Id: "hello1",
	// })
	//
	sg.Add(&sleep.Sleep{
		Id:      "sleep1",
		Pass:    false,
		Seconds: 2,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep2",
		Pass:    false,
		Seconds: 4,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep3",
		Pass:    false,
		Seconds: 6,
	})

	sg.Start()
	sg.Wait()

	for _, err := range sg.Status() {
		fmt.Printf("ERRS from SG: %v\n", err)
	}
}
