.PHONY:keygen-openssh
keygen-openssh:
	rm -rf ./key
	mkdir -p key
	ssh-keygen -t rsa -b 4096 -C "super@example.com" -f ./key/id_rsa -N ""

.PHONY:keygen-rsa
keygen-rsa:
	rm -rf ./key
	mkdir -p key
	ssh-keygen -t rsa -b 4096 -C "super@example.com" -f ./key/id_rsa -N "" -m PEM

NAME=private-git-testenv

.PHONY:up
up:
	docker compose -p ${NAME} up -d --build

.PHONY:down
down:
	docker compose -p ${NAME} down

.PHONY:clean
clean:
	docker compose -p ${NAME} down -v

.PHONY:log
log:
	docker compose -p ${NAME} logs -f
