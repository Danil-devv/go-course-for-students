package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
)

func NewService(a app.App) *AdService {
	return &AdService{app: a}
}

type AdService struct {
	app app.App
}

func (s *AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.app.CreateAd(request.Title, request.Text, request.UserId)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id: ad.ID, Title: ad.Title,
		Text:      ad.Text,
		Published: ad.Published,
		AuthorId:  ad.AuthorID,
	}, nil
}

func (s *AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.app.ChangeAdStatus(request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id: ad.ID, Title: ad.Title,
		Text:      ad.Text,
		Published: ad.Published,
		AuthorId:  ad.AuthorID,
	}, nil
}

func (s *AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.app.UpdateAd(request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id: ad.ID, Title: ad.Title,
		Text:      ad.Text,
		Published: ad.Published,
		AuthorId:  ad.AuthorID,
	}, nil
}

func (s *AdService) ListAds(ctx context.Context, empty *emptypb.Empty) (*ListAdResponse, error) {
	ads, err := s.app.GetAds()
	if err != nil {
		return nil, err
	}
	res := ListAdResponse{
		List: make([]*AdResponse, 0),
	}
	for _, ad := range ads {
		res.List = append(res.List, &AdResponse{
			Id: ad.ID, Title: ad.Title,
			Text:      ad.Text,
			Published: ad.Published,
			AuthorId:  ad.AuthorID,
		})
	}
	return &res, nil
}

func (s *AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := s.app.CreateUser(request.Id, request.Name, request.Email)
	if err != nil {
		return nil, err
	}
	return &UserResponse{Name: user.Nickname, Id: user.ID}, nil
}

func (s *AdService) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := s.app.GetUser(request.Id)
	if err != nil {
		return nil, err
	}
	return &UserResponse{Name: user.Nickname, Id: user.ID}, nil
}

func (s *AdService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := s.app.DeleteUser(request.Id)
	return &emptypb.Empty{}, err
}

func (s *AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	_, err := s.app.DeleteAd(request.AdId, request.AuthorId)
	return &emptypb.Empty{}, err
}
