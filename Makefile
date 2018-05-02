build:
	docker build -t uphy/drone-image-copy .
run:
	docker run -it --rm --privileged -e "PLUGIN_REGISTRY=registry:5000" -e "PLUGIN_IMAGES=bash,hello-world" uphy/drone-image-copy
