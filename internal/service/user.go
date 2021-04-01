package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/pkg/auth"
	"github.com/aintsashqa/go-simple-blog/pkg/hash"
	uuid "github.com/satori/go.uuid"
)

type UserService struct {
	repo             repository.User
	hasher           hash.HashProvider
	auth             auth.AuthorizationProvider
	tokenExpiresTime time.Duration
}

func NewUserService(
	repo repository.User,
	hasher hash.HashProvider,
	auth auth.AuthorizationProvider,
	tokenExpiresTime time.Duration,
) *UserService {
	return &UserService{repo: repo, hasher: hasher, auth: auth, tokenExpiresTime: tokenExpiresTime}
}

func (s *UserService) SignUp(ctx context.Context, input SignUpUserInput) error {
	user := domain.User{
		Email:    input.Email,
		Username: fmt.Sprintf("Username%d", time.Now().Unix()),
		Password: s.hasher.Make(input.Password),
	}
	user.Init()

	return s.repo.Create(ctx, user)
}

func (s *UserService) SignIn(ctx context.Context, input SignInUserInput) (Tokens, error) {
	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return Tokens{}, err
	}

	if err := s.hasher.Compare(user.Password, input.Password); err != nil {
		return Tokens{}, err
	}

	accessToken, err := s.auth.NewToken(auth.TokenParams{
		UserID:    user.ID,
		ExpiresAt: s.tokenExpiresTime,
	})

	return Tokens{AccessToken: accessToken}, err
}

func (s *UserService) Find(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.Find(ctx, id)
}

func (s *UserService) Update(ctx context.Context, input UpdateUserInput) (domain.User, error) {
	user, err := s.repo.Find(ctx, input.ID)
	if err != nil {
		return domain.User{}, err
	}

	user.Username = input.Username
	user.Update()

	if err := s.repo.Update(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, input AuthenticateUserInput) (uuid.UUID, error) {
	return s.auth.Parse(input.Token)
}
