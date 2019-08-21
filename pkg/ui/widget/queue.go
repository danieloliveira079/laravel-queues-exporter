package widget

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"math/rand"
	"time"
)

type QueueWidget struct {
	Queue string
	data  func(queue string) int
	*SparklineGroup
}

func (q *QueueWidget) update() {
	jobs := q.data(q.Queue)
	q.Lines[0].Data = append(q.Lines[0].Data, jobs)
	q.Title = fmt.Sprintf(" %s: %d ", q.Queue, jobs)
}

func NewQueueWidget(queue string, data func(queue string) int) *QueueWidget {
	queueWidget := new(QueueWidget)

	sentSparkline := NewSparkline()
	sentSparkline.Data = []int{}
	max := 7
	min := 0
	color := rand.Intn(max-min) + min

	sparkLine := NewSparkline()
	sparkLine.TitleColor = ui.ColorWhite
	sparkLine.LineColor = ui.Color(color)
	sparkLine.Data = []int{}

	spark := NewSparklineGroup(sparkLine)
	spark.PaddingLeft = 1
	spark.Title = queue
	queueWidget.SparklineGroup = spark
	queueWidget.Queue = queue
	queueWidget.data = data

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				queueWidget.update()
			}
		}
	}()

	return queueWidget
}
