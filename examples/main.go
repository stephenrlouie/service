package main

import (
	"fmt"

	"github.com/stephenrlouie/service"
	"github.com/stephenrlouie/service/examples/hello"
	"github.com/stephenrlouie/service/examples/sleep"
)

func main() {
	sg := service.New()
	sg.HandleSigint(func() {
		fmt.Printf("SIGINT Received. Shutting down...\n")
	})

	sg.Add(&hello.Hello{
		Id: "hello",
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep-2",
		Pass:    true,
		Seconds: 2,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep-4",
		Pass:    false,
		Seconds: 4,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep-6",
		Pass:    true,
		Seconds: 6,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep-8",
		Pass:    true,
		Panic:   true,
		Seconds: 8,
	})

	sg.Start()
	sg.Wait()

	errs := sg.Status()
	if len(errs) != 0 {
		fmt.Printf("*** Service Group Errors ***\n")
		for i, err := range sg.Status() {
			fmt.Printf("\t%d: %v\n", i, err)
		}
	}
}
