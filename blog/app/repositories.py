from typing import Any
import uuid
from sqlalchemy.ext.asyncio import AsyncSession
from models import PostDTO, UserDTO
import sqlalchemy as sa


class PostgresPostsRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def create_post(self, new_post_data: dict[str, Any]) -> PostDTO:
        new_post = PostDTO(**new_post_data)
        self.session.add(new_post)
        await self.session.commit()
        return new_post

    async def list_all(self) -> list[PostDTO]:
        stmt = sa.Select(PostDTO)
        posts = (await self.session.execute(stmt)).scalars().all()
        return posts


class PostgresUserRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_id_or_none(self, user_id: uuid.UUID | str) -> UserDTO | None:
        stmt = sa.select(UserDTO).filter_by(id=user_id)
        user = (await self.session.execute(stmt)).scalars().one_or_none()
        return user
