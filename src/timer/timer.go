package timer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Timer struct {
	stageNames     []string
	StageTimes     map[string]time.Time
	StageDurations map[string]time.Duration
	Total          int64
}

func NewTimer() *Timer {
	t := &Timer{
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
	t.stageNames = append(t.stageNames, stage)
}

func (t *Timer) Stop() {
	for i, stage := range t.stageNames {
		if i == 0 {
			t.StageDurations[stage] = 0
		} else {
			t.StageDurations[stage] = t.StageTimes[stage].Sub(t.StageTimes[t.stageNames[i-1]])
		}
	}
	// From the last stage, subtract the first stage, to get total duration
	t.Total = t.StageTimes[t.stageNames[len(t.stageNames)-1]].Sub(t.StageTimes[t.stageNames[0]]).Nanoseconds()
}

func (t *Timer) PrintJson() {
	json, _ := json.MarshalIndent(t, "", strings.Repeat(" ", 2))
	fmt.Println(string(json))
}
