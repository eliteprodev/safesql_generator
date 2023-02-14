# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.17.0
# source: query.sql
from typing import Optional

import sqlalchemy
import sqlalchemy.ext.asyncio

from querytest import models


DELETE_BAR_BY_ID = """-- name: delete_bar_by_id \\:one
DELETE FROM bars WHERE id = :p1 RETURNING id, name
"""


DELETE_EXCLUSION_BY_ID = """-- name: delete_exclusion_by_id \\:one
DELETE FROM exclusions WHERE id = :p1 RETURNING id, name
"""


DELETE_MY_DATA_BY_ID = """-- name: delete_my_data_by_id \\:one
DELETE FROM my_data WHERE id = :p1 RETURNING id, name
"""


class Querier:
    def __init__(self, conn: sqlalchemy.engine.Connection):
        self._conn = conn

    def delete_bar_by_id(self, *, id: int) -> Optional[models.Bar]:
        row = self._conn.execute(sqlalchemy.text(DELETE_BAR_BY_ID), {"p1": id}).first()
        if row is None:
            return None
        return models.Bar(
            id=row[0],
            name=row[1],
        )

    def delete_exclusion_by_id(self, *, id: int) -> Optional[models.Exclusions]:
        row = self._conn.execute(sqlalchemy.text(DELETE_EXCLUSION_BY_ID), {"p1": id}).first()
        if row is None:
            return None
        return models.Exclusions(
            id=row[0],
            name=row[1],
        )

    def delete_my_data_by_id(self, *, id: int) -> Optional[models.MyData]:
        row = self._conn.execute(sqlalchemy.text(DELETE_MY_DATA_BY_ID), {"p1": id}).first()
        if row is None:
            return None
        return models.MyData(
            id=row[0],
            name=row[1],
        )


class AsyncQuerier:
    def __init__(self, conn: sqlalchemy.ext.asyncio.AsyncConnection):
        self._conn = conn

    async def delete_bar_by_id(self, *, id: int) -> Optional[models.Bar]:
        row = (await self._conn.execute(sqlalchemy.text(DELETE_BAR_BY_ID), {"p1": id})).first()
        if row is None:
            return None
        return models.Bar(
            id=row[0],
            name=row[1],
        )

    async def delete_exclusion_by_id(self, *, id: int) -> Optional[models.Exclusions]:
        row = (await self._conn.execute(sqlalchemy.text(DELETE_EXCLUSION_BY_ID), {"p1": id})).first()
        if row is None:
            return None
        return models.Exclusions(
            id=row[0],
            name=row[1],
        )

    async def delete_my_data_by_id(self, *, id: int) -> Optional[models.MyData]:
        row = (await self._conn.execute(sqlalchemy.text(DELETE_MY_DATA_BY_ID), {"p1": id})).first()
        if row is None:
            return None
        return models.MyData(
            id=row[0],
            name=row[1],
        )
