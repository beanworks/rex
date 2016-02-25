.PHONY: cluster

COMPOSE_ENV = env COMPOSE_PROJECT_NAME=rex \
				  COMPOSE_FILE=cluster/docker-compose.yml

RABBIT_VHOST = rex
RABBIT_USER = rex
RABBIT_PASS = pwd

RABBIT_1_CTL = docker exec rex_rabbitmq1_1 rabbitmqctl
RABBIT_2_CTL = docker exec rex_rabbitmq2_1 rabbitmqctl
RABBIT_3_CTL = docker exec rex_rabbitmq3_1 rabbitmqctl

cluster: cluster-build cluster-up cluster-config

cluster-build:
	$(COMPOSE_ENV) docker-compose build

cluster-up:
	$(COMPOSE_ENV) docker-compose up -d --force-recreate

cluster-config: rabbit-join rabbit-user rabbit-perm rabbit-ha

rabbit-join:
	# Reset rabbitmq1
	$(RABBIT_1_CTL) stop_app
	$(RABBIT_1_CTL) reset
	$(RABBIT_1_CTL) start_app
	# Reset rabbitmq2 and join cluster with rabbitmq1
	$(RABBIT_2_CTL) stop_app
	$(RABBIT_2_CTL) reset
	$(RABBIT_2_CTL) join_cluster rabbit@rabbitmq1
	$(RABBIT_2_CTL) start_app
	# Reset rabbitmq3 and join cluster with rabbitmq1
	$(RABBIT_3_CTL) stop_app
	$(RABBIT_3_CTL) reset
	$(RABBIT_3_CTL) join_cluster rabbit@rabbitmq1
	$(RABBIT_3_CTL) start_app

rabbit-user:
	$(RABBIT_1_CTL) add_user $(RABBIT_USER) $(RABBIT_PASS)
	$(RABBIT_1_CTL) set_user_tags $(RABBIT_USER) administrator

rabbit-perm:
	$(RABBIT_1_CTL) add_vhost $(RABBIT_VHOST)
	$(RABBIT_1_CTL) set_permissions -p $(RABBIT_VHOST) $(RABBIT_USER) ".*" ".*" ".*"

rabbit-ha:
	$(RABBIT_1_CTL) set_policy -p $(RABBIT_VHOST) rex-ha-all "^rex\.ha\." '{"ha-mode":"all"}'
