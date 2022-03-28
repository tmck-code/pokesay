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
	t.stageNames = append(t.stageNames, stage)
}

func (t *Timer) StopTimer() {
	idx := 0
	for name, stage := range t.StageTimes {
		fmt.Println(name, stage)
		if idx == 0 {
			t.StageDurations[name] = 0
			idx += 1
			continue
		}
		t.StageDurations[name] = stage.Sub(t.StageTimes[t.stageNames[idx-1]])
		t.Total += t.StageDurations[name].Nanoseconds()
		idx += 1
	}
}

func (t *Timer) PrintJson() {
	json.NewEncoder(os.Stderr).Encode(t)
}
