REGISTRY_NAMESPACE ?= nuancemobility
IMAGE_NAME := $(REGISTRY_NAMESPACE)/azure-request-limitometer
TAG ?= "dev"

build:
	@docker build --no-cache -t $(IMAGE_NAME):$(TAG) .

upload:
	@docker push $(IMAGE_NAME):$(TAG)

publish:
	@docker build -t $(IMAGE_NAME):$(TAG) .
	@docker push $(IMAGE_NAME):$(TAG)