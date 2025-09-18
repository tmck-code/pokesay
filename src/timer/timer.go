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
	DebugTimer *Timer = NewTimer("Debug Timer", true)
)

// Timer is a simple timer that can be used to time the duration of program stages.
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

// NewTimer creates a new Timer with the given name.
// If align is true, values will be vertially aligned in the JSON output
func NewTimer(name string, boolArgs ...bool) *Timer {
	align := false
	if len(boolArgs) == 1 {
		align = boolArgs[0]
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

// Mark records the current time as the end of a stage.
// e.g. to time a block of code:
// t = NewTimer("MyTimer");
// SomeCode();
// t.Mark("Some Code");
func (t *Timer) Mark(stage string) {
	if !t.Enabled {
		return
	}
	now := time.Now()
	if t.AlignKeys {
		stage = fmt.Sprintf("%02d.%-20s", len(t.stageNames), stage)
	} else {
		stage = fmt.Sprintf("%02d.%s", len(t.stageNames), stage)
	}

	t.StageTimes[stage] = now
	t.stageNames = append(t.stageNames, stage)
}

// Stop records the current time as the end of the last stage,
// calculates the duration of each stage, the total duration,
// and the percentage of the total duration that each stage took.
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

// Print prints the timer's stage names, durations, and percentages to stderr as
// indented JSON
func (t *Timer) PrintJson() {
	if !t.Enabled {
		return
	}
	json, _ := json.MarshalIndent(t, "", strings.Repeat(" ", 2))
	fmt.Fprintln(os.Stderr, string(json))
}
