package service

import (
	"errors"
	"fmt"

	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/model"
	r "github.com/hugaojanuario/crud_golang_testing_example_private/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{
	r *r.Repository
}

func NewService (r *r.Repository) *Service{
	return &Service{r: r}
}

func (s *Service) CreateUser(req model.CreateUserRequest)(*model.User, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	err != nil{
		return nil, fmt.Errorf("Erro ao encripitar a senha fornecida pelo usuario: %w", err)
	}

	user, err := s.r.CreateUser(req, string(hashedPassword))
	if err != nil{
		return nil, err
	}

	return user, nil
}

func (s *Service) FindAllUsers()([]model.User, error){
	return s.r.FindAllUsers();
}

func (s *Service) FindUserById(id int) (*model.User, error){
	user, err := s.r.FindByIdUser(id)
	if err != nil{
		return nil, err
	}
	if user == nil{
		return nil, errors.New("Erro: usuario nao encontrado")
	}

	return user, nil
}

func (s *Service) UpdateUser (id int, req model.UpdateUserRequest)(*model.User, error){
	existing, err := s.r.FindByIdUser(id)
	if err != nil{
		return nil, err
	}
	if existing == nil{
		return nil, errors.New("Erro: usuario nao encontrado")
	}

	return s.r.UpdateUser(id, req)
}

func (s *Service) DeleteUser(id int) error {
    err := s.repo.Delete(id)
    if errors.Is(err, sql.ErrNoRows) {
        return errors.New("usuário não encontrado")
    }

    return err
}
