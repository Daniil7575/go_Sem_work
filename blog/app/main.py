import aiohttp
from fastapi import FastAPI, Request, HTTPException, Response, Cookie, Form, status as http_status
from fastapi.responses import RedirectResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates


from database import async_session_factory
from repositories import PostgresPostsRepository, PostgresUserRepository
from schemas import SUserAuth, SPostCreate
import settings


app = FastAPI(
    title="Blog",
    description="A simple blog with auth microservice on go",
    docs_url="/docs",
    redoc_url=None,
)
app.mount("/static", StaticFiles(directory="static"), name="static")

templates = Jinja2Templates(directory="templates")


@app.post("/users/signup")
async def user_register(response: Response, username: str = Form(), password: str = Form()):
    async with aiohttp.ClientSession() as session:
        async with session.post(
            settings.AUTH_MICROSERVICE_DOMAIN + "/signup",
            json={"username": username, "password": password},
        ) as resp:
            text = await resp.text()
            if (status := resp.status) != 200:
                raise HTTPException(status)
    response.set_cookie(key=settings.COOKIE_ACCESS_TOKEN_NAME, value=text)
    return RedirectResponse("/", status_code=http_status.HTTP_302_FOUND, headers=response.headers)

@app.get("/users/signup")
async def user_register_html(request: Request):
    return templates.TemplateResponse("signup.html", context={"request": request})


@app.post("/users/signin")
async def user_login(
    response: Response, username: str = Form(), password: str = Form()
):
    async with aiohttp.ClientSession() as session:
        async with session.post(
            settings.AUTH_MICROSERVICE_DOMAIN + "/signin",
            json={"username": username, "password": password},
        ) as resp:
            text = await resp.text()
            if (status := resp.status) != 200:
                raise HTTPException(status)
    response.set_cookie(key=settings.COOKIE_ACCESS_TOKEN_NAME, value=text)
    return RedirectResponse("/", status_code=http_status.HTTP_302_FOUND, headers=response.headers)

@app.get("/users/signin")
async def user_login_html(request: Request):
    return templates.TemplateResponse("signin.html", context={"request": request})


async def check_token(token: str | None, silent: bool = False) -> str:
    if not token:
        token = ""
    async with aiohttp.ClientSession() as session:
        print(token)
        async with session.post(
            settings.AUTH_MICROSERVICE_DOMAIN + "/access-token",
            json={"access_token": token},
        ) as resp:
            text = await resp.text()
            if (status := resp.status) != 200:
                if silent:
                    return ""
                raise HTTPException(status)
    return text


@app.post("/posts")
async def create_post(
    title: str = Form(),
    description: str = Form(),
    blogaccesstoken: str | None = Cookie(default=None),
):
    author_id = await check_token(blogaccesstoken)
    async with async_session_factory() as session:
        posts_repo = PostgresPostsRepository(session)
        new_post = await posts_repo.create_post(
            {"title": title, "description": description, "author_id": author_id}
        )
    return RedirectResponse("/", status_code=302)


@app.get("/posts")
async def create_post(
    request: Request,
    blogaccesstoken: str | None = Cookie(default=None),
):
    return templates.TemplateResponse("create_post.html", context={"request": request})


@app.get("/")
async def main_page(
    request: Request,
    blogaccesstoken: str | None = Cookie(default=None),
):
    user_id = await check_token(blogaccesstoken, silent=True)
    username = None
    async with async_session_factory() as session:
        if user_id:
            user_repo = PostgresUserRepository(session)
            username = await user_repo.get_by_id_or_none(user_id)
        posts_repo = PostgresPostsRepository(session)
        posts = await posts_repo.list_all()
    return templates.TemplateResponse(
        "base.html", context={"posts": posts, "request": request, "user_id": str(user_id)}
    )
