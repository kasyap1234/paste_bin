package repositories

import (
	"context"
	"fmt"
	"pastebin/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func(u*UserRepository)GetUserByID(ctx context.Context,userID uuid.UUID)(*models.User,error){
	query :=`SELECT * FROM users WHERE user_id=$1`
	rows,err :=u.db.Query(ctx,query,userID)
	if err !=nil{
		return nil,fmt.Errorf("failed to get user by id %w",err)
	}
	defer rows.Close()
	user,err :=pgx.CollectOneRow(rows,pgx.RowToStructByName[models.User])
	if err !=nil{
		return nil,fmt.Errorf("unable to collect row %w",err)

	}
return user,nil

}
func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email=$1`
	rows, err := u.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	defer rows.Close()
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &user, nil
}

func (u *UserRepository) ExistsUser(ctx context.Context, email string) (bool, error) {
	_, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to get user by email %s", email)
	}
	return true, nil
}

func(u*UserRepository)CheckUserByUserID(ctx context.Context,userID uuid.UUID)(bool,error){
	_,err :=
}
