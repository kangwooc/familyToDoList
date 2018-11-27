# Use the name of RDBMSs base image, mysql, in the FROM statement,
# and then COPY schema.sql into /docker-entrypoint-initdb.d/schema.sql.
FROM mysql
COPY schema.sql /docker-entrypoint-initdb.d/schema.sql