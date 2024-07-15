package dto

type Response struct {
	Type    string `json:"type"`
	Sender  string `json:"sender"`
	Status  bool   `json:"status"`
	Content any    `json:"content"`
}

type Score struct {
	RoomID     string   `json:"roomId"`
	Responses  []string `json:"responses"`
	TotalWins  int      `json:"totalWins"`
	QuestionID int      `json:"questionId"`
}

type Question struct {
	ID      int    `json:"id"`
	Detail  string `json:"detail"`
	Options []struct {
		Value    string `json:"value"`
		Label    string `json:"label"`
		ImageURL string `json:"imageUrl"`
	} `json:"options"`
}
