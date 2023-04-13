FROM golang:1.19-buster AS build

WORKDIR /src

RUN export GOPRIVATE="bitbucket.org/ntuclink"
ARG BITBUCKET_CREDENTIAL
RUN git config --global url."https://${BITBUCKET_CREDENTIAL}@bitbucket.org/".insteadOf "https://bitbucket.org/" \
	&& git config --global user.email "developers@ntucenterprise.sg" \
	&& git config --global user.name "developers"
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go install ./cmd/api


FROM debian:buster
RUN apt-get update && apt-get install -y ca-certificates tzdata
COPY --from=build /go/bin/* /bin/

ARG CLIENT_CERT
ARG CLIENT_KEY
ARG SERVER_CA

RUN echo ${CLIENT_CERT} | base64 -d > /client-cert.pem
RUN echo ${CLIENT_KEY} | base64 -d > /client-key.pem
RUN echo ${SERVER_CA} | base64 -d > /server-ca.pem

CMD ["/bin/api"]