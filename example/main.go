package main

import (
	"fmt"
	"github.com/violetpay-org/go-saga/messageRelayer"
	"log"
	"time"
)

func main() {
	relayer := messageRelayer.New(1, channelRegistry, UnitOfWorkFactory)
	go messageRelayer.StartBatchRun(1*time.Second, relayer)

	saga := NewExampleSaga()
	saga.ApplySchemaTo(registry)

	err := registry.StartSaga(saga.Name(), map[string]interface{}{
		"id": "ExampleSaga-9cda1798-c3ad-463a-8a25-70c2416fed13",
	})
	if err != nil {
		log.Print(err)
	}

	go func() {
		for {
			load, err := exampleSessionRepository.Load("ExampleSaga-9cda1798-c3ad-463a-8a25-70c2416fed13")
			if err != nil {
				return
			}
			log.Print("", load.currentStep.Name()+" ", load.State(), load.IsPending())
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Hello, World!")

	<-time.After(2 * 1000 * time.Millisecond)

	<-time.After(2 * 1100 * time.Millisecond)

	log.Print("End relayer")
}
