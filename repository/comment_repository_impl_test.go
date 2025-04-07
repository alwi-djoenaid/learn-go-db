package repository

import (
	"context"
	"fmt"
	dbconn "learn-go-db/db"
	"learn-go-db/entity"
	"testing"
)

func TestCommentInsert(t *testing.T) {
	commentRepository := NewCommentRepository(dbconn.GetConnection())

	ctx := context.Background()
	comment := entity.Comment{
		Email:   "unittest@email.com",
		Comment: "Test comment",
	}

	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
