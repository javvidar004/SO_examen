FROM mysql:8.0

COPY ./mysqlsampledatabase.sql /docker-entrypoint-initdb.d/

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_USER=api
ENV MYSQL_PASSWORD=api



