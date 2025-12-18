package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/srinivasarynh/age_calculator/internal/models"
	"github.com/srinivasarynh/age_calculator/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidDate  = errors.New("invalid date format")
)

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (*models.UserResponse, error)
	ListUsers(ctx context.Context, params *models.PaginationParams) (*models.UserListResponse, error)
	UpdateUser(ctx context.Context, id int32, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid DOB format", zap.Error(err))
		return nil, ErrInvalidDate
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format("2006-01-02"),
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (*models.UserResponse, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	age := CalculateAge(user.DOB)
	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format("2006-01-02"),
		Age:  &age,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, params *models.PaginationParams) (*models.UserListResponse, error) {
	params.SetDefaults()

	users, err := s.repo.List(ctx, params.GetLimit(), params.GetOffset())
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		age := CalculateAge(user.DOB)
		userResponses = append(userResponses, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  user.DOB.Format("2006-01-02"),
			Age:  &age,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	return &models.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid DOB format", zap.Error(err))
		return nil, ErrInvalidDate
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format("2006-01-02"),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return age
}
