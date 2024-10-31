#!/bin/bash

echo "Stopping all containers..."
docker-compose down

echo "Removing database volume..."
docker volume rm $(docker volume ls -q | grep db_postgres)

echo "Starting services..."
docker-compose up -d

echo "Waiting for database to be ready..."
sleep 10

echo "Running seeder..."
docker-compose --profile seeder up seeder

echo "Database reinitialized and seeded!"