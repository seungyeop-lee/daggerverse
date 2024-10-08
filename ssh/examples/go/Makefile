.PHONY:use-password
use-password:
	dagger call -v ssh-use-password-with-option --sshd=tcp://localhost:8022 --password=admin stdout

.PHONY:use-key
use-key:
	dagger call -v ssh-use-key-with-option --sshd=tcp://localhost:8022 --key=file:../helper/key/id_rsa stdout
