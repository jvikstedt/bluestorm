package main

import (
	"log"
	"os"

	"github.com/jvikstedt/bluestorm"
	"github.com/jvikstedt/bluestorm/network"
	"github.com/jvikstedt/bluestorm/network/json"
)

type Greet struct {
	Name string `json:"name"`
}

func (g *Greet) Run(a *network.Agent, i interface{}) {
	greet, ok := i.(*Greet)
	if !ok {
		log.Println("not right type")
	}

	log.Printf("Greetings %s\n", greet.Name)
	a.Conn.Write([]byte(`{"Greet": { "name": "Alice" }}`))
}

type Tester interface {
	Run(a *network.Agent, i interface{})
}

func test(a *network.Agent, i interface{}) {
	tester, ok := i.(Tester)
	if !ok {
		log.Println("not right type")
	}

	tester.Run(a, i)
}

func main() {
	processor := json.NewProcessor()
	processor.Register(&Greet{}, test)

	servers := bluestorm.Servers{
		&network.WSServer{
			Addr:      ":8081",
			Processor: processor,
		},
	}

	bluestorm.Run(bluestorm.CloseOnSignal(os.Interrupt), servers)
}
