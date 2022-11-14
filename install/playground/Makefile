#-----------------------------------------------------------------------------
# Meshery Playground Production Environment
#-----------------------------------------------------------------------------

DEBUG=false
PROD_PLAYGROUND_PORT=80
REMOTE_PROVIDER=https://meshery.layer5.io

#-----------------------------------------------------------------------------
# Meshery Playground Production Deployment
#-----------------------------------------------------------------------------
.PHONY: play-prod-run 

## Deploys Meshery Playground in Prod
play-prod-run:
	(docker rm -f meshery) || true
	docker run layer5/meshery:playground-latest \
	--name meshery -d \
	-e PROVIDER_BASE_URLS=$(REMOTE_PROVIDER) \
	-e DEBUG=$(DEBUG) \
	-e PLAYGROUND=true \ 
	-p $(PROD_PLAYGROUND_PORT):8080 \
	