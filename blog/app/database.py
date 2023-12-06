from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import DeclarativeBase, sessionmaker

import settings


class Base(DeclarativeBase):
    pass


engine = create_async_engine(
    settings.URI,
    pool_size=settings.POOL_SIZE,
    max_overflow=settings.MAX_OVERFLOW,
    pool_pre_ping=True,
)

async_session_factory = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
    autocommit=False,
    autoflush=False,
)

from models import *  # pylint: disable=C0413  # isort:skip  # noqa: F403, E402
