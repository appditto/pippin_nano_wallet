FROM python:3.7

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --trusted-host pypi.org --no-cache-dir pippin-wallet==1.1.3
# Temporary workaround for https://github.com/tortoise/tortoise-orm/issues/359
RUN pip install --trusted-host pypi.org --no-cache-dir -U tortoise-orm==0.15.5

RUN mkdir PippinData
COPY docker.config.yaml PippinData/config.yaml

CMD ["pippin-server"]