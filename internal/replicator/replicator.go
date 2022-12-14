package replicator

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Replicator struct {
	client       client.Client
	syncInterval time.Duration
}

func New(replicatorClient client.Client, syncInterval time.Duration) *Replicator {
	return &Replicator{
		client:       replicatorClient,
		syncInterval: syncInterval,
	}
}

func (r *Replicator) Run(ctx context.Context) error {
	return nil
}
