#!/bin/bash

chmod 777 /pg_logs

chown postgres /pg_setup/server.key
chown postgres /pg_setup/server.crt
su postgres -c "chmod 600 /pg_setup/server.key"
su postgres -c "chmod 600 /pg_setup/server.crt"

su postgres -c "postgres -c config_file=/pg_setup/postgresql.conf -c hba_file=/pg_setup/pg_hba.conf"

echo "======> POSTGRES EXIT CODE $?"
