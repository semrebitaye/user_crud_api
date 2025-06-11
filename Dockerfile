# use official golang image
FROM golang:1.23-alpine

# set working directory
WORKDIR /app

# copy the csource code
COPY . .

# download and install all dependencies
RUN go get -d -v ./...

# build the go app
RUN go build -o user_crud_api .

# expose with port number
EXPOSE 8080

# run the excutable
CMD [ "./user_crud_api" ]



