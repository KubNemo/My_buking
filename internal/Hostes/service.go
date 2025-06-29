package hostes

import (
	"context"
	"duking/internal/models"
)

type Service interface {
	HotelCreate(ctx context.Context, hotel *models.Hotel) error
	HotelGetOne(ctx context.Context, id uint) (*models.Hotel, error)
	HotelGetAll(ctx context.Context) ([]models.Hotel, error)
	HotelUpdate(ctx context.Context, id uint, hotel *models.Hotel) (*models.Hotel, error)
	HotelDelete(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// HotelCreate implements Service.
func (s *service) HotelCreate(ctx context.Context, hotel *models.Hotel) error {
	if err := s.repo.Create(ctx, hotel); err != nil {
		return err
	}
	return nil
}

// HotelDelete implements Service.
func (s *service) HotelDelete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

// HotelGetAll implements Service.
func (s *service) HotelGetAll(ctx context.Context) ([]models.Hotel, error) {
	hotels, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

// HotelGetOne implements Service.
func (s *service) HotelGetOne(ctx context.Context, id uint) (*models.Hotel, error) {
	hotel, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

// HotelUpdate implements Service.
func (s *service) HotelUpdate(ctx context.Context, id uint, hotel *models.Hotel) (*models.Hotel, error) {
	hotel, err := s.repo.Update(ctx, id, hotel)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}
