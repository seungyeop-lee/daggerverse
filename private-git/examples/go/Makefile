.PHONY:clone-http
clone-http:
	dagger call -v private-git-clone-by-http \
	--git=tcp://localhost:3000 \
	export --path ${PWD}/../helper/repo/http/test

.PHONY:clone-ssh
clone-ssh:
	dagger call -v private-git-clone-by-ssh \
	--git=tcp://localhost:22 \
	--key=file:../helper/key/id_rsa \
	export --path ${PWD}/../helper/repo/ssh/test

.PHONY:push
push:
	dagger call -v private-git-push \
	--git=tcp://localhost:3000 \
	--repo=${PWD}/../helper/repo/http/test \
	stdout
