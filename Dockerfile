#####################################
############# GO LANG BUILD #########
#####################################
FROM golang:latest AS BUILD
WORKDIR /src
ENV LAST_UPDATE=20190817
ADD go.mod .
ADD go.sum .
RUN go get -d
ADD . /src
# RUN go test --cover ./...
RUN go build -v -tags netgo -o jpb-mon-go
RUN ./jpb-mon-go --help
#####################################

#####################################
############## BUILD RUNTIME IMAGE ##
#####################################
FROM alpine:3.10
ENTRYPOINT ["/bin/jpb-mon"]
CMD ["--help"]
COPY letsencryptauthorityx3.pem /usr/src
COPY --from=BUILD /src/jpb-mon-go /bin/jpb-mon
#####################################
