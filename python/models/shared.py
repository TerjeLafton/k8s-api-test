from __future__ import annotations

from datetime import datetime
from enum import Enum
from typing import Annotated, Any

from pydantic import BaseModel, ConfigDict, Field


class DeletionPolicy(str, Enum):
    delete = 'Delete'
    orphan = 'Orphan'


class Kind(str, Enum):
    prefix_set = 'PrefixSet'
    route_policy = 'RoutePolicy'
    vrf = 'VRF'
    bgpvrf = 'BGPVRF'


class Dependency(BaseModel):
    model_config = ConfigDict(
        populate_by_name=True,
    )
    kind: Annotated[Kind, Field(description='Kind of the dependent resource')]
    name: Annotated[str, Field(description='Name of the dependent resource')]
    namespace: Annotated[
        str | None,
        Field(
            description='Namespace of the dependent resource (optional, defaults to "default")'
        ),
    ] = None


class LifecyclePolicy(str, Enum):
    create = 'Create'
    update = 'Update'
    delete = 'Delete'


class ConditionStatus(str, Enum):
    true = 'True'
    false = 'False'
    unknown = 'Unknown'


class Condition(BaseModel):
    model_config = ConfigDict(
        populate_by_name=True,
    )
    last_transition_time: Annotated[
        datetime,
        Field(
            alias='lastTransitionTime',
            description='lastTransitionTime is the last time the condition transitioned from one status to another.',
        ),
    ]
    message: Annotated[
        str,
        Field(
            description='message is a human readable message indicating details about the transition.',
            max_length=32768,
        ),
    ]
    observed_generation: Annotated[
        int | None,
        Field(
            alias='observedGeneration',
            description='observedGeneration represents the .metadata.generation that the condition was set based upon.',
            ge=0,
        ),
    ] = None
    reason: Annotated[
        str,
        Field(
            description="reason contains a programmatic identifier indicating the reason for the condition's last transition.",
            max_length=1024,
            min_length=1,
            pattern=r'^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$',
        ),
    ]
    status: Annotated[
        ConditionStatus,
        Field(description='status of the condition, one of True, False, Unknown.'),
    ]
    type: Annotated[
        str,
        Field(
            description='type of condition in CamelCase or in foo.example.com/CamelCase.',
            max_length=316,
            pattern=r'^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$',
        ),
    ]


class Status(BaseModel):
    model_config = ConfigDict(
        populate_by_name=True,
    )
    conditions: Annotated[
        list[Condition] | None,
        Field(
            description='Conditions is a list of conditions that describe the current state of the resource.'
        ),
    ] = None
