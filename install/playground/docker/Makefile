#-----------------------------------------------------------------------------
# Meshery Playground Production Environment
#-----------------------------------------------------------------------------

DEBUG=true
PROD_PLAYGROUND_PORT=80
# Despite the name, REMOTE_PROVIDER here holds the PROVIDER_BASE_URLS list
# (comma-separated remote-provider endpoints), NOT a provider name. The
# provider *name* in install/Makefile.core.mk is the empty string by default.
REMOTE_PROVIDER=https://cloud.meshery.io

#-----------------------------------------------------------------------------
# Meshery Playground Production Deployment
#-----------------------------------------------------------------------------
.PHONY: play-prod-run play-prod-stop 

## Deploys Meshery Playground in Prod
play-prod-run:
	(docker rm -f meshery) || true
	docker run -d --name meshery \
	-e PROVIDER_BASE_URLS=$(REMOTE_PROVIDER) \
	-e DEBUG=$(DEBUG) \
	-e PLAYGROUND=true \
	-e MESHERY_SERVER_CALLBACK_URL=https://playground.meshery.io/api/user/token \ 
	-p $(PROD_PLAYGROUND_PORT):8080 \
	meshery/meshery:playground-latest ./meshery
	
# docker run -d --name meshery \
# 	-e PROVIDER_BASE_URLS=$(REMOTE_PROVIDER) \
# 	-e PLAYGROUND=true \
# 	-p 80:8080 \
# 	meshery/meshery:playground-latest ./meshery

play-prod-stop:
	docker rm -f meshery
