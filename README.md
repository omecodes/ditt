# Data Impact Technical Test

Solution to the technical test

## Install and build

```bash
git clone https://github.com/omecodes/ditt
cd ditt

go get -v -t -d ./...
cd bin

make

chmod u+x ditt-api-server
```

## Run

Assuming you are in the `bin` folder, run the following command:

```
./ditt-api-server start --port=8080
```

### Database target

By default the program target a mongo database running at localhost. If you want it to connect to another mongo database
you can specify it this way:

```
./ditt-api-server start --port=8080 --db-uri=<target-db-uri>
```

### Authentication

When starting the app, it displays the path of the configuration directory. Inside that directory there is a file
named `admin-auth` that contains the admin password. Please use this password to login as admin and add users.

## Comments

### Testing

The API is pretty well covered. Tests for implementations of `UserDataStore` and `Files` are the big absents but this is
due to lack of time.

### Specification

Some part of the specification were not very clear to me. Maybe it'll be discussed during our next call

### Difficulties

I spent time reading about routine/channels and Mongo. About Mongo, I tried up to 2 different go implementations. 