package services

import (
	"github.com/Cosmin2410/proxy-backend-test/internal/core/domain"
	"github.com/Cosmin2410/proxy-backend-test/internal/core/ports"
)

type ProxyLogService struct {
	repo ports.ProxyLogRepository
}

func NewProxyLogService(repo ports.ProxyLogRepository) *ProxyLogService {
	return &ProxyLogService{repo: repo}
}

func (s *ProxyLogService) SaveLog(log *domain.ProxyLog) error {
	return s.repo.SaveLog(log)
}
