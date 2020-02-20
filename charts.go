package sws

import (
	"fmt"
	"io"
	"time"

	gochart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type Chartable interface {
	YMax() int
	XSeries() []Countable
}

type Countable interface {
	Label() string
	YValue() int
}

/*
type TimeChartable interface {
	XMax() int
	Series() []TimeCountable
}

type TimeCountable interface {
	XValue() time.Time
	YValue() int
}
*/

type chart struct {
	width, height int
	data          TimeBuckets
}

func NewChart(data TimeBuckets, opts ...ChartOption) (*chart, error) {
	out := &chart{data: data}
	for _, o := range opts {
		if err := o(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

type ChartOption func(*chart) error

func Dimentions(height, width int) ChartOption {
	return func(c *chart) error {
		if height < 0 || width < 0 {
			return fmt.Errorf("invalid chart dimensions")
		}
		c.height = height
		c.width = width
		return nil
	}
}

func SparklineSVG(w io.Writer, data TimeBuckets, d time.Duration) error {
	hits := gochart.TimeSeries{
		//Name: "Hits",
		Style: gochart.Style{
			Show:        true,
			StrokeWidth: 2.0,
			StrokeColor: drawing.Color{R: 0, G: 0, B: 255, A: 100},
		},
	}

	hits.XValues, hits.YValues = data.XYValues()

	graph := gochart.Chart{
		Width:  300,
		Height: 50,
		Series: []gochart.Series{hits},
		// Background: gochart.Style{
		// 	Padding: gochart.Box{Top: 10, Right: 10, Bottom: 10, Left: 27},
		// },
		// XAxis: gochart.XAxis{
		// 	Style: gochart.Style{Show: true},
		// 	ValueFormatter: func(v interface{}) string {
		// 		return time.Unix(0, int64(v.(float64))).Format("Jan 02")
		// 	},
		// },
		// YAxis: gochart.YAxis{
		// 	NameStyle: gochart.StyleShow(),
		// 	Style:     gochart.Style{Show: true},
		// 	ValueFormatter: func(v interface{}) string {
		// 		return fmt.Sprintf("%.0f", v.(float64))
		// 	},
		// },
	}
	//graph.YAxis.Range.SetMin(0)

	graph.Render(gochart.SVG, w)
	return nil
}

func HitChartSVG(w io.Writer, data TimeBuckets, d time.Duration) error {
	hits := gochart.TimeSeries{
		Name: "Hits",
		Style: gochart.Style{
			Show:        true,
			StrokeWidth: 2.0,
			StrokeColor: drawing.Color{R: 0, G: 0, B: 255, A: 100},
			DotColorProvider: func(_, _ gochart.Range, _ int, _, _ float64) drawing.Color {
				return drawing.Color{R: 21, G: 198, B: 148, A: 100}
			},
			DotWidthProvider: func(_, _ gochart.Range, _ int, _, _ float64) float64 {
				return 5
			},
		},
		//YAxis:   gochart.YAxisSecondary,
	}

	hits.XValues, hits.YValues = data.XYValues()

	graph := gochart.Chart{
		Width:  400,
		Height: 200,
		Series: []gochart.Series{hits},
		Background: gochart.Style{
			Padding: gochart.Box{Top: 10, Right: 10, Bottom: 10, Left: 27},
		},
		XAxis: gochart.XAxis{
			Style: gochart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				return time.Unix(0, int64(v.(float64))).Format("2015-01-02 15:04")
			},
		},
		YAxis: gochart.YAxis{
			//NameStyle: gochart.StyleShow(),
			Style: gochart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f", v.(float64))
			},
		},
		// YAxis: gochart.YAxis{
		// 	Name:      "Unique visitors",
		// 	NameStyle: gochart.StyleShow(),
		// 	Style:     gochart.Style{Show: true},
		// 	ValueFormatter: func(v interface{}) string {
		// 		return fmt.Sprintf("%.1f", v.(float64))
		// 	},
		// },
	}
	//graph.YAxis.Range.SetMin(0)
	// graph.Elements = []gochart.Renderable{
	// 	gochart.Legend(&graph),
	// }

	graph.Render(gochart.SVG, w)
	return nil
}
