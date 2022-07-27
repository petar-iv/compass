package cronjob

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCronJob(t *testing.T) {
	templateCronJob := CronJob{
		Name:           "Test",
		Fn:             nil,
		SchedulePeriod: time.Nanosecond,
	}
	cfgNoElection := ElectionConfig{ElectionEnabled: false}

	t.Run("Should start executing cronJob", func(t *testing.T) {
		const expectedCronJobRuns = 3
		cronJob := templateCronJob
		counter := 0
		ctx, cancel := context.WithCancel(context.Background())

		cronJob.Fn = func(ctx context.Context) {
			counter += 1
			if counter == expectedCronJobRuns {
				cancel()
			}
		}
		RunCronJob(ctx, cfgNoElection, cronJob)

		assert.Equal(t, counter, expectedCronJobRuns)
	})

}
