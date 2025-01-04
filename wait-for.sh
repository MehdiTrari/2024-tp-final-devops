#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

until nc -z "$host" 5432; do
    echo "Waiting for PostgreSQL..."
    sleep 1
done

echo "PostgreSQL is up - executing command"
exec $cmd
