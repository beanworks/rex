.PHONY: cluster

cluster: cluster-up cluster-join

cluster-up:
	env COMPOSE_PROJECT_NAME=rex \
		COMPOSE_FILE=cluster/docker-compose.yml \
		docker-compose up -d --force-recreate

cluster-join:
	docker exec -it rex_rabbitmq1_1 sh -c 'rabbitmqctl stop_app'
	docker exec -it rex_rabbitmq1_1 sh -c 'rabbitmqctl reset'
	docker exec -it rex_rabbitmq1_1 sh -c 'rabbitmqctl start_app'
	docker exec -it rex_rabbitmq2_1 sh -c 'rabbitmqctl stop_app'
	docker exec -it rex_rabbitmq2_1 sh -c 'rabbitmqctl reset'
	docker exec -it rex_rabbitmq2_1 sh -c 'rabbitmqctl join_cluster rabbit@rabbitmq1'
	docker exec -it rex_rabbitmq2_1 sh -c 'rabbitmqctl start_app'
	docker exec -it rex_rabbitmq3_1 sh -c 'rabbitmqctl stop_app'
	docker exec -it rex_rabbitmq3_1 sh -c 'rabbitmqctl reset'
	docker exec -it rex_rabbitmq3_1 sh -c 'rabbitmqctl join_cluster rabbit@rabbitmq1'
	docker exec -it rex_rabbitmq3_1 sh -c 'rabbitmqctl start_app'
