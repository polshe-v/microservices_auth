package user

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/repository/user/converter"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/user/model"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	passwordColumn  = "password"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"

	bcryptCost = 12
)

var errQueryBuild = errors.New("failed to build query")

type repo struct {
	db *pgxpool.Pool
}

// NewRepository creates new object of repository layer.
func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.UserCreate) (int64, error) {
	// Hashing the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil {
		log.Printf("%v", err)
		return 0, errors.New("failed to process password")
	}

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, roleColumn, emailColumn, passwordColumn).
		Values(user.Name, user.Role, user.Email, hashedPassword).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return 0, errQueryBuild
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Printf("%v", err)
		return 0, errors.New("failed to create user")
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return nil, errQueryBuild
	}

	var user modelRepo.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("%v", err)
			return nil, errors.New("no user with given id")
		}
		log.Printf("%v", err)
		return nil, errors.New("failed to read user info")
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, user *model.UserUpdate) error {
	builderUpdate := sq.Update(tableName).
		SetMap(map[string]interface{}{
			nameColumn:      user.Name,
			emailColumn:     user.Email,
			roleColumn:      user.Role,
			updatedAtColumn: sq.Expr("NOW()"),
		}).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return errQueryBuild
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%v", err)
		return errors.New("failed to update user info")
	}
	log.Printf("Result: %v", res)
	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return errQueryBuild
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%v", err)
		return errors.New("failed to delete user")
	}
	log.Printf("Result: %v", res)
	return nil
}
