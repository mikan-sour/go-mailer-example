current_dir = $(shell pwd)
main_path = ./src/main.go
kafka_network = kafka-network 
kafka_host_port = localhost:9092
kafka_topic = errors

# Add coloration
ifneq (,$(findstring xterm,${TERM}))
	RED          := $(shell tput -Txterm setaf 1)
	GREEN        := $(shell tput -Txterm setaf 2)
	LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
	PURPLE       := $(shell tput -Txterm setaf 5)
	YELLOW       := $(shell tput -Txterm setaf 3)
	BLUE         := $(shell tput -Txterm setaf 6)
	RESET := $(shell tput -Txterm sgr0)
else
	RED          := ""
	GREEN        := ""
	LIGHTPURPLE  := ""
	PURPLE       := ""
	YELLOW 		 := ""
	BLUE         := ""
	RESET        := ""
endif

# set target color

POUND = \#

env:
	@echo "${PURPLE}Setup .env ...${RESET}"
	cp ./.env.sample ./.env
	@echo "${PURPLE}.Env setup!...${RESET}"

kafka:
	@echo "${BLUE"}Making network...${RESET}"
	docker network create -d bridge $(kafka_network)
	@echo "${BLUE}Setup Zookeeper...${RESET}"
	docker run -d -p 2181:2181 \
	-e "ALLOW_ANONYMOUS_LOGIN=yes" \
	--network $(kafka_network) \
	--name zookeeper bitnami/zookeeper:3.8
	@echo "${BLUE}Kafka setup!${RESET}"

	@echo "${BLUE}Setup Kafka...${RESET}"
	docker run -d -p 9092:9092 \
	-v $(current_dir)/kafka_data:/bitnami/kafka \
	-e "KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181" \
	-e "ALLOW_PLAINTEXT_LISTENER=yes" \
	-e "KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CLIENT:PLAINTEXT,EXTERNAL:PLAINTEXT" \
	-e "KAFKA_CFG_LISTENERS=CLIENT://:9092,EXTERNAL://:9093" \
	-e "KAFKA_CFG_ADVERTISED_LISTENERS=CLIENT://$(kafka_host_port)" \
	-e "KAFKA_CFG_INTER_BROKER_LISTENER_NAME=CLIENT" \
	--network $(kafka_network) \
	--name kafka bitnami/kafka:3.2
	@echo "${BLUE}Kafka setup!${RESET}"

mailhog:
	@echo "${GREEN}Setup mailhog...${RESET}"
	docker run -d -e "MH_STORAGE=maildir" \
	-v $(current_dir)/maildir:/maildir \
	-p 1025:1025 -p 8025:8025 --name mailhog mailhog/mailhog
	@echo "${GREEN}Mailhog setup!${RESET}"

clean:
	@echo "${RED}Cleaning...${RESET}"
	docker rm -v --force $(shell docker ps -a -q -f name=mailhog)
	rm -rf maildir
	docker rm -v --force $(shell docker ps -a -q -f name=kafka)
	rm -rf kafka_data
	docker rm -v --force $(shell docker ps -a -q -f name=zookeeper)
	docker network rm $(kafka_network) 
	@echo "${RED}All clean!${RESET}"

run:
	@echo "${LIGHTPURPLE}Running mailer app...${RESET}"
	@go run $(main_path)

send:
	@echo "${YELLOW}Sending a mail...${RESET}"
	sh ./scripts/send-test.sh
	@echo "${YELLOW}Mail sent!${RESET}"

kafka_sh:
	@printf "CREATE TOPIC: \n ${GREEN}sh /opt/bitnami/kafka/bin/kafka-topics.sh  --create --topic ${kafka_topic}  --partitions 3  --replication-factor 1 --bootstrap-server ${kafka_host_port}${RESET}\n\n" \

	@printf "LIST TOPIC: \n ${YELLOW}sh /opt/bitnami/kafka/bin/kafka-topics.sh  --list --bootstrap-server ${kafka_host_port}${RESET}\n\n" \

	@printf "PRODUCE TOPIC: \n ${PURPLE}sh /opt/bitnami/kafka/bin/kafka-console-producer.sh  --broker-list ${kafka_host_port} --topic ${kafka_topic}${RESET}\n\n" \

	@printf "SAMPLE MSG PAYLOAD: \n ${RED}{\"from\": \"jed@mail.com\", \"fromName\":\"jed\",\"to\":\"molly@molly.com\",\"subject\":\"a simple hello\", \"messageBody\":{\"errorMessage\":\"hi!\",\"url\":\"www.com\"}} ${RESET}\n\n" \

	@printf "${LIGHTPURPLE}Running producer mode...${RESET}\n"
	docker exec -it kafka bash
	unset JXM_PORT