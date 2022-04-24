package bin

import (
	"fmt"
	"os"
	"time"
	"github.com/schollz/progressbar/v3"
)

func NewProgressBar(max int, newlineMode bool) progressbar.ProgressBar {
	if newlineMode {
		max = -1
	}
	return *progressbar.NewOptions(
		max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "█", SaucerPadding: "░", BarStart: "╢", BarEnd: "╟"}),
	)
}
