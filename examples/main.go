package main

import (
	"github.com/stephenrlouie/service"
	"github.com/stephenrlouie/service/examples/hello"
	"github.com/stephenrlouie/service/examples/sleep"
)

func main() {
	sg := service.New()
	svc1 := &hello.Hello{
		Id: "hello1",
	}

	sg.Add(svc1)
	svc2 := &sleep.Sleep{
		Id:      "sleep1",
		Pass:    true,
		Seconds: 5,
	}
	sg.Add(svc2)

	svc3 := &sleep.Sleep{
		Id:      "sleep2",
		Pass:    true,
		Seconds: 2,
	}
	sg.Add(svc3)

	sg.Start()
	sg.Wait()
}
