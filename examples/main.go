package main

import (
	"github.com/stephenrlouie/service"
	"github.com/stephenrlouie/service/examples/hello"
	"github.com/stephenrlouie/service/examples/sleep"
)

func main() {
	sg := service.New()

	sg.Add(&hello.Hello{
		Id: "hello1",
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep1",
		Pass:    true,
		Seconds: 5,
	})

	sg.Add(&sleep.Sleep{
		Id:      "sleep2",
		Pass:    false,
		Seconds: 6,
	})

	sg.Start()
	sg.Wait()
}
