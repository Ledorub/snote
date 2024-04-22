FROM postgres:16

WORKDIR /usr/local/bin

COPY scripts/db_*.sh .
RUN ["chmod", "+x", "db_startup.sh", "db_healthcheck.sh"]

ENTRYPOINT ["db_startup.sh", "docker-entrypoint.sh"]
CMD ["postgres"]
