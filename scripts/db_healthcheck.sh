#!/bin/bash

pg_isready -h 127.0.0.1 -p $PGPORT -U $POSTGRES_USER -d $POSTGRES_DB
