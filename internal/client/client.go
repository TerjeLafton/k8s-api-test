package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/terjelafton/k8s-api-test/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

const namespace = "default"

var (
	group   = "xr.crossnet.network.intility.io"
	version = "v1"

	vrfGVR         = schema.GroupVersionResource{Group: group, Version: version, Resource: "vrfs"}
	bgpvrfGVR      = schema.GroupVersionResource{Group: group, Version: version, Resource: "bgpvrfs"}
	routePolicyGVR = schema.GroupVersionResource{Group: group, Version: version, Resource: "routepolicies"}
	prefixSetGVR   = schema.GroupVersionResource{Group: group, Version: version, Resource: "prefixsets"}
)

type FetchResult[T any] struct {
	Items             []T
	Count             int
	FetchDuration     time.Duration
	SerializeDuration time.Duration
}

type Client struct {
	dyn     dynamic.Interface
	Context string
}

func NewClient(kubeContext string) (*Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: kubeContext,
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules, configOverrides,
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("building kubeconfig for context %q: %w", kubeContext, err)
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating dynamic client: %w", err)
	}

	return &Client{dyn: dyn, Context: kubeContext}, nil
}

func (c *Client) ListVRFs(ctx context.Context) (*FetchResult[types.VRF], error) {
	fetchStart := time.Now()
	list, err := c.dyn.Resource(vrfGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	fetchDur := time.Since(fetchStart)
	if err != nil {
		return nil, fmt.Errorf("listing VRFs: %w", err)
	}

	serStart := time.Now()
	result := make([]types.VRF, 0, len(list.Items))
	for _, item := range list.Items {
		raw, err := item.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshalling VRF: %w", err)
		}
		var v types.VRF
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("unmarshalling VRF: %w", err)
		}
		result = append(result, v)
	}
	serDur := time.Since(serStart)

	return &FetchResult[types.VRF]{
		Items:             result,
		Count:             len(result),
		FetchDuration:     fetchDur,
		SerializeDuration: serDur,
	}, nil
}

func (c *Client) ListBGPVRFs(ctx context.Context) (*FetchResult[types.BGPVRF], error) {
	fetchStart := time.Now()
	list, err := c.dyn.Resource(bgpvrfGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	fetchDur := time.Since(fetchStart)
	if err != nil {
		return nil, fmt.Errorf("listing BGPVRFs: %w", err)
	}

	serStart := time.Now()
	result := make([]types.BGPVRF, 0, len(list.Items))
	for _, item := range list.Items {
		raw, err := item.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshalling BGPVRF: %w", err)
		}
		var v types.BGPVRF
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("unmarshalling BGPVRF: %w", err)
		}
		result = append(result, v)
	}
	serDur := time.Since(serStart)

	return &FetchResult[types.BGPVRF]{
		Items:             result,
		Count:             len(result),
		FetchDuration:     fetchDur,
		SerializeDuration: serDur,
	}, nil
}

func (c *Client) ListRoutePolicies(ctx context.Context) (*FetchResult[types.RoutePolicy], error) {
	fetchStart := time.Now()
	list, err := c.dyn.Resource(routePolicyGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	fetchDur := time.Since(fetchStart)
	if err != nil {
		return nil, fmt.Errorf("listing RoutePolicies: %w", err)
	}

	serStart := time.Now()
	result := make([]types.RoutePolicy, 0, len(list.Items))
	for _, item := range list.Items {
		raw, err := item.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshalling RoutePolicy: %w", err)
		}
		var v types.RoutePolicy
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("unmarshalling RoutePolicy: %w", err)
		}
		result = append(result, v)
	}
	serDur := time.Since(serStart)

	return &FetchResult[types.RoutePolicy]{
		Items:             result,
		Count:             len(result),
		FetchDuration:     fetchDur,
		SerializeDuration: serDur,
	}, nil
}

func (c *Client) ListPrefixSets(ctx context.Context) (*FetchResult[types.PrefixSet], error) {
	fetchStart := time.Now()
	list, err := c.dyn.Resource(prefixSetGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	fetchDur := time.Since(fetchStart)
	if err != nil {
		return nil, fmt.Errorf("listing PrefixSets: %w", err)
	}

	serStart := time.Now()
	result := make([]types.PrefixSet, 0, len(list.Items))
	for _, item := range list.Items {
		raw, err := item.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshalling PrefixSet: %w", err)
		}
		var v types.PrefixSet
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("unmarshalling PrefixSet: %w", err)
		}
		result = append(result, v)
	}
	serDur := time.Since(serStart)

	return &FetchResult[types.PrefixSet]{
		Items:             result,
		Count:             len(result),
		FetchDuration:     fetchDur,
		SerializeDuration: serDur,
	}, nil
}
