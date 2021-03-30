FROM python:3.9-slim-buster

RUN apt-get update \
 && apt-get install -y build-essential \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

COPY . .

RUN pip install --no-cache-dir .

RUN mkdir PippinData
COPY docker.config.yaml PippinData/config.yaml

CMD ["pippin-server"]