package dnsproxy

import (
	"context"
	"fmt"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/miekg/dns"
	"net"
	"strings"
)

// handleAAAA handles session registration. Registration looks like this:
// <ProjectID>.<SessionID>.<domain>
func (h *Handler) handleAAAA(msg *dns.Msg, r *dns.Msg) error {
	dataFromMessage := h.removeDomain(r.Question[0].Name)
	if dataFromMessage == "" {
		return fmt.Errorf("invalid domain")
	}

	ids := strings.Split(dataFromMessage, ".")
	if len(ids) != 2 {
		return fmt.Errorf("invalid data")
	}

	// grpc call to register session
	_, err := h.grpcClient.RegisterNewSession(context.Background(), &gateway.RegisterNewSessionRequest{
		SessionId: ids[1],
		ProjectId: ids[0],
	})
	if err != nil {
		return err
	}

	msg.Answer = append(msg.Answer, &dns.AAAA{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		AAAA: net.ParseIP("4efc:3425:412b:0bc0:99a1:9f87:d39f:9c84"),
	})

	return nil
}
