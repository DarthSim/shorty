# Shorty
[![Build Status](https://travis-ci.org/DarthSim/shorty.svg)](https://travis-ci.org/DarthSim/shorty)

Tiny url shortener

## Installation
You will need Go and [Gom](https://github.com/mattn/gom) to build the project and PostgreSQL to make Shorty fly.

```bash
make

# First launch
RESET_DB=1 bin/shorty

# Futher launches
bin/shorty
```

#### Configuration

You can specify DB connection string by setting DB_CONN variable:

```bash
DB_CONN="dbname=my_db sslmode=disable" bin/shorty
```

You can specify server address by setting ADDRESS variable:

```bash
ADDRESS="192.168.1.1:4321" bin/shorty
```

## API

````
POST /shorten (url=http://url_to_short.com/?lorem=ipsum)
# => http://domain.com/:code

GET /expand/:code
# => http://url_to_short.com/?lorem=ipsum

GET /:code
# => Redirect to http://url_to_short.com/?lorem=ipsum

GET /statistics/:code
# => Count of redirects to http://url_to_short.com/?lorem=ipsum
````

## How to run tests
1. First of all you need a public PostgreSQL DB named `shorty_test`.

2. Next install testing packages with
```bash
gom -test install
```

3. And finally run the following
```bash
gom test src/*
```
