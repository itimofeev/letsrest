

clean:
	rm -r tools/target/*

build-image:
	docker build -t letsrest/build --no-cache --force-rm=true -f ./tools/build/build.Dockerfile .

cp-configs:
	cp ./tools/prod-files/* ./tools/target/

build-letsrest: build-image
	docker run -v letsrest:/go/src -t letsrest/build make build
	docker cp $(shell docker ps -q -n=1):/letsrest ./tools/letsrest
	docker build -t letsrest --no-cache --force-rm=true -f ./tools/letsrest.Dockerfile tools
	rm ./tools/letsrest
	docker save -o ./tools/target/letsrest.img letsrest

build-frontend:
	cp -r /Users/ilyatimofee/prog/js/letsrest-ui/build ./tools/target/frontend
	tar -jcvf ./tools/target/frontend.tar.bz2 ./tools/target/frontend
	rm -r ./tools/target/frontend

build: clean build-letsrest cp-configs build-frontend

deploy:
	scp -r tools/target/* ilyaufo@188.166.26.165:/home/ilyaufo/letsrest