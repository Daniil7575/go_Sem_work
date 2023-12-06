import uuid
from pydantic import BaseModel


class SUserAuth(BaseModel):
    username: str
    password: str


class SAccessToken(BaseModel):
    value: str


class SPostCreate(BaseModel):
    title: str
    description: str
