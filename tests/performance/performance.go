package performance

import (
	"context"
	"fmt"
	"math"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
)

type DataPoint struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Request   any
	Response  any
	Error     errwrap.Error
}

type Statistics struct {
	Count                uint
	SvcUnderLoadDuration time.Duration
	P50                  time.Duration
	P90                  time.Duration
	P95                  time.Duration
	P99                  time.Duration
	Avg                  time.Duration
	Std                  time.Duration
}

type Pass struct {
	P50 bool
	P90 bool
	P95 bool
	P99 bool
}

type Report struct {
	Name              string
	RPS               uint
	DurationSec       uint
	RequestCount      uint
	NFRMs             uint
	StatisticsSuccess Statistics
	StatisticsError   Statistics
	ThroughPut        float64 // extrapolated req/s
	Pass              Pass
	Datapoints        []DataPoint
	Notes             []string
}

type Plan struct {
	T               *testing.T
	RPS             uint
	DurationSec     uint
	RampDurationSec uint
	Setup           func() (ctxFunc func() context.Context, request any, err errwrap.Error)
	Test            func(ctx context.Context, request any, errIn errwrap.Error) (response any, errOut errwrap.Error)
	NFRMs           uint
	Notes           []string
}

func lerp(x1 float64, x2 float64, y1 float64, y2 float64, x float64) float64 {
	if x2-x1 == 0 {
		return 0
	}
	return max(y1+((x-x1)*(y2-y1))/(x2-x1), 1)
}

type input struct {
	ctxFunc  func() context.Context
	request  any
	response any
	err      errwrap.Error
}

type setupsResults struct {
	rampup   [][]input
	test     [][]input
	rampdown [][]input
}

func (p *Plan) runSetups() setupsResults {
	var setupRes setupsResults
	if p.Setup == nil {
		return setupRes
	}

	for i := uint(0); i < p.RampDurationSec; i++ {
		count := uint(lerp(0, float64(p.RampDurationSec), 0, float64(p.DurationSec), float64(i+1)))

		stepInput := make([]input, count)
		for j := uint(0); j < count; j++ {
			ctxFunc, req, err := p.Setup()
			stepInput[j] = input{
				ctxFunc: ctxFunc,
				request: req,
				err:     err,
			}
		}
		setupRes.rampup = append(setupRes.rampup, stepInput)

	}

	setupRes.rampdown = make([][]input, len(setupRes.rampup))
	for i := 0; i < len(setupRes.rampup); i++ {
		setupRes.rampdown[len(setupRes.rampup)-1-i] = setupRes.rampup[i]
	}

	for i := uint(0); i < p.DurationSec; i++ {
		stepInput := make([]input, p.RPS)
		for j := uint(0); j < p.RPS; j++ {
			ctxFunc, req, err := p.Setup()
			stepInput[j] = input{
				ctxFunc: ctxFunc,
				request: req,
				err:     err,
			}
		}
		setupRes.test = append(setupRes.test, stepInput)
	}
	return setupRes
}

func avg(dps []DataPoint) time.Duration {
	var average float64
	for _, dp := range dps {
		average += float64(dp.Duration) / float64(len(dps))
	}
	return time.Duration(average)
}

func stddev(dps []DataPoint, average time.Duration) time.Duration {
	var variance float64
	sampleSize := len(dps) - 1
	for _, dp := range dps {
		variance += math.Pow(float64(dp.Duration-average), 2) / float64(sampleSize)
	}
	return time.Duration(math.Sqrt(variance))
}

func calcStatistics(datapoints []DataPoint) Statistics {
	if len(datapoints) == 0 {
		return Statistics{}
	}
	m := avg(datapoints)
	slices.SortFunc(datapoints, func(a DataPoint, b DataPoint) int {
		return int(a.Duration - b.Duration)
	})
	var loadDur time.Duration
	for _, dp := range datapoints {
		loadDur += dp.Duration
	}
	return Statistics{
		Count:                uint(len(datapoints)),
		SvcUnderLoadDuration: loadDur,
		P50:                  datapoints[len(datapoints)*50/100].Duration,
		P90:                  datapoints[len(datapoints)*90/100].Duration,
		P95:                  datapoints[len(datapoints)*95/100].Duration,
		P99:                  datapoints[len(datapoints)*99/100].Duration,
		Avg:                  m,
		Std:                  stddev(datapoints, m),
	}
}

func (p *Plan) Run() *Report {
	if p == nil || p.Test == nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: no plan or testing function defined; run aborted")
		return nil
	}

	p.T.Helper()

	setupRes := p.runSetups()

	iterate := func(iterSetups [][]input, collector chan DataPoint) {
		for i, setups := range iterSetups {
			s := time.Now()
			go func(ssetups []input) {
				dist := int64(time.Second) / int64(len(ssetups))
				for j, setup := range ssetups {
					ctx := setup.ctxFunc()
					start := time.Now()
					resp, err := p.Test(ctx, setup.request, setup.err)
					end := time.Now()
					dur := end.Sub(start)
					fmt.Printf("%v-%v: start: '%v', stop: '%v', err: %s\n", i, j, start, end, err)
					if collector != nil {
						collector <- DataPoint{
							StartTime: start,
							EndTime:   end,
							Duration:  dur,
							Request:   setup.request,
							Response:  resp,
							Error:     err,
						}
					}
					<-time.After(time.Duration(dist))
				}
			}(setups)
			<-time.After(max(1*time.Second-time.Since(s), 0)) // wait for next second
		}
	}

	collector := make(chan DataPoint, p.RPS*p.DurationSec)

	iterate(setupRes.rampup, nil)
	iterate(setupRes.test, collector)
	iterate(setupRes.rampdown, nil)

	close(collector)

	var datapoints []DataPoint
	var datapointsSuccess []DataPoint
	var datapointsError []DataPoint
	for dp := range collector {
		datapoints = append(datapoints, dp)
		if dp.Error == nil {
			datapointsSuccess = append(datapointsSuccess, dp)
		} else {
			datapointsError = append(datapointsError, dp)
		}
	}

	statsSuccess := calcStatistics(datapointsSuccess)
	statsError := calcStatistics(datapointsError)

	pass := Pass{
		P50: statsSuccess.P50.Milliseconds() < int64(p.NFRMs),
		P90: statsSuccess.P90.Milliseconds() < int64(p.NFRMs),
		P95: statsSuccess.P95.Milliseconds() < int64(p.NFRMs),
		P99: statsSuccess.P99.Milliseconds() < int64(p.NFRMs),
	}

	for _, edp := range datapointsError {
		p.Notes = append(p.Notes, fmt.Sprintf("trace_id: %v", edp.Error.GetTraceID()))
	}

	report := Report{
		Name:              p.T.Name(),
		RPS:               p.RPS,
		DurationSec:       p.DurationSec,
		RequestCount:      uint(len(datapoints)),
		NFRMs:             p.NFRMs,
		StatisticsSuccess: statsSuccess,
		StatisticsError:   statsError,
		ThroughPut:        float64(len(datapoints)) / float64(statsSuccess.SvcUnderLoadDuration.Seconds()+statsError.SvcUnderLoadDuration.Seconds()),
		Pass:              pass,
		Datapoints:        datapoints,
		Notes:             p.Notes,
	}
	return &report
}
