package timer

import (
	"encoding/json"
	"os"
	"time"
)

type finishedTimer struct {
	OrderedDurations []int64
	StageDurations   map[string]int64
}

type Timer struct {
	Start          time.Time
	StageNames     []string
	StageTimes     map[string]time.Time
	StageDurations map[string]time.Duration
	FinishedTimer  finishedTimer
}

func NewTimer() *Timer {
	now := time.Now()
	return &Timer{
		Start:          now,
		StageNames:     []string{"start"},
		StageTimes:     map[string]time.Time{"start": now},
		StageDurations: map[string]time.Duration{"start": now.Sub(now)},
	}
}

func (t *Timer) lastStageTime() time.Time {
	return t.StageTimes[t.StageNames[len(t.StageNames)-1]]
}

func (t *Timer) Mark(stage string) {
	now := time.Now()
	t.StageDurations[stage] = now.Sub(t.lastStageTime())
	t.StageTimes[stage] = now
	t.StageNames = append(t.StageNames, stage)
}

func (t *Timer) StopTimer() {
	t.FinishedTimer = finishedTimer{StageDurations: map[string]int64{}}
	var total int64 = 0
	for _, stage := range t.StageNames {
		n := t.StageDurations[stage].Nanoseconds()
		t.FinishedTimer.OrderedDurations = append(t.FinishedTimer.OrderedDurations, n)
		t.FinishedTimer.StageDurations[stage] = n
		total += n
	}
}

func (t *Timer) PrintJson() {
	json.NewEncoder(os.Stderr).Encode(t)
}
