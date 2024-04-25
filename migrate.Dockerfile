FROM migrate/migrate:v4.17.1

WORKDIR /usr/local/bin

RUN ["apk", "add", "--no-cache", "bash", "jq"]

COPY scripts/db_secret.sh scripts/db_migrate.sh ./
RUN ["chmod", "+x", "db_migrate.sh"]

ENTRYPOINT ["db_migrate.sh"]
