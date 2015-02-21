# Shorty
[![Build Status](https://travis-ci.org/DarthSim/shorty.svg)](https://travis-ci.org/DarthSim/shorty)

Tiny url shortener

## Installation
You will need Go and [Gom](https://github.com/mattn/gom) to build the project and PostgreSQL to make Shorty fly.

```bash
make
cp shorty.conf.sample shorty.conf

# First launch
bin/shorty --init-db-schema

# Futher launches
bin/shorty
```

#### Configuration

You can specify the path to the config file using `--config` key:

```bash
bin/shorty --config /etc/shorty/shorty.conf
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
1. First of all you need a PostgreSQL DB named `shorty_test`.

2. Next install testing packages with
```bash
gom -test install
```

3. And finally run the following
```bash
DBUSER=your_db_user DBPASS=your_db_user_password gom test src/*
```
