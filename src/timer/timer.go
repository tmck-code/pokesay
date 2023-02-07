package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	DEBUG bool = os.Getenv("DEBUG") != ""
)

type Timer struct {
	Name             string
	stageNames       []string
	StageTimes       map[string]time.Time
	StageDurations   map[string]time.Duration
	StagePercentages map[string]string
	Total            int64
	AlignKeys        bool
	Enabled          bool
}

func NewTimer(name string, alignKeys ...bool) *Timer {
	align := false
	if len(alignKeys) == 1 {
		align = alignKeys[0]
	}
	t := &Timer{
		Name:             name,
		stageNames:       make([]string, 0),
		StageTimes:       make(map[string]time.Time),
		StageDurations:   make(map[string]time.Duration),
		StagePercentages: make(map[string]string),
		Total:            0,
		AlignKeys:        align,
		Enabled:          DEBUG,
	}
	t.Mark("Start")
	return t
}

func (t *Timer) Mark(stage string) {
	if !t.Enabled {
		return
	}
	now := time.Now()
	if t.AlignKeys {
		stage = fmt.Sprintf("%02d.%-15s", len(t.stageNames), stage)
	} else {
		stage = fmt.Sprintf("%02d.%s", len(t.stageNames), stage)
	}

	t.StageTimes[stage] = now
	t.stageNames = append(t.stageNames, stage)
}

func (t *Timer) Stop() {
	if !t.Enabled {
		return
	}
	for i, stage := range t.stageNames {
		if i == 0 {
			t.StageDurations[stage] = 0
		} else {
			t.StageDurations[stage] = t.StageTimes[stage].Sub(t.StageTimes[t.stageNames[i-1]])
		}
	}
	// From the last stage, subtract the first stage, to get total duration
	t.Total = t.StageTimes[t.stageNames[len(t.stageNames)-1]].Sub(t.StageTimes[t.stageNames[0]]).Nanoseconds()

	for i, stage := range t.stageNames {
		if i == 0 {
			t.StagePercentages[stage] = "0.0%"
		} else {
			t.StagePercentages[stage] = fmt.Sprintf(
				"%.2f%%",
				float64(t.StageDurations[stage].Nanoseconds())*100.0/float64(t.Total),
			)
		}
	}
}

func (t *Timer) PrintJson() {
	if !t.Enabled {
		return
	}
	json, _ := json.MarshalIndent(t, "", strings.Repeat(" ", 2))
	fmt.Println(string(json))
}
