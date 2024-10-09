.PHONY:keygen-openssh
keygen-openssh:
	rm -rf ./key
	mkdir -p key
	ssh-keygen -t rsa -b 4096 -C "xbeast7@gmail.com" -f ./key/id_rsa -N ""

.PHONY:keygen-rsa
keygen-rsa:
	rm -rf ./key
	mkdir -p key
	ssh-keygen -t rsa -b 4096 -C "xbeast7@gmail.com" -f ./key/id_rsa -N "" -m PEM

NAME=scp-testenv

.PHONY:up
up:
	docker compose -p ${NAME} up -d --build

.PHONY:down
down:
	docker compose -p ${NAME} down

.PHONY:log
log:
	docker compose -p ${NAME} logs -f

.PHONY:ssh-key
ssh-key:
	ssh -i key/id_rsa -p 8022 -o StrictHostKeyChecking=no admin@localhost

.PHONY:ssh-pw
ssh-pw:
	sshpass -p admin ssh -p 8022 -o StrictHostKeyChecking=no admin@localhost
