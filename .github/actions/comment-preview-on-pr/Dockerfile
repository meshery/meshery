FROM python:3.7@sha256:eedf63967cdb57d8214db38ce21f105003ed4e4d0358f02bedc057341bcf92a0

RUN pip install httpx "pydantic==1.5.1" pygithub

COPY ./app /app

CMD ["python", "/app/main.py"]