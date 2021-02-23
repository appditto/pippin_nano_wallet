FROM python:3.8

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --trusted-host pypi.org --no-cache-dir pippin-wallet==1.1.18

RUN mkdir PippinData
COPY docker.config.yaml PippinData/config.yaml

CMD ["pippin-server"]