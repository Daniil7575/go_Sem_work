HOST = "localhost"
PORT = "5432"
USER = "postgres"
PASSWORD = "123"
NAME = "go_sem"
MAX_OVERFLOW = 15
POOL_SIZE = 15
URI = f"postgresql+asyncpg://{USER}:{PASSWORD}@{HOST}:{PORT}/{NAME}"
AUTH_MICROSERVICE_DOMAIN = "http://localhost:8001"
COOKIE_ACCESS_TOKEN_NAME = "blogaccesstoken"
DOMAIN = "http://localhost:8000"