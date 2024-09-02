package main

import (
	"fmt"
	"github.com/violetpay-org/go-saga/messageRelayer"
	"log"
	"time"
)

func main() {
	relayer := messageRelayer.New(1, channelRegistry, UnitOfWorkFactory)

	saga := NewExampleSaga()
	saga.ApplySchemaTo(registry)

	err := registry.StartSaga(saga.Name(), map[string]interface{}{
		"id": "ExampleSaga-9cda1798-c3ad-463a-8a25-70c2416fed13",
	})
	if err != nil {
		log.Print(err)
	}

	fmt.Println("Hello, World!")

	go messageRelayer.StartBatchRun(1*time.Second, relayer)

	<-time.After(1.5 * 1000 * time.Millisecond)

	log.Print("End relayer")

	load, err := exampleSessionRepository.Load("ExampleSaga-9cda1798-c3ad-463a-8a25-70c2416fed13")
	if err != nil {
		return
	}

	log.Print(load.currentStep.Name()+" ", load.State(), load.IsPending())
}
