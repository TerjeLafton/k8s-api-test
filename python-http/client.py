from __future__ import annotations

import ssl
import time
from dataclasses import dataclass, field

import httpx
from kubernetes import config as k8s_config
from kubernetes.client import Configuration
from pydantic import BaseModel

from models import BGPVRF, VRF, PrefixSet, RoutePolicy

NAMESPACE = "default"
GROUP = "xr.crossnet.network.intility.io"
VERSION = "v1"

PLURAL_MAP: dict[str, str] = {
    "VRF": "vrfs",
    "BGPVRF": "bgpvrfs",
    "RoutePolicy": "routepolicies",
    "PrefixSet": "prefixsets",
}


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

        k8s_config.load_kube_config(context=kube_context)
        self._config = Configuration.get_default_copy()

        ssl_context = ssl.create_default_context()
        if self._config.ssl_ca_cert:
            ssl_context.load_verify_locations(self._config.ssl_ca_cert)
        if self._config.cert_file:
            ssl_context.load_cert_chain(
                certfile=self._config.cert_file,
                keyfile=self._config.key_file,
            )

        self._client = httpx.Client(
            base_url=self._config.host,
            verify=ssl_context,
            timeout=30.0,
        )

    def _get_auth_headers(self) -> dict[str, str]:
        token = self._config.get_api_key_with_prefix("authorization")
        if token:
            return {"Authorization": token}
        return {}

    def _list_resource[T: BaseModel](
        self,
        kind: str,
        model: type[T],
    ) -> FetchResult[T]:
        plural = PLURAL_MAP[kind]
        url = f"/apis/{GROUP}/{VERSION}/namespaces/{NAMESPACE}/{plural}"

        fetch_start = time.perf_counter()
        response = self._client.get(url, headers=self._get_auth_headers())
        response.raise_for_status()
        data = response.json()
        fetch_duration = time.perf_counter() - fetch_start

        ser_start = time.perf_counter()
        items = [model.model_validate(item) for item in data["items"]]
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
