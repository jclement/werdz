FROM node:lts-alpine as client-build-step
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY ./client  ./
RUN npm install
RUN npm run build

FROM golang as server-build-step
WORKDIR /build
COPY ./server ./
RUN go build -tags netgo -a -v .

FROM alpine
WORKDIR /app
EXPOSE 80
RUN mkdir data
COPY ./server/data /app/data
COPY --from=client-build-step /app/build /app/static
COPY --from=server-build-step /build/werdz /app
COPY ./server/werdz.yaml.sample /app/werdz.yaml
CMD ./werdz