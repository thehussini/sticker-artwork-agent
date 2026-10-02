package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type ArtworkMetadata struct {
	OrderID      string  `json:"order_id"`
	FileName     string  `json:"file_name"`
	DPI          int     `json:"dpi"`
	WidthInches  float64 `json:"width_inches"`
	HeightInches float64 `json:"height_inches"`
	ColorSpace   string  `json:"color_space"`
	HasBleed     bool    `json:"has_bleed"`
}

type VerificationResult struct {
	Approved       bool     `json:"approved"`
	Confidence     float64  `json:"confidence"`
	Issues         []string `json:"issues"`
	ActionRequired string   `json:"action_required"`
}

type ArtworkAgent struct {
	MinDPI int
}

func NewArtworkAgent(minDPI int) *ArtworkAgent {
	return &ArtworkAgent{MinDPI: minDPI}
}

func (a *ArtworkAgent) EvaluateArtwork(ctx context.Context, artwork ArtworkMetadata) (*VerificationResult, error) {
	log.Printf("[Agent] Analyzing Artwork for Order: %s (%s)", artwork.OrderID, artwork.FileName)

	var issues []string

	if artwork.DPI < a.MinDPI {
		issues = append(issues, fmt.Sprintf("Resolution too low: %d DPI (Minimum required: %d DPI)", artwork.DPI, a.MinDPI))
	}
	if artwork.ColorSpace != "CMYK" {
		issues = append(issues, fmt.Sprintf("Color space is %s; recommended print format is CMYK", artwork.ColorSpace))
	}
	if !artwork.HasBleed {
		issues = append(issues, "Missing print bleed margins around custom cut line")
	}

	result := &VerificationResult{
		Approved:   len(issues) == 0,
		Confidence: 0.98,
		Issues:     issues,
	}

	if result.Approved {
		result.ActionRequired = "Send directly to automated laser-cutting pipeline"
	} else {
		result.ActionRequired = "Flag for manual customer proofing review"
	}

	return result, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent := NewArtworkAgent(300)

	sampleArtwork := ArtworkMetadata{
		OrderID:      "ORD-2026-8891",
		FileName:     "die_cut_sticker_design.png",
		DPI:          150,
		WidthInches:  3.5,
		HeightInches: 3.5,
		ColorSpace:   "RGB",
		HasBleed:     false,
	}

	result, err := agent.EvaluateArtwork(ctx, sampleArtwork)
	if err != nil {
		log.Fatalf("Agent evaluation failed: %v", err)
	}

	outputBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("=== AI Agent Evaluation Summary ===")
	fmt.Println(string(outputBytes))
}
