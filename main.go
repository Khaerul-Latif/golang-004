package main

import (
	"encoding/json"
	"fmt"
	"os"
	"math/rand"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func updateStatus(filePath string) {
	for {
		water := rand.Intn(100) + 1
		wind := rand.Intn(100) + 1
		status := Status{
			Water: water,
			Wind:  wind,
		}

		data, err := json.MarshalIndent(status, "", "   ")
		if err != nil {
			fmt.Println("Error encoding status:", err)
			continue
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			fmt.Println("Error writing status to file:", err)
		}

		time.Sleep(15 * time.Second)
	}
}

func getStatus(filePath string) (Status, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return Status{}, err
	}

	var status Status
	err = json.Unmarshal(file, &status)
	if err != nil {
		return Status{}, err
	}

	return status, nil
}

func determineStatus(water, wind int) string {
	waterStatus := "Aman"
	windStatus := "Aman"

	if water < 5 {
		waterStatus = "Aman"
	} else if water >= 6 && water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	if wind < 6 {
		windStatus = "Aman"
	} else if wind >= 7 && wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	return fmt.Sprintf("Status Air: %s, Status Angin: %s", waterStatus, windStatus)
}

func main() {
	filePath := "status.json"

	go updateStatus(filePath)

	for {
		status, err := getStatus(filePath)
		if err != nil {
			fmt.Println("Error reading status:", err)
			os.Exit(1)
		}

		statusStr := determineStatus(status.Water, status.Wind)
		fmt.Println(statusStr)

		time.Sleep(2 * time.Second) // Update status every 2 seconds
	}
}
