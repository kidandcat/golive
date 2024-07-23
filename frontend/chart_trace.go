package frontend

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/exp/trace"
)

type Trace struct {
	app.Compo
	chartComponent
	Xaxis     []string
	Stats     Stats
	functions map[string]*Function
	nodes     []opts.SankeyNode
	links     []opts.SankeyLink
}

func (h *Trace) OnMount(ctx app.Context) {
	h.functions = map[string]*Function{}
}

func (h *Trace) ProcessData(functions []*Function) {
	h.nodes = []opts.SankeyNode{}
	h.links = []opts.SankeyLink{}
	noduplicate := map[string]struct{}{}
	for _, f := range functions {
		ss := strings.Split(f.Stack, " | ")
		for _, s := range ss {
			if _, ok := noduplicate[s]; ok {
				continue
			}
			h.nodes = append(h.nodes, opts.SankeyNode{Name: fmt.Sprintf("%s (%v)", s, f.Status)})
			noduplicate[s] = struct{}{}
		}
		for i := 0; i < len(ss)-1; i++ {
			if _, ok := noduplicate[ss[i]+"-"+ss[i+1]]; ok {
				continue
			}
			h.links = append(h.links, opts.SankeyLink{
				Source: fmt.Sprintf("%s (%v)", ss[i], f.Status),
				Target: fmt.Sprintf("%s (%v)", ss[i+1], f.Status),
				Value:  float32(f.Time.Seconds()),
			})
			noduplicate[ss[i]+"-"+ss[i+1]] = struct{}{}
		}
	}
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
			// Get the stack.
			var stack string
			var stacks []string
			ev.StateTransition().Stack.Frames(func(f trace.StackFrame) bool {
				if f.Func == "" {
					return true
				}
				stacks = append(stacks, fmt.Sprintf("%s:%d", f.Func, f.Line))
				return true
			})
			m := len(stacks)
			if m > 5 {
				m = 5
			}
			stacks = stacks[:m]
			stack = strings.Join(stacks, " | ")
			// Process goroutine state transition events.
			if ev.StateTransition().Resource.Kind == trace.ResourceGoroutine {
				// Create a sortable goroutine ID.
				gid := int64(ev.StateTransition().Resource.Goroutine())
				gID := fmt.Sprintf("%05d\n", int64(ev.StateTransition().Resource.Goroutine())*-1)
				_, gto := ev.StateTransition().Goroutine()
				// Remove goroutines in the "NotExist" state.
				if gto.String() == "NotExist" {
					delete(gmap, gID)
					continue
				}
				g, ok := gmap[gID]
				var oldStack string
				if ok {
					oldStack = g.Stack
					// If the stack is not empty, update the goroutine stack.
					if stack != "" {
						g.Stack = stack
					}
					g.State = gto.String()
					g.Reason = ev.StateTransition().Reason
					g.Timestamp = int64(ev.Time())
				} else {
					// If the goroutine is not in the map, create it.
					g = &Goroutine{
						ID:        gID,
						Stack:     stack,
						State:     gto.String(),
						Reason:    ev.StateTransition().Reason,
						Timestamp: int64(ev.Time()),
					}
					gmap[gID] = g
				}
				for i, s := range stacks {
					s = strings.Join(stacks[i:], " | ")
					// If the goroutine is running a function, start the function.
					f, ok := h.functions[s]
					if !ok {
						f = &Function{
							Stack:   s,
							Started: map[int64]time.Time{},
						}
						h.functions[s] = f
					}
					if f.Started[gid] == (time.Time{}) {
						f.Started[gid] = time.Unix(0, g.Timestamp)
					}
					f.Status = g.State
					// If the goroutine is not running a function, stop the function.
					if g.Stack != oldStack && oldStack != "" {
						oldStacks := strings.Split(oldStack, " | ")
						for ii, ss := range oldStacks {
							ss = strings.Join(oldStacks[ii:], " | ")
							f := h.functions[ss]
							if f != nil && f.Started[gid] != (time.Time{}) {
								if time.Unix(0, g.Timestamp).Sub(f.Started[gid]) < 0 {
									panic("negative time")
								}
								f.Time += time.Unix(0, g.Timestamp).Sub(f.Started[gid])
								f.Started[gid] = time.Time{}
							}
						}
					}
				}
			}
		}
	}

	// Sort functions by time.
	var functions []*Function
	for _, f := range h.functions {
		functions = append(functions, f)
	}
	slices.SortFunc(functions, func(i, j *Function) int {
		return int(j.Time - i.Time)
	})

	// Keep only the top MAX_FUNCTIONS functions.
	if len(functions) > MAX_FUNCTIONS {
		functions = functions[:MAX_FUNCTIONS]
	}

	// Process the data.
	h.ProcessData(functions)
	sankey := charts.NewSankey()
	width, height := app.Window().Size()
	sankey.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "dark",
			Height: fmt.Sprintf("%dpx", height),
			Width:  fmt.Sprintf("%dpx", width),
		}),
		charts.WithAnimation(false),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)
	sankey.AddSeries("sankey", h.nodes, h.links,
		charts.WithLabelOpts(opts.Label{Show: opts.Bool(true)}),
		charts.WithSeriesAnimation(false))

	return h.renderChart(sankey)

	// 	app.Div().ID("goroutines").
	// 		Style("flex", "1").
	// 		Body(
	// 			app.Range(gmap).Map(func(id string) app.UI {
	// 				return gmap[id]
	// 			}),
	// 		),
}
