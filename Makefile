.PHONY: cluster

cluster: cluster-up cluster-join

cluster-up:
	env COMPOSE_PROJECT_NAME=rex \
		COMPOSE_FILE=cluster/docker-compose.yml \
		docker-compose up -d --force-recreate

cluster-join:
	# Reset rabbitmq1
	docker exec rex_rabbitmq1_1 rabbitmqctl stop_app
	docker exec rex_rabbitmq1_1 rabbitmqctl reset
	docker exec rex_rabbitmq1_1 rabbitmqctl start_app
	# Reset rabbitmq2 and join cluster with rabbitmq1
	docker exec rex_rabbitmq2_1 rabbitmqctl stop_app
	docker exec rex_rabbitmq2_1 rabbitmqctl reset
	docker exec rex_rabbitmq2_1 rabbitmqctl join_cluster rabbit@rabbitmq1
	docker exec rex_rabbitmq2_1 rabbitmqctl start_app
	# Reset rabbitmq3 and join cluster with rabbitmq1
	docker exec rex_rabbitmq3_1 rabbitmqctl stop_app
	docker exec rex_rabbitmq3_1 rabbitmqctl reset
	docker exec rex_rabbitmq3_1 rabbitmqctl join_cluster rabbit@rabbitmq1
	docker exec rex_rabbitmq3_1 rabbitmqctl start_app
