
dev:
	docker-compose -f compose.yml up --build -d
	-docker exec -it mpsse-dev /bin/bash
	docker-compose -f compose.yml kill

docker:
	docker build -f build.Dockerfile -t vaporio/libmpsse-base .

build:
	cd src ; ./configure --disable-python
	cd src ; make
	cd src ; make install
	go build

lint:
	golint .