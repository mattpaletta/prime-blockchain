FROM alpine:latest
RUN apk add yarn --no-cache

RUN yarn global add serve
COPY ./package.json /prime/package.json
WORKDIR /prime
RUN yarn install
COPY ./ /prime
EXPOSE 4000
RUN yarn --version && node --version && yarn run build
ENTRYPOINT ["serve", "-s", "build", "-p", "4000"]
