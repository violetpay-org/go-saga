package main

import "github.com/violetpay-org/go-saga"

var orchestrator = saga.NewOrchestrator[ExampleTxContext](UnitOfWorkFactory)
var registry = saga.NewRegistry(orchestrator)
