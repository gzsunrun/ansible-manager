# Ansible-Manager
#
# VERSION               2.0

FROM centos

RUN yum -y install ansible
RUN yum -y install git
COPY ansible-manager /usr/local/bin/ansible-manager
EXPOSE 8090
CMD /usr/local/bin/ansible-manager