package repository

import (
	"context"
	"database/sql"
	"errors"
	"learn-go-db/entity"
	"strconv"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repository *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	query := "INSERT INTO comments(email,  comment) values ($1, $2) returning id"
	res := repository.DB.QueryRowContext(ctx, query, comment.Email, comment.Comment)
	if res.Err() != nil {
		return comment, res.Err()
	}

	var lastInsertId int
	err := res.Scan(&lastInsertId)
	if err != nil {
		return comment, err
	}

	comment.Id = int32(lastInsertId)
	return comment, nil
}

func (repository *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments WHERE id = $1 LIMIT 1"
	res := repository.DB.QueryRowContext(ctx, query, id)
	comment := entity.Comment{}

	if res.Err() != nil {
		return comment, res.Err()
	}

	if res != nil {
		res.Scan(&comment.Id, &comment.Email, &comment.Comment)
		return comment, nil
	} else {
		return comment, errors.New("Id " + strconv.Itoa(int(id)) + " tidak ditemukan")
	}
}

func (repository *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments"
	res, err := repository.DB.QueryContext(ctx, query)
	var comments []entity.Comment

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		comment := entity.Comment{}
		res.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}

	return comments, nil
}
