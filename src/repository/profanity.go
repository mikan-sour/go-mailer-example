package repository

import (
	"context"
	"database/sql"
	"time"
)

type ProfanityRepo interface {
	GetProfanities(ctx context.Context, region string) ([]*BadWord, error)
}

type ProfanityRepoImpl struct {
	db *sql.Conn
}

type BadWord struct {
	Id        int
	Word      string
	Region    string
	Active    bool
	Created   time.Time
	CreatedBy string
	Updated   time.Time
	UpdatedBy string
}

var getProfanitiesQuery = `SELECT id, word, region, active, created, created_by FROM profanity WHERE region = $1 AND active = true;`

func NewRepo(ctx context.Context, db *sql.DB) (*ProfanityRepoImpl, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &ProfanityRepoImpl{db: conn}, nil
}

func (p *ProfanityRepoImpl) GetProfanities(ctx context.Context, region string) ([]*BadWord, error) {
	rows, err := p.db.QueryContext(ctx, getProfanitiesQuery, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badwords []*BadWord
	for rows.Next() {
		var b BadWord
		err = rows.Scan(&b.Id, &b.Word, &b.Region, &b.Active, &b.Created, &b.CreatedBy)
		if err != nil {
			return nil, err
		}
		badwords = append(badwords, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return badwords, nil
}
