FROM postgres:15

WORKDIR /migrations

COPY migrations/ /migrations/
COPY run-migrations.sh /run-migrations.sh

RUN chmod +x /run-migrations.sh

ENTRYPOINT ["/run-migrations.sh"]
