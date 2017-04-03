package main

import (
	"github.com/stephenrlouie/service"
	"github.com/stephenrlouie/service/examples/hello"
	"github.com/stephenrlouie/service/examples/sleep"
)

func main() {
	sg := service.New()
	svc1 := &hello.Hello{}
	sg.Add(svc1)
	svc2 := &sleep.Sleep{
		Msg:     "sleep1",
		Pass:    true,
		Seconds: 20,
	}
	sg.Add(svc2)

	svc3 := &sleep.Sleep{
		Msg:     "sleep2",
		Pass:    false,
		Seconds: 2,
	}
	sg.Add(svc3)

	sg.Start()
	sg.Wait()
}
