package handlers

import (
	"encoding/base64"
	"html/template"
	rdb "main/ridership_db"
	"main/utils"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the selected chart from the query parameter
	selectedChart := r.URL.Query().Get("line")
	if selectedChart == "" {
		selectedChart = "red"
	}

	// instantiate ridershipDB
	var db rdb.RidershipDB = &rdb.SqliteRidershipDB{} // Sqlite implementation
	// var db rdb.RidershipDB = &rdb.CsvRidershipDB{} // CSV implementation

	// Get the chart data from RidershipDB
	err := db.Open("../mbta.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ridership, err := db.GetRidership(selectedChart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Plot the bar chart using utils.GenerateBarChart. The function will return the bar chart
	// as PNG byte slice. Convert the bytes to a base64 string, which is used to embed images in HTML.
	chart, err := utils.GenerateBarChart(ridership)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	imgChart := base64.StdEncoding.EncodeToString(chart)

	// Get path to the HTML template for our web app
	_, currentFilePath, _, _ := runtime.Caller(0)
	templateFile := filepath.Join(filepath.Dir(currentFilePath), "template.html")

	// Read and parse the HTML so we can use it as our web app template
	html, err := os.ReadFile(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.New("line").Parse(string(html))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We now want to create a struct to hold the values we want to embed in the HTML
	data := struct {
		Image string
		Chart string
	}{
		Image: imgChart,
		Chart: selectedChart,
	}

	// Use tmpl.Execute to generate the final HTML output and send it as a response
	// to the client's request.
	tmpl.Execute(w, data)
}
