FROM python:3.7

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --trusted-host pypi.org --no-cache-dir pippin-wallet==1.1.17

RUN mkdir PippinData
COPY docker.config.yaml PippinData/config.yaml

CMD ["pippin-server"]