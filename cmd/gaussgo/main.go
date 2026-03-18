package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Unit struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Concepts    []Concept `json:"concepts"`
}

type Concept struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Explanation string   `json:"explanation"`
	Exercise    Exercise `json:"exercise"`
}

type Exercise struct {
	Type     string `json:"type"`
	Question string `json:"question,omitempty"`
}

func main() {
	fmt.Println("🚀 GaussGo - Linear Algebra CLI Tutor")
	fmt.Println("=====================================")

	data, err := os.ReadFile("data/units/vectors.json")
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	var unit Unit
	if err := json.Unmarshal(data, &unit); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Loaded unit: %s\n", unit.Title)
	fmt.Printf("Concepts: %d\n", len(unit.Concepts))
	fmt.Println("\nData-driven learning ready. Next: implement Bubble Tea TUI.")

	fmt.Println("\nRun with Docker for full Go 1.23+ support:")
	fmt.Println("  docker build -t gaussgo . && docker run -it gaussgo")
}
