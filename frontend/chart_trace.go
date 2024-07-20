package frontend

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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

	gmap := map[string]*Goroutine{}
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
			var stack string
			var stacks []string
			ev.StateTransition().Stack.Frames(func(f trace.StackFrame) bool {
				if f.Func == "" {
					return true
				}
				stacks = append(stacks, f.Func)
				return true
			})
			stack = strings.Join(stacks, " | ")
			if ev.StateTransition().Resource.Kind == trace.ResourceGoroutine {
				gID := fmt.Sprintf("%05d\n", int64(ev.StateTransition().Resource.Goroutine())*-1)
				_, gto := ev.StateTransition().Goroutine()
				if gto.String() == "NotExist" {
					delete(gmap, gID)
					continue
				}
				g, ok := gmap[gID]
				if ok {
					if stack != "" {
						g.Stack = stack
					}
					g.State = gto.String()
					g.Reason = ev.StateTransition().Reason
					g.Timestamp = int64(ev.Time())
				} else {
					gmap[gID] = &Goroutine{
						ID:        gID,
						Stack:     stack,
						State:     gto.String(),
						Reason:    ev.StateTransition().Reason,
						Timestamp: int64(ev.Time()),
					}
				}
			}
		}
	}
	// Print what we found.
	return app.Div().ID("trace").
		Style("margin-left", "50px").
		Body(
			app.Range(gmap).Map(func(id string) app.UI {
				return gmap[id]
			}),
		)
}
