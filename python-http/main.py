from __future__ import annotations

import time
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass, field

from client import Client, FetchResult

ITERATIONS = 10

CONTEXTS = [
    "default/api-crossnet-z2-staging-1-m5yc9s-apps-intilitycloud-com:443/terje.lafton@intility.no",
    "default/api-crossnet-z2-staging-2-u5fc1m-apps-intilitycloud-com:443/terje.lafton@intility.no",
]

RESOURCE_METHODS = ["list_vrfs", "list_bgpvrfs", "list_route_policies", "list_prefix_sets"]
RESOURCE_LABELS = ["VRFs", "BGPVRFs", "RoutePolicies", "PrefixSets"]


@dataclass
class ResourceStats:
    count: int = 0
    total_fetch: float = 0.0
    total_serialize: float = 0.0

    def avg_fetch_ms(self) -> float:
        return (self.total_fetch / ITERATIONS) * 1000

    def avg_serialize_ms(self) -> float:
        return (self.total_serialize / ITERATIONS) * 1000


@dataclass
class ClusterStats:
    context: str
    resources: dict[str, ResourceStats] = field(default_factory=dict)

    def __post_init__(self) -> None:
        if not self.resources:
            self.resources = {name: ResourceStats() for name in RESOURCE_METHODS}


def add_result(stats: ResourceStats, result: FetchResult) -> None:
    stats.count = result.count
    stats.total_fetch += result.fetch_duration
    stats.total_serialize += result.serialize_duration


def fetch_cluster(client: Client) -> dict[str, FetchResult]:
    results: dict[str, FetchResult] = {}
    with ThreadPoolExecutor(max_workers=len(RESOURCE_METHODS)) as pool:
        futures = {
            pool.submit(getattr(client, method)): method
            for method in RESOURCE_METHODS
        }
        for future in futures:
            results[futures[future]] = future.result()
    return results


def print_stats(label: str, s: ResourceStats) -> None:
    print(
        f"  {label:<14} {s.count:4d} items | "
        f"avg fetch {s.avg_fetch_ms():8.0f}ms | "
        f"avg serialize {s.avg_serialize_ms():8.3f}ms"
    )


def main() -> None:
    clients = [Client(ctx) for ctx in CONTEXTS]
    stats = [ClusterStats(context=c.context) for c in clients]

    for iteration in range(ITERATIONS):
        run_start = time.perf_counter()

        with ThreadPoolExecutor(max_workers=len(clients)) as pool:
            cluster_results = list(pool.map(fetch_cluster, clients))

        elapsed = (time.perf_counter() - run_start) * 1000
        print(f"Run {iteration + 1}/{ITERATIONS}: {elapsed:.0f}ms")

        for cs, result_map in zip(stats, cluster_results):
            for method, result in result_map.items():
                add_result(cs.resources[method], result)

    print(f"\nAverages over {ITERATIONS} runs:\n")
    for cs in stats:
        print(f"[{cs.context}]")
        for method, label in zip(RESOURCE_METHODS, RESOURCE_LABELS):
            print_stats(label, cs.resources[method])
        print()


if __name__ == "__main__":
    main()
