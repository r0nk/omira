FROM golang

COPY . omira/

WORKDIR omira/

RUN go build .


CMD [ "/bin/bash" ]

