current_dir = $(shell pwd)
main_path = ./src/main.go

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
	@echo "${RED}All clean!${RESET}"

run:
	@echo "${LIGHTPURPLE}Running mailer app...${RESET}"
	@go run $(main_path)

send:
	@echo "${YELLOW}Sending a mail...${RESET}"
	sh ./scripts/send-test.sh
	@echo "${YELLOW}Mail sent!${RESET}"