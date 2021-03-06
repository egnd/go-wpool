package wpool_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/egnd/go-wpool/v2"
	"github.com/egnd/go-wpool/v2/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Worker_Do(t *testing.T) {
	cases := []struct {
		buffSize int
		tasksCnt int
	}{
		{
			buffSize: 10,
			tasksCnt: 21,
		},
		{
			buffSize: 20,
			tasksCnt: 10,
		},
	}
	for k, test := range cases {
		t.Run(fmt.Sprint(k), func(tt *testing.T) {
			worker := wpool.NewWorker(test.buffSize)

			var wg sync.WaitGroup
			for i := 0; i <= test.tasksCnt; i++ {
				wg.Add(1)

				task := &interfaces.MockTask{}
				defer task.AssertExpectations(tt)
				task.On("Do").Once().After(time.Duration(rand.Intn(10)) * time.Millisecond).Run(func(_ mock.Arguments) { wg.Done() })

				assert.NoError(tt, worker.Do(task))
			}

			wg.Wait()
			assert.NoError(tt, worker.Close())
		})
	}
}

func Test_Worker_Close(t *testing.T) {
	worker := wpool.NewWorker(0)
	worker.Close()
	assert.EqualError(t, worker.Do(nil), "worker is closed")
}
