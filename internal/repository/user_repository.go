package repository

import (
	"database/sql"
	"fmt"

	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepostory(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(req model.CreateUserRequest, hashedPassword string) (*model.User, error) {
	query := `INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING users (name, email, created_at)`

	user := &model.User{}
	err := r.db.QueryRow(query, req.Name, req.Email, hashedPassword).
		Scan(&user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("Erro ao cadastrar o usuario no banco: %w", err)
	}

	return user, nil
}

func (r *Repository) FindAllUsers() ([]model.User, error) {
	query := `SELECT id, name, email, created_at
		FROM users
		ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar os usuarios no banco: %w", err)
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("Erro ao ler usuario: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) FindByIdUser(id int) (*model.User, error) {
	query := `SELECT id, nome, email, created_at
		FROM users
		WHERE id = $1`

	user := &model.User{}

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Error ao filtrar o usuario pelo ID fornecido: %w", err)

	}

	return user, nil
}

func (r *Repository) UpdateUser(id int, req model.UpdateUserRequest) (*model.User, error) {
	query := `UPDATE users
		SET email = $1, password = $2
		WHERE id = $3
		REETUNING id, name, email`

	user := &model.User{}

	err := r.db.QueryRow(query, req.Email, req.Password, id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Erro ao atualizar os dados do usuario com o id fornecido: %w", err)
	}

	return user, nil
}

func (r *Repository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
