package frontend

import (
	"bytes"
	"fmt"
	"io"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/exp/trace"
)

type Trace struct {
	app.Compo
	Xaxis []string
	Stats Stats
}

func (h *Trace) Render() app.UI {
	// Start reading from STDIN.
	r, err := trace.NewReader(bytes.NewReader(h.Stats.TraceData))
	if err != nil {
		return app.Div().ID("trace").Text(fmt.Sprintf("failed to read trace: %v", err))
	}

	var events []string
	for {
		// Read the event.
		ev, err := r.ReadEvent()
		if err == io.EOF {
			break
		} else if err != nil {
			return app.Div().ID("trace").Text(fmt.Sprintf("failed to read event: %v", err))
		}
		switch ev.Kind() {
		case trace.EventStateTransition:
			var from, to trace.ProcState
			var gfrom, gto trace.GoState
			if ev.StateTransition().Resource.Kind == trace.ResourceProc {
				from, to = ev.StateTransition().Proc()
			}
			if ev.StateTransition().Resource.Kind == trace.ResourceGoroutine {
				gfrom, gto = ev.StateTransition().Goroutine()
			}
			var stack string
			if !ev.StateTransition().Stack.Frames(func(f trace.StackFrame) bool {
				if f.Func != "" {
					stack = f.Func
				}
				return true
			}) {
				stack = "failed to get stack"
			}
			if stack == "" {
				continue
			}
			if from != to {
				events = append(events, fmt.Sprintf("(%v -> %v) %v / %v  %v", from.String(), to.String(), stack, ev.StateTransition().Resource.Kind, ev.Time()))
			} else {
				events = append(events, fmt.Sprintf("(%v -> %v) %v / %v  %v", gfrom.String(), gto.String(), stack, ev.StateTransition().Reason, ev.Time()))
			}
		}
	}
	// Print what we found.
	var lines []app.UI
	for _, ev := range events {
		lines = append(lines, app.Li().Text(fmt.Sprintf("%v", ev)))
	}
	return app.Div().ID("trace").Body(
		lines...,
	)
}
