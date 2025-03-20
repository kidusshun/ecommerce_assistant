package chat

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/llmclient"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) WriteMessage(message string, sessionID uuid.UUID) error {
	_, err := s.db.Query("INSERT INTO chat (message, session_id) VALUES (?,?)", message, sessionID)
	if err != nil {
		return err
	}
	return nil
}

var Messsages []llmclient.Message
var MessageHistory []ChatResponse

func (s *Store) GetChatHistory(sessionID uuid.UUID) ([]llmclient.Message, error) {
	return Messsages, nil
}

