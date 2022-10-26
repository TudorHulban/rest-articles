#!/bin/sh
psql -h localhost -U user -W -d database_name -f sql/init.sql