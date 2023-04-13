go = go

cr_user = k8scat
cr_server = ghcr.io
cr_password =

image_version = latest
image_name = resium-downhub
image_dockerfile = Dockerfile

image_tag = $(cr_server)/$(cr_user)/$(image_name):$(image_version)
image_tag_latest = $(cr_server)/$(cr_user)/$(image_name):latest

.PHONY: run-dev
run-dev:
	$(go) run main.go

.PHONY: build
build:
	$(go) build -trimpath -o bin/downhub main.go

.PHONY: build-image
build-image:
	docker build \
		--no-cache \
		-f Dockerfile \
		-t $(image_tag) .
	docker tag $(image_tag) $(image_tag_latest)

.PHONY: push-image
push-image:
	docker push $(image_tag)
	docker push $(image_tag_latest)

.PHONY: login-cr
login-cr:
	docker login -u $(cr_user) -p $(cr_password) $(cr_server)

.PHONY: login-cr
logout-cr:
	docker logout $(cr_server)
