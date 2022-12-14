package replicator

import (
	"context"
	"fmt"
	"time"

	naisiov1 "nais/replicator/api/v1"

	corev1 "k8s.io/api/core/v1"
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
	go func() {
		for {
			<-time.After(r.syncInterval)
			var rcs naisiov1.ReplicatorConfigurationList
			err := r.client.List(ctx, &rcs, &client.ListOptions{})
			if err != nil {
				fmt.Println("Error listing ReplicatorConfigurations: ", err)
			}
			err = r.client.List(ctx, &corev1.NamespaceList{}, &client.ListOptions{})
			if err != nil {
				fmt.Println("Error listing ReplicatorConfigurations: ", err)
			}
			fmt.Println("Found ", len(rcs.Items), " ReplicatorConfigurations")
		}
	}()
	return nil
}
