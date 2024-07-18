package frontend

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"slices"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/exp/trace"
)

type Trace struct {
	app.Compo
	Xaxis      []string
	Stats      Stats
	functions  map[string]*Function
	goroutines map[int64]*Goroutine
}

func (h *Trace) OnMount(ctx app.Context) {
	h.functions = make(map[string]*Function)
	h.goroutines = make(map[int64]*Goroutine)
}

func (h *Trace) Render() app.UI {
	// Start reading from STDIN.
	r, err := trace.NewReader(bytes.NewReader(h.Stats.TraceData))
	if err != nil {
		return app.Div().ID("trace").Text(fmt.Sprintf("failed to read trace: %v", err))
	}

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
			ev.StateTransition().Stack.Frames(func(f trace.StackFrame) bool {
				if f.Func != "" {
					stack = f.Func
				} else if stack == "" {
					stack = fmt.Sprintf("%s:%d", f.File, f.Line)
				}
				return true
			})
			if ev.StateTransition().Resource.Kind == trace.ResourceProc {
				_, to := ev.StateTransition().Proc()
				h.functions[stack] = &Function{
					Name:      stack,
					State:     to.String(),
					Timestamp: int64(ev.Time()),
				}
			}
			if ev.StateTransition().Resource.Kind == trace.ResourceGoroutine {
				_, gto := ev.StateTransition().Goroutine()
				if gto.String() == "NotExist" {
					delete(h.goroutines, int64(ev.StateTransition().Resource.Goroutine()))
					continue
				}
				h.goroutines[int64(ev.StateTransition().Resource.Goroutine())] = &Goroutine{
					ID:        int64(ev.StateTransition().Resource.Goroutine()),
					State:     gto.String(),
					Reason:    ev.StateTransition().Reason,
					Timestamp: int64(ev.Time()),
				}
			}
		}
	}
	// Print what we found.
	var functions []*Function
	for _, ev := range h.functions {
		functions = append(functions, ev)
	}
	var goroutines []*Goroutine
	for _, ev := range h.goroutines {
		goroutines = append(goroutines, ev)
	}
	slices.SortFunc(functions, func(i, j *Function) int {
		return cmp.Compare(i.Name, j.Name)
	})
	slices.SortFunc(goroutines, func(i, j *Goroutine) int {
		return cmp.Compare(i.ID, j.ID)
	})
	return app.Div().ID("trace").Body(
		app.Range(functions).Slice(func(i int) app.UI {
			return functions[i]
		}),
		app.Range(goroutines).Slice(func(i int) app.UI {
			return goroutines[i]
		}),
	)
}
