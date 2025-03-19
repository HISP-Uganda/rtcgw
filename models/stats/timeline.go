package stats

type ChartConfig struct {
	XAxisCategories []string     `json:"xAxisCategories"` // Each day represented once
	Series          []SeriesData `json:"series"`          // Data series for the days
}

type SeriesData struct {
	Name string   `json:"name"` // e.g., "Created TEs" or "Updated Events"
	Data []string `json:"data"` // One data point per day
}
