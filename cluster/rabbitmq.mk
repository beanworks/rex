.PHONY: cluster

COMPOSE_ENV = env COMPOSE_PROJECT_NAME=rex \
				  COMPOSE_FILE=cluster/docker-compose.yml

RABBIT_VHOST = rex
RABBIT_USER = rex
RABBIT_USER_PASS = $(RABBIT_USER) passw0rd

cluster: cluster-build cluster-up cluster-config

cluster-build:
	$(COMPOSE_ENV) docker-compose build

cluster-up:
	$(COMPOSE_ENV) docker-compose up -d --force-recreate

cluster-config: rabbit-join rabbit-user rabbit-perm rabbit-ha

rabbit-join:
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

rabbit-user:
	docker exec rex_rabbitmq1_1 rabbitmqctl add_user $(RABBIT_USER_PASS)

rabbit-perm:
	docker exec rex_rabbitmq1_1 rabbitmqctl add_vhost $(RABBIT_VHOST)

rabbit-ha:
	echo