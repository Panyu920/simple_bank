package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	UserIp    string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	userIPHeader               = "x-forwarded-for"
)

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtda := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("%+v", md)
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtda.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtda.UserAgent = userAgents[0]
		}

		if userIPs := md.Get(userIPHeader); len(userIPs) > 0 {
			mtda.UserIp = userIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtda.UserIp = p.Addr.String()
	}

	return mtda
}
