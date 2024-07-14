package dto

var Questions = []Question{
	{
		ID:     0,
		Detail: "What is the capital of France?",
		Options: []struct {
			Value    string `json:"value"`
			Label    string `json:"label"`
			ImageURL string `json:"imageUrl"`
		}{
			{Value: "A", Label: "Paris", ImageURL: "https://example.com/paris.jpg"},
			{Value: "B", Label: "London", ImageURL: "https://example.com/london.jpg"},
			{Value: "C", Label: "Berlin", ImageURL: "https://example.com/berlin.jpg"},
			{Value: "D", Label: "Madrid", ImageURL: "https://example.com/madrid.jpg"},
		},
	},
	{
		ID:     1,
		Detail: "Which planet is known as the Red Planet?",
		Options: []struct {
			Value    string `json:"value"`
			Label    string `json:"label"`
			ImageURL string `json:"imageUrl"`
		}{
			{Value: "A", Label: "Earth", ImageURL: "https://example.com/earth.jpg"},
			{Value: "B", Label: "Mars", ImageURL: "https://example.com/mars.jpg"},
			{Value: "C", Label: "Jupiter", ImageURL: "https://example.com/jupiter.jpg"},
			{Value: "D", Label: "Saturn", ImageURL: "https://example.com/saturn.jpg"},
		},
	},
}
