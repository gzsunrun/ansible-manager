FROM 192.168.1.100:5000/library/alpine:3.8-ansible

RUN apk add -U -q --no-progress ansible

COPY ./bin/ansible-manager /usr/local/bin/

ENV PATH=$PATH:/usr/local/bin

WORKDIR /usr/local/bin

CMD ["./ansible-manager"]
