all: clean build

clean:
	rm -rf bin/

build:
	gom install && gom build -o bin/shorty src/*
