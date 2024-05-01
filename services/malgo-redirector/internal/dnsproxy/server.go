package dnsproxy

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/miekg/dns"
)

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	var err error

	switch r.Question[0].Qtype {
	case dns.TypeTXT:
		// command query
		err = h.handleTXT(msg, r)
	case dns.TypeAAAA:
		// session registration
		err = h.handleAAAA(msg, r)
	case dns.TypeA:
		// result retrieval
		err = h.handleA(msg, r)
	default:
		msg.SetRcode(r, dns.RcodeNameError)
	}

	if err != nil {
		log.FromContext(context.Background()).Errorf("failed to handle DNS query: %v", err)
		msg.SetRcode(r, dns.RcodeRefused)
	}

	w.WriteMsg(msg)
}
