FROM alpine as builder
RUN apk --update add curl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /bin/kubectl

FROM alpine

RUN apk --update add ca-certificates

COPY ./build/lessor-linux-amd64 /bin/lessor
COPY --from=builder /bin/kubectl /bin/kubectl

CMD ["lessor"]
