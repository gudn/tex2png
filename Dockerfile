FROM golang:1.17 as builder

WORKDIR /app
COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

COPY ./tex2png.go tex2png.go
RUN go build -o /tex2png

FROM ubuntu:20.04

RUN apt update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install tzdata

RUN apt install -y texlive texlive-lang-cyrillic texlive-latex-extra poppler-utils &&\
  apt-get clean && rm -rf /var/cache/apt/*

WORKDIR /app

COPY ./templates /app/templates
COPY --from=builder /tex2png /app/tex2png

RUN useradd -U -M t2p
RUN chown -R t2p:t2p /app

USER t2p
ENV T2P_PORT=8080
EXPOSE 8080

CMD '/app/tex2png'
