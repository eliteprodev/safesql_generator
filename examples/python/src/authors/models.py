# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.17.0
import dataclasses
from typing import Optional


@dataclasses.dataclass()
class Author:
    id: int
    name: str
    bio: Optional[str]
