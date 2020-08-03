FROM centos:centos7

RUN echo "root ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Add source repositories
ADD src/epel7.repo /etc/yum.repos.d/epel7.repo
ADD src/datastax.repo /etc/yum.repos.d/datastax.repo

RUN yum clean all

# Install Java, Install packages (sshd + supervisord + monitoring tools + cassandra)
RUN yum install -y wget tar openssh-server openssh-clients supervisor sysstat sudo which openssl hostname
RUN yum install -y java-1.8.0-openjdk-headless
RUN yum install -y dsc22

RUN mkdir -p /var/run/sshd && chmod -rx /var/run/sshd && \
	ssh-keygen -t rsa -N '' -f /etc/ssh/ssh_host_rsa_key && \
	sed -ri 's/#PermitRootLogin yes/PermitRootLogin yes/g' /etc/ssh/sshd_config && \
	sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config && \
	sed -ri 's/#UsePAM no/UsePAM no/g' /etc/ssh/sshd_config && \
	useradd -m -G users,root -p $(openssl passwd -1 "opscenter") opscenter && \
	echo "%root ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

ADD src/cassandra.yaml /etc/cassandra/conf/cassandra.yaml

# Configure supervisord
ADD src/supervisord.conf /etc/supervisord.conf
RUN mkdir -p /var/log/supervisor

# Deploy startup script
ADD src/start.sh /usr/local/bin/start

# Clone the CRUD service from git
RUN yum install -y golang
RUN yum install -y git

RUN git clone https://github.com/rherbut/simple_rest_api.git
RUN chmod 777 /simple_rest_api/*
ENV GOPATH=/simple_rest_api
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gocql/gocql
RUN go get github.com/hashicorp/consul/api

EXPOSE 7199 7000 7001 9160 9042
EXPOSE 22 8012 61621 
EXPOSE 8080

RUN chmod 777 /usr/local/bin/start
CMD start
