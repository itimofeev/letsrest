

clean:
	rm -rf tools/target/*

build-image:
	docker build -t letsrest/build --no-cache --force-rm=true -f ./tools/build/build.Dockerfile .

cp-configs:
	cp ./tools/prod-files/* ./tools/target/

cp-makefile:
	cp Makefile ./tools/target

build-letsrest: build-image
	docker run -v letsrest:/go/src -t letsrest/build make build
	docker cp $(shell docker ps -q -n=1):/letsrest ./tools/letsrest
	docker build -t letsrest --no-cache --force-rm=true -f ./tools/letsrest.Dockerfile tools
	rm ./tools/letsrest
	docker save -o ./tools/target/letsrest.img letsrest

build-frontend:
	cp -r /Users/ilyatimofee/prog/js/letsrest-ui/build ./tools/target/frontend
	cd ./tools/target ; tar -jcvf frontend.tar.bz2 frontend
	rm -r ./tools/target/frontend

build: clean build-letsrest cp-configs cp-makefile build-frontend

upload:
	rsync -ravezP ./tools/target/* -e ssh ilyaufo@188.166.26.165:/home/ilyaufo/letsrest


prepare-run:
	docker load -i letsrest.img
	tar -jxvf frontend.tar.bz2
	chown ilyaufo frontend

run: prepare-run
	docker-compose -p letsrest -f prod.docker-compose.yml up -d --build

stop:
	docker-compose -p letsrest -f prod.docker-compose.yml stop

# =====================

run-mongo:
	docker run -p 27017:27017 --name mg-letsrest -d mongo:3.5.9
