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
	var users []models.User
	query := `
        SELECT user_address, MAX(CAST(user_nonce AS INTEGER)) AS count
        FROM events
        GROUP BY user_address
        ORDER BY count DESC
		LIMIT 10
    `
	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Address, &user.Count)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
