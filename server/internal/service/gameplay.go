package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mrspec7er/matchmind/server/internal/dto"
)

type Client struct {
	conn *websocket.Conn
}

type Service struct {
	Clients []*Client
	Scores  []*dto.Score
}

func (s *Service) WebsocketConnection() websocket.Upgrader {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return upgrader
}

func (s *Service) BroadcastMessage(msgType int, message []byte) {
	for _, client := range s.Clients {
		err := client.conn.WriteMessage(msgType, message)
		if err != nil {
			fmt.Println("error while broadcasting: ", err)
			client.conn.Close()
		}
	}
}

func (s *Service) CreateRoom() (*dto.Score, error) {
	newServer := &dto.Score{
		RoomID:     "ROOM_01",
		TotalWins:  0,
		QuestionID: 0,
	}
	s.Scores = append(s.Scores, newServer)

	return newServer, nil
}

func (s *Service) ProcessMessage(conn *websocket.Conn, roomId string) (int, error) {
	defer conn.Close()

	client := &Client{
		conn: conn,
	}
	s.Clients = append(s.Clients, client)

	defer func() {
		client.conn.Close()

		for i, c := range s.Clients {
			if c == client {
				s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
				break
			}
		}
	}()

	for {
		var resp *dto.Response

		msgType, data, err := conn.ReadMessage()
		if err != nil {
			return 500, err
		}

		err = json.Unmarshal(data, &resp)
		if err != nil {
			return 500, err
		}

		for i, score := range s.Scores {
			switch resp.Type {
			case "Question":
				s.RetrieveQuestion(msgType, s.Scores[i].QuestionID)

			case "Response":
				if score.RoomID == roomId {
					playerResponse, ok := resp.Content.(string)
					if !ok {
						s.SendMessage("Error", "Server", false, "Bad request")
					}
					score.Responses = append(score.Responses, playerResponse)
					s.MatchingResponse(msgType, score)
				}
			}
		}
	}
}

func (s *Service) RetrieveQuestion(msgType int, questionId int) error {
	question := s.FilterQuestion(questionId)
	s.SendMessage("Question", "Server", true, question)

	return nil
}

func (s *Service) FilterQuestion(questionId int) *dto.Question {
	for _, q := range dto.Questions {
		if q.ID == questionId {
			return &q
		}
	}

	return nil
}

func (s *Service) MatchingResponse(msgType int, score *dto.Score) {
	if len(score.Responses) < len(s.Clients) { // TODO checking client per room
		return
	}
	matchStatus := true

	for _, r := range score.Responses {
		if r != score.Responses[0] {
			matchStatus = false
		}

	}

	if matchStatus {
		score.TotalWins += 1
		// msgPayload, err := json.Marshal("Match")
		// if err != nil {
		// 	s.SendMessage("Error", "Server", false, "Cannot process result")
		// }
		// s.BroadcastMessage(msgType, msgPayload)

		s.SendMessage("Result", "Server", true, "Result Match")
	}

	if !matchStatus {
		// msgPayload, err := json.Marshal("Not Match")
		// if err != nil {
		// 	s.SendMessage("Error", "Server", false, "Cannot process result")
		// }
		// s.BroadcastMessage(msgType, msgPayload)
		s.SendMessage("Result", "Server", false, "Result Not Match")
	}

	// reset response
	score.Responses = []string{}
	score.QuestionID += 1

	fmt.Println("Question: ", score.QuestionID-1, "Total Win", score.TotalWins)

}

func (s *Service) SendMessage(respType string, sender string, status bool, content any) {
	payload := &dto.Response{
		Type:    respType,
		Sender:  sender,
		Status:  status,
		Content: content,
	}
	msgPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	s.BroadcastMessage(1, msgPayload)
}
