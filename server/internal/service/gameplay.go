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
				msgPayload, err := json.Marshal(fmt.Sprintf("Send question with id: %v", s.Scores[i].QuestionID))
				if err != nil {
					return 500, err
				}
				s.BroadcastMessage(msgType, msgPayload)
				continue

			case "Response":
				if score.RoomID == roomId {
					s.Scores[i].Responses = append(s.Scores[i].Responses, resp.Response)
				}

				if score.RoomID == roomId && len(score.Responses) == len(s.Clients) {

					matchStatus := true

					for _, r := range score.Responses {
						if r != score.Responses[0] {
							matchStatus = false
						}

					}

					if matchStatus {
						msgPayload, err := json.Marshal("Match")
						if err != nil {
							return 500, err
						}
						s.Scores[i].TotalWins += 1
						s.BroadcastMessage(msgType, msgPayload)
					}

					if !matchStatus {
						msgPayload, err := json.Marshal("Not Match")
						if err != nil {
							return 500, err
						}
						s.BroadcastMessage(msgType, msgPayload)
					}

					// reset score
					s.Scores[i].Responses = []string{}
					s.Scores[i].QuestionID += 1

					fmt.Println("Question: ", s.Scores[i].QuestionID-1, "Total Win", s.Scores[i].TotalWins)

				}

			}
		}
	}
}
