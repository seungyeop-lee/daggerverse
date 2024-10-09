.PHONY:copy-to-remote
copy-to-remote:
	dagger call -v scp-copy-to-remote-with-option \
	--sshd=tcp://localhost:8022 \
	--key=file:../helper/key/id_rsa \
	--file=../helper/testdata/hello.txt \
	--target=/home/admin \
	stdout

.PHONY:copy-from-remote
copy-from-remote:
	dagger call -v scp-copy-from-remote-with-option \
	--sshd=tcp://localhost:8022 \
	--key=file:../helper/key/id_rsa \
	--path=/home/admin/hello.txt \
	contents
