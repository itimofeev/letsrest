

build-image:
	docker build -t letsrest/build --no-cache --force-rm=true -f ./tools/build/build.Dockerfile .

build: build-image
	docker run -v letsrest:/go/src -t letsrest/build make build
	docker cp $(shell docker ps -q -n=1):/letsrest ./tools/letsrest
	docker build -t letsrest --no-cache --force-rm=true -f ./tools/letsrest.Dockerfile tools
	rm ./tools/letsrest


