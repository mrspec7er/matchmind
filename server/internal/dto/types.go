package dto

type Response struct {
	Type     string `json:"type"`
	PlayerID string `json:"playerId"`
	Response string `json:"response"`
}

type Score struct {
	RoomID     string   `json:"roomId"`
	Responses  []string `json:"responses"`
	TotalWins  int      `json:"totalWins"`
	QuestionID int      `json:"questionId"`
}
