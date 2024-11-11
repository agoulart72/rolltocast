package main

import (
	"fmt"
	"log"
	"net/http"
	"rolltocast/app"
	"strconv"
	"text/template"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	fmt.Println("Starting Server")

	http.HandleFunc("/", mainHandler)

	log.Fatal(http.ListenAndServe(":8100", nil))

}

func mainHandler(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Loading template
	var tmplFile = "stats.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	var level int = 1
	var runs int = 1
	var strategy app.Strategies = app.Strategies{
		MaxFirst: false,
	}
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		level, runs, strategy = parseForm(req)
	}

	r := app.RunSorcerer(level, strategy, runs)
	//		io.WriteString(w, "Hello, world!\n")
	err = tmpl.Execute(w, r)
	if err != nil {
		fmt.Fprintf(w, "Fatal Error: %v", err)
		return
	}

	// 1st Chart
	bar := generateBarChart(level, r, runs)
	bar.Render(w)

	// 2nd Chart
	line := generateLineChart(level, r, runs)
	line.Render(w)
}

func parseForm(req *http.Request) (int, int, app.Strategies) {

	var err error
	var level int = 1
	var runs int = 1
	var strategy app.Strategies = app.Strategies{
		MaxFirst: false,
	}

	level, err = strconv.Atoi(req.FormValue("level"))
	if err != nil {
		fmt.Println("Not a level")
		level = 1
	}
	runs, err = strconv.Atoi(req.FormValue("runs"))
	if err != nil {
		fmt.Println("Not a number")
		runs = 1
	}
	maxFirst := req.FormValue("strategy_max_first")
	if maxFirst == "on" {
		strategy.MaxFirst = true
	}

	return level, runs, strategy
}

func generateBarChart(level int, r app.Stats, runs int) *charts.Bar {
	// create a new bar instance
	bar := charts.NewBar()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Roll to Cast",
			Subtitle: fmt.Sprintf("Level %v", level),
		}),
		charts.WithColorsOpts(opts.Colors{"green", "lightgreen", "lightgray", "red", "pink"}),
	)

	// Put data into instance
	sData := make([]opts.BarData, 0)
	for _, sd := range r.NumberOfSuccess {
		sData = append(sData, opts.BarData{Value: sd / runs})
	}
	fData := make([]opts.BarData, 0)
	for _, fd := range r.NumberOfFailures {
		fData = append(fData, opts.BarData{Value: fd / runs})
	}
	maxSData := make([]opts.BarData, 0)
	for _, sd := range r.MaxSuccess {
		maxSData = append(maxSData, opts.BarData{Value: sd})
	}
	maxFData := make([]opts.BarData, 0)
	for _, fd := range r.MaxFailures {
		maxFData = append(maxFData, opts.BarData{Value: fd})
	}
	modeSData := make([]opts.BarData, 0)
	for _, msd := range r.SuccessTimes {
		maxSuccessTimes := 0
		mode := 0
		for c, st := range msd {
			if st > maxSuccessTimes {
				maxSuccessTimes = st
				mode = c
			}
		}
		modeSData = append(modeSData, opts.BarData{Value: mode})
	}

	bar.SetXAxis([]string{"1st", "2nd", "3rd", "4th", "5th", "6th", "7th", "8th", "9th"}).
		AddSeries("Avg Success", sData).
		AddSeries("Mode Success", modeSData).
		AddSeries("Max Success", maxSData).
		AddSeries("Failures", fData).
		AddSeries("Max Failures", maxFData)
	return bar
}

func generateLineChart(level int, r app.Stats, runs int) *charts.Line {

	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Roll to Cast",
			Subtitle: fmt.Sprintf("Level %v", level),
		}))

	// Put data into instance
	var xAxis []int
	for k := range r.SuccessTimes[0] {
		xAxis = append(xAxis, k)
	}

	line = line.SetXAxis(xAxis)

	for i := 0; i < 9; i++ {
		items := make([]opts.LineData, 0)
		for _, val := range r.SuccessTimes[i] {
			items = append(items, opts.LineData{Value: val})
		}
		if r.SuccessTimes[i][0] < runs {
			line = line.AddSeries(fmt.Sprintf("Level %d", i+1), items)
		}
	}

	// line = line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}
