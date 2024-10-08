FROM alpine:3.20.3

# 필요한 패키지 설치
RUN apk update && \
    apk add --no-cache openssh-server sudo

# sshd 설정
RUN mkdir /var/run/sshd \
    && echo 'PubkeyAuthentication yes' >> /etc/ssh/sshd_config \
    && echo 'AuthorizedKeysFile /etc/ssh/authorized_keys' >> /etc/ssh/sshd_config \
    && echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config \
    && ssh-keygen -A # ssh 접속용 호스트키 생성, /etc/ssh/에 저장 됨.

# ssh 접속 유저 추가
RUN adduser -h /home/admin -s /bin/sh -D admin \
    && echo "admin ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/admin \
    && echo -n 'admin:admin' | chpasswd

# 공개키 복사
COPY key/id_rsa.pub /home/admin/.ssh/authorized_keys
RUN chown -R admin:admin /home/admin/.ssh \
    && chmod 700 /home/admin/.ssh \
    && chmod 600 /home/admin/.ssh/authorized_keys

EXPOSE 22

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
