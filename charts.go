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
	Count() int
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
	data          HitSet
}

func NewChart(data HitSet, opts ...ChartOption) (*chart, error) {
	out := &chart{data: data}
	for _, o := range opts {
		if err := o(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

type ChartOption func(*chart) error

func Dimensions(height, width int) ChartOption {
	return func(c *chart) error {
		if height < 0 || width < 0 {
			return fmt.Errorf("invalid chart dimensions")
		}
		c.height = height
		c.width = width
		return nil
	}
}

func SparklineSVG(w io.Writer, data *HitSet, d time.Duration) error {
	hits := gochart.TimeSeries{
		//Name: "Hits",
		Style: gochart.Style{
			Show:        true,
			StrokeWidth: 2.0,
			StrokeColor: drawing.Color{R: 0, G: 0, B: 255, A: 100},
		},
	}

	data.SortByDate()
	var xVals []time.Time
	var yVals []float64
	tmp := data.XSeries()
	fmt.Println("xseries", len(tmp))
	direction := 0
	lastV := float64(0)
	for i := range tmp {
		v := tmp[i].Count()
		switch {
		case i == 0:
			fallthrough
		case v > tmp[i-1].Count():
			direction = 1
		case v < tmp[i-1].Count():
			direction = -1
		default:
			direction = 0
		}
		if direction != 0 {
			yVals = append(yVals, float64(v))
			lastV = float64(v)
		} else {
			yVals = append(yVals, lastV)
		}
		xVals = append(xVals, tmp[i].Time())
	}
	hits.XValues, hits.YValues = xVals, yVals

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

func HitChartSVG(w io.Writer, data *HitSet, d time.Duration) error {
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
