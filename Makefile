VERSION=$(shell cat VERSION)

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t <registery>:latest .
	docker tag <registery>:latest <registery>:$(VERSION)

docker-push:
	docker push <registery>:latest
	docker push <registery>:$(VERSION)

docker-clean:
	docker rmi <registery>:latest || true
	docker rmi <registery>:$(VERSION) || true
	docker rm -v $(shell docker ps --filter status=exited -q 2>/dev/null) 2>/dev/null || true
	docker rmi $(shell docker images --filter dangling=true -q 2>/dev/null) 2>/dev/null || true

check-version-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi

update-tag:
	git pull --tags
	if git --no-pager tag --list | grep $(VERSION) -q ; then echo "$(VERSION) already exsits"; exit 1; fi
	git tag $(VERSION)
	git push origin $(VERSION)

unit-test:
	GO111MODULE=on go test -v ./...

test: unit-test

ci-build: check-version-tag docker-build docker-push docker-clean update-tag
