FROM ubuntu:18.04
MAINTAINER Vlad
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y gnupg
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

# Клонируем проект
USER root
RUN git clone https://github.com/Vlad104/TP_DB_RK2.git
WORKDIR TP_DB_RK2

# Устанавливаем PostgreSQL
RUN apt-get -y update
RUN apt-get -y install apt-transport-https git wget
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main' >> /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
ENV PGVERSION 11
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib

# Подключаемся к PostgreSQL и создаем БД
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql docker -a -f  database/sql/init.sql &&\
    /etc/init.d/postgresql stop

USER root
# Настраиваем сеть и БД
# COPY database/pg_hba.conf /etc/postgresql/$PGVERSION/main/pg_hba.conf
# COPY database/postgresql.conf /etc/postgresql/$PGVERSION/main/postgresql.conf
RUN echo "local all all md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf

RUN echo "\nlisten_addresses='*'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nport = 5432\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_connections = 100\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nunix_socket_directories = '/var/run/postgresql'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nshared_buffers = '256 MB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nhuge_pages = off\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nwork_mem = '32 MB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmaintenance_work_mem = '256 MB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\neffective_cache_size = '1 GB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nwal_level = minimal\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_wal_senders = 0\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nsynchronous_commit = off\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nfsync = off\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nfull_page_writes = off\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\ncheckpoint_timeout  = '15 min'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\ncheckpoint_completion_target = 0.9'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_wal_size = '1024 MB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmin_wal_size = '512 MB'\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nwal_compression = on\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nwal_buffers = -1\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_worker_processes = 8\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_parallel_maintenance_workers = 4\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_parallel_workers_per_gather = 4\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nparallel_leader_participation = on\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nmax_parallel_workers = 8\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nenable_partitionwise_join = on\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nenable_partitionwise_aggregate = on\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\njit = on\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "\nlogging_collector = off\n" >> /etc/postgresql/$PGVERSION/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432

# Устанавливаем Golang 
ENV GOVERSION 1.11.1
USER root
RUN wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg
ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN apt-get -y install gcc musl-dev && GO11MODULE=on
ENV GOBIN $GOPATH/bin
RUN go get
RUN go build .
EXPOSE 5000
# RUN echo "./config/postgresql.conf" >> /etc/postgresql/$PGVERSION/main/postgresql.conf

# Запускаем PostgreSQL и api сервер
CMD service postgresql start && go run main.go