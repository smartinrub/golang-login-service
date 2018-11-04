FROM golang:latest 
WORKDIR /go/src/golang-bcrypt
RUN go get -d -v golang.org/x/crypto/bcrypt
RUN go get -d -v github.com/go-sql-driver/mysql
RUN go get -d -v github.com/gorilla/mux
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV MYSQL_URL=localhost
ENV MYSQL_PORT=3306
ENV MYSQL_DB=credentials
ENV MYSQL_USER=newuser
ENV MYSQL_PASSWORD=password
WORKDIR /root/
# copies the first build into this stage
EXPOSE 8090
COPY --from=0 /go/src/golang-bcrypt/main .
CMD ["./main"]
