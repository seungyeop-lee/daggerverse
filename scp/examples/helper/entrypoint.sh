#!/bin/sh

exec /usr/sbin/sshd -D -e "$@"
