FROM postgres:16

WORKDIR /usr/local/bin

COPY scripts/db_healthcheck.sh .
RUN ["chmod", "+x", "db_healthcheck.sh"]
