package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a new Prometheus gauge
var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU in Celsius.",
})

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(cpuTemp)
}

func readCPUTemperature() (float64, error) {
	files, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil {
		return 0, err
	}

	if len(files) == 0 {
		return 0, fmt.Errorf("no thermal zone files found")
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Println("Error reading file:", err)
			continue
		}

		tempStr := strings.TrimSpace(string(content))
		tempMilli, err := strconv.Atoi(tempStr)
		if err != nil {
			log.Println("Error converting temperature:", err)
			continue
		}

		// Return temperature in Celsius
		return float64(tempMilli) / 1000.0, nil
	}

	return 0, fmt.Errorf("failed to read any temperature")
}

func updateTemperature() {
	temp, err := readCPUTemperature()
	if err != nil {
		log.Println("Error updating temperature:", err)
	} else {
		cpuTemp.Set(temp)
	}
}

func main() {
	// Expose the /metrics endpoint for Prometheus scraping
	http.Handle("/metrics", promhttp.Handler())

	// Periodically update the temperature
	go func() {
		for {
			updateTemperature()
			// Update every 5 seconds
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
