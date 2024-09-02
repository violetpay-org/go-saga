package messageRelayer

import (
	"fmt"
	"github.com/thanos-io/thanos/pkg/runutil"
	"log"
	"time"
)

type BatchJob interface {
	Execute() error
}

type Logger string

func (l Logger) Log(keyvals ...interface{}) error {
	fmt.Println(fmt.Sprintf("[%s] ", l), keyvals)
	return nil
}

func StartBatchRun(interval time.Duration, job BatchJob) chan struct{} {
	log.Print("Starting batch run")
	haultChan := make(chan struct{})
	runutil.RepeatInfinitely(
		Logger("Batch"),
		interval,
		haultChan,
		job.Execute,
	)

	return haultChan
}
