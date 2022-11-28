export KO_DOCKER_REPO = kickbeak

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ cluster
.PHONY: create
create: ## Create dev cluster
	kind create cluster --name k8s-notion --config=kind-config.yaml || true 
.PHONY: delete
delete: ## Create dev cluster
	kind delete	cluster --name k8s-notion

##@ deploy

.PHONY: deploy
deploy: ## Deploy to cluster using ko
	ko apply -f k8s/

.PHONY: undeploy
undeploy: ## Undeploy from cluster using ko
	ko delete -f k8s/
