package repository

import (
	"github.com/jmoiron/sqlx"

	"github.com/klausborkowski/vote/models"
)

type Repository interface {
	InsertEvents(events []models.Event) error
	GetTopVoters() (*[]models.User, error)
}

type rep struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &rep{db}
}

func (r *rep) InsertEvents(events []models.Event) error {
	query := `
	INSERT INTO events (user_address, nft_ids, user_nonce, time,hash)
	VALUES (:user_address, :nft_ids, :user_nonce, :time,:hash)
	ON CONFLICT (hash) DO UPDATE
	SET nft_ids = excluded.nft_ids,
		user_nonce = excluded.user_nonce,
		time = excluded.time,
		user_address=excluded.user_address
`
	_, err := r.db.NamedExec(query, events)
	if err != nil {
		return err
	}

	return nil
}

func (r *rep) GetTopVoters() (*[]models.User, error) {
	// user := &models.User{}
	// err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	// if err != nil {
	// 	return nil, err
	// }
	// return user, nil
	return nil, nil
}
