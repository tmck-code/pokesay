package timer

import (
	"fmt"
	"encoding/json"
	"os"
	"time"
)

type Timer struct {
	Start          time.Time
	stageNames     []string
	StageTimes     map[string]time.Time
	StageDurations map[string]time.Duration
	Total          int64
}

func NewTimer() *Timer {
	now := time.Now()
	t := &Timer{
		Start:          now,
		stageNames:     make([]string, 0),
		StageTimes:     make(map[string]time.Time),
		StageDurations: make(map[string]time.Duration),
		Total:          0,
	}
	t.Mark("Start")
	return t
}

func (t *Timer) Mark(stage string) {
	now := time.Now()

	stage = fmt.Sprintf("%d.%s", len(t.stageNames), stage)

	t.StageTimes[stage] = now
	if len(t.stageNames) == 0 {
		t.StageDurations[stage] = 0
	} else {
		t.StageDurations[stage] = now.Sub(t.StageTimes[t.stageNames[len(t.stageNames)-1]])
	}
	t.stageNames = append(t.stageNames, stage)
}

func (t *Timer) Stop() {
	// From the last stage, subtract the first stage, to get total duration
	t.Total = t.StageTimes[t.stageNames[len(t.stageNames)-1]].Sub(t.StageTimes[t.stageNames[0]]).Nanoseconds()
}

func (t *Timer) PrintJson() {
	json.NewEncoder(os.Stderr).Encode(t)
}
