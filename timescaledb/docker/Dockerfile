FROM timescale/timescaledb:2.15.1-pg16

# or
# ARG BUILD_FROM
# FROM ${BUILD_FROM}

ENV LANG=C.UTF-8

# Install dependencies (optional)
RUN apt-get update && apt-get install -y \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Copy init scripts
COPY run.sh /run.sh
COPY init.sql /docker-entrypoint-initdb.d/init.sql
RUN chmod a+x /run.sh

# Use the official entrypoint
ENTRYPOINT ["/run.sh"]

# Expose Postgres port
EXPOSE 5432
