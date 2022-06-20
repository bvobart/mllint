ARG python_version=3.10
FROM python:${python_version}-alpine
WORKDIR /mllint

RUN apk update && apk add --no-cache gcc musl-dev python3-dev libffi-dev libgit2-dev

COPY build/requirements-tools.txt .
RUN pip install --no-cache-dir -r requirements-tools.txt

COPY mllint .

WORKDIR /app
VOLUME /app
ENTRYPOINT [ "/mllint/mllint" ]
