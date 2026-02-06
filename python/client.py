from __future__ import annotations

import time
from dataclasses import dataclass, field

from kubernetes import config as k8s_config
from kubernetes.dynamic import DynamicClient
from pydantic import BaseModel

from models import BGPVRF, VRF, PrefixSet, RoutePolicy

NAMESPACE = "default"
GROUP = "xr.crossnet.network.intility.io"
VERSION = "v1"


@dataclass
class FetchResult[T]:
    items: list[T] = field(default_factory=list)
    fetch_duration: float = 0.0
    serialize_duration: float = 0.0

    @property
    def count(self) -> int:
        return len(self.items)


class Client:
    def __init__(self, kube_context: str) -> None:
        self.context = kube_context
        api_client = k8s_config.new_client_from_config(context=kube_context)
        self._dyn = DynamicClient(api_client)

    def _list_resource[T: BaseModel](
        self,
        kind: str,
        model: type[T],
    ) -> FetchResult[T]:
        api = self._dyn.resources.get(
            api_version=f"{GROUP}/{VERSION}",
            kind=kind,
        )

        fetch_start = time.perf_counter()
        result = api.get(namespace=NAMESPACE)
        fetch_duration = time.perf_counter() - fetch_start

        ser_start = time.perf_counter()
        items = [model.model_validate(item.to_dict()) for item in result.items]
        serialize_duration = time.perf_counter() - ser_start

        return FetchResult(
            items=items,
            fetch_duration=fetch_duration,
            serialize_duration=serialize_duration,
        )

    def list_vrfs(self) -> FetchResult[VRF]:
        return self._list_resource("VRF", VRF)

    def list_bgpvrfs(self) -> FetchResult[BGPVRF]:
        return self._list_resource("BGPVRF", BGPVRF)

    def list_route_policies(self) -> FetchResult[RoutePolicy]:
        return self._list_resource("RoutePolicy", RoutePolicy)

    def list_prefix_sets(self) -> FetchResult[PrefixSet]:
        return self._list_resource("PrefixSet", PrefixSet)
