package main

type (
	ServiceConfig struct {
	}

	Service struct {
		mode ModeEnum
	}
)

func NewService() *Service {
	var s Service

	return &s
}

func (s *Service) Run() {

}

func (s *Service) Shutdown() {

}
