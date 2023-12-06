from datetime import datetime
import uuid
from sqlalchemy.orm import Mapped, mapped_column
import sqlalchemy as sa
from database import Base


class UserDTO(Base):
    __tablename__ = "users"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True)
    username: Mapped[str] = mapped_column(unique=True)
    password_hash: Mapped[str]


class PostDTO(Base):
    __tablename__ = "posts"

    id: Mapped[uuid.UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    title: Mapped[str] = mapped_column(sa.String(200))
    description: Mapped[str]
    author_id: Mapped[uuid.UUID] = mapped_column(sa.ForeignKey("users.id"))
    created_at: Mapped[datetime] = mapped_column(server_default=sa.func.now())
