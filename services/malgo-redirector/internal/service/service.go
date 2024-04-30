package service

import (
	"context"
	"fmt"
	log2 "github.com/VipWW/malgo-c2/services/common/log"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/VipWW/malgo-c2/services/malgo-redirector/internal/dnsproxy"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log2.Init(logrus.InfoLevel)
}

type Service struct {
	dnsTcpServer *dns.Server
	dnsUdpServer *dns.Server
}

func New(
	dnsAddr string,
	grpcClient gateway.GatewayServiceClient,
) Service {
	dnsHandler := dnsproxy.NewHandler(grpcClient)

	dnsTcpServer := &dns.Server{
		Addr:      dnsAddr,
		Net:       "tcp",
		Handler:   dnsHandler,
		UDPSize:   65535,
		ReusePort: true,
	}
	dnsUdpServer := &dns.Server{
		Addr:      dnsAddr,
		Net:       "udp",
		Handler:   dnsHandler,
		UDPSize:   65535,
		ReusePort: true,
	}

	return Service{
		dnsTcpServer: dnsTcpServer,
		dnsUdpServer: dnsUdpServer,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		fmt.Printf("Starting DNS TCP server on %s\n", s.dnsTcpServer.Addr)
		return s.dnsTcpServer.ListenAndServe()
	})
	errgrp.Go(func() error {
		fmt.Printf("Starting DNS UDP server on %s\n", s.dnsUdpServer.Addr)
		return s.dnsUdpServer.ListenAndServe()
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down DNS TCP server\n")
		return s.dnsTcpServer.Shutdown()
	})
	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down DNS UDP server\n")
		return s.dnsUdpServer.Shutdown()
	})

	return errgrp.Wait()
}
