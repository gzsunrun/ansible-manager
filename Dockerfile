# Ansible-Manager
#
# VERSION               2.0

FROM centos

RUN yum -y install ansible
RUN yum -y install git
COPY ansible-manager /usr/local/bin/ansible-manager
COPY public/ /usr/local/html/ansible-manager/public/
EXPOSE 8090
CMD /usr/local/bin/ansible-manager