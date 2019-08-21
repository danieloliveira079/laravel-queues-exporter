// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary sparklinedemo displays a couple of SparkLine widgets.
// Exist when 'q' is pressed.
package widget

import (
	"context"
	queuesmetrics "github.com/danieloliveira079/laravel-queues-exporter/pkg/proto"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
)

type DashboardWidget struct {
	updateInterval time.Duration
	data           func() map[string]*queuesmetrics.Metric
	queuesWidgets  []*QueueWidget
}

func NewDashboardWidget(updateInterval time.Duration, data func() map[string]*queuesmetrics.Metric) *DashboardWidget {
	dashboard := &DashboardWidget{
		updateInterval: updateInterval,
		data:           data,
	}

	return dashboard
}

func (d *DashboardWidget) initWidgets() {
	data := d.data()
	for _, m := range data {
		queue := m.Queue
		qWidget := NewQueueWidget(queue, func(queue string) int {
			if q, ok := d.data()[queue]; ok {
				return int(q.Jobs)
			}

			return 0
		})
		d.queuesWidgets = append(d.queuesWidgets, qWidget)
	}
}

func (d *DashboardWidget) Render(ctx context.Context, cancel context.CancelFunc) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	d.initWidgets()

	grid := ui.NewGrid()
	gridRows := []interface{}{}
	for i, _ := range d.queuesWidgets {
		gridRows = append(gridRows, ui.NewRow(1.0/float64(len(d.queuesWidgets)), d.queuesWidgets[i]))
	}
	grid.Set(gridRows...)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(grid)

	ticker := time.NewTicker(d.updateInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ui.Render(grid)
			case <-ctx.Done():
				return
			}
		}
	}()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			cancel()
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			termWidth, termHeight := payload.Width, payload.Height
			grid.SetRect(0, 0, termWidth, termHeight-1)

			ui.Clear()
		}

	}
}
