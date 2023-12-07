package ports

import "github.com/Cosmin2410/proxy-backend-test/core/domain"

type ProxyLogRepository interface {
	SaveLog(*domain.ProxyLog) error
}
