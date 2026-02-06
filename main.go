package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/terjelafton/k8s-api-test/internal/client"
	"github.com/terjelafton/k8s-api-test/internal/types"
	"golang.org/x/sync/errgroup"
)

const iterations = 10

var contexts = []string{
	"default/api-crossnet-z2-staging-1-m5yc9s-apps-intilitycloud-com:443/terje.lafton@intility.no",
	"default/api-crossnet-z2-staging-2-u5fc1m-apps-intilitycloud-com:443/terje.lafton@intility.no",
}

type resourceStats struct {
	Count        int
	TotalFetch   time.Duration
	TotalSerialize time.Duration
}

func (s *resourceStats) AvgFetch() time.Duration     { return s.TotalFetch / time.Duration(iterations) }
func (s *resourceStats) AvgSerialize() time.Duration  { return s.TotalSerialize / time.Duration(iterations) }

type clusterStats struct {
	Context       string
	VRFs          resourceStats
	BGPVRFs       resourceStats
	RoutePolicies resourceStats
	PrefixSets    resourceStats
}

func addResult[T any](stats *resourceStats, r *client.FetchResult[T]) {
	stats.Count = r.Count
	stats.TotalFetch += r.FetchDuration
	stats.TotalSerialize += r.SerializeDuration
}

type clusterResult struct {
	Context       string
	VRFs          *client.FetchResult[types.VRF]
	BGPVRFs       *client.FetchResult[types.BGPVRF]
	RoutePolicies *client.FetchResult[types.RoutePolicy]
	PrefixSets    *client.FetchResult[types.PrefixSet]
}

func fetchCluster(ctx context.Context, c *client.Client) (*clusterResult, error) {
	res := &clusterResult{Context: c.Context}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		r, err := c.ListVRFs(ctx)
		if err != nil {
			return err
		}
		res.VRFs = r
		return nil
	})

	g.Go(func() error {
		r, err := c.ListBGPVRFs(ctx)
		if err != nil {
			return err
		}
		res.BGPVRFs = r
		return nil
	})

	g.Go(func() error {
		r, err := c.ListRoutePolicies(ctx)
		if err != nil {
			return err
		}
		res.RoutePolicies = r
		return nil
	})

	g.Go(func() error {
		r, err := c.ListPrefixSets(ctx)
		if err != nil {
			return err
		}
		res.PrefixSets = r
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

func printStats(resource string, s *resourceStats) {
	fmt.Printf("  %-14s %4d items | avg fetch %8s | avg serialize %8s\n",
		resource, s.Count, s.AvgFetch().Round(time.Millisecond), s.AvgSerialize().Round(time.Microsecond))
}

func main() {
	ctx := context.Background()

	clients := make([]*client.Client, len(contexts))
	for i, name := range contexts {
		c, err := client.NewClient(name)
		if err != nil {
			log.Fatalf("creating client for %q: %v", name, err)
		}
		clients[i] = c
	}

	stats := make([]clusterStats, len(contexts))
	for i, c := range clients {
		stats[i].Context = c.Context
	}

	for iter := range iterations {
		runStart := time.Now()

		results := make([]*clusterResult, len(clients))
		g, gctx := errgroup.WithContext(ctx)

		for i, c := range clients {
			g.Go(func() error {
				res, err := fetchCluster(gctx, c)
				if err != nil {
					return fmt.Errorf("cluster %q: %w", c.Context, err)
				}
				results[i] = res
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Run %d/%d: %s\n", iter+1, iterations, time.Since(runStart).Round(time.Millisecond))

		for i, r := range results {
			addResult(&stats[i].VRFs, r.VRFs)
			addResult(&stats[i].BGPVRFs, r.BGPVRFs)
			addResult(&stats[i].RoutePolicies, r.RoutePolicies)
			addResult(&stats[i].PrefixSets, r.PrefixSets)
		}
	}

	fmt.Printf("\nAverages over %d runs:\n\n", iterations)
	for _, s := range stats {
		fmt.Printf("[%s]\n", s.Context)
		printStats("VRFs", &s.VRFs)
		printStats("BGPVRFs", &s.BGPVRFs)
		printStats("RoutePolicies", &s.RoutePolicies)
		printStats("PrefixSets", &s.PrefixSets)
		fmt.Println()
	}
}
