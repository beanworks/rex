.PHONY: cluster

cluster: cluster-up

cluster-up:
	env COMPOSE_PROJECT_NAME=rex \
		COMPOSE_FILE=cluster/docker-compose.yml \
		docker-compose up -d --force-recreate
