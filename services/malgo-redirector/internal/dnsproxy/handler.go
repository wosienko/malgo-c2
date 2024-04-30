package dnsproxy

import (
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"strings"
)

const maxDNSMessageSize = 176

type Handler struct {
	grpcClient gateway.GatewayServiceClient
}

func NewHandler(
	grpcClient gateway.GatewayServiceClient,
) *Handler {
	return &Handler{
		grpcClient: grpcClient,
	}
}

// removeDomain removes the last 4 subdomains from the domain
// e.g. "projectId.sessionId.<4 random characters>.subdomain.domain.com." -> "projectId.sessionId"
func (h *Handler) removeDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 5 {
		return ""
	}
	return strings.Join(parts[:len(parts)-5], ".")
}
