#!/bin/bash

echo "Stopping all containers..."
docker-compose down

echo "Starting services..."
docker-compose up -d

echo "Waiting for database to be ready..."
sleep 10

echo "Running seeder..."
docker-compose --profile seeder up seeder

echo "Database reinitialized and seeded!"