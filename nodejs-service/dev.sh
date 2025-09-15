#!/bin/bash

case "$1" in
  up)
    echo "Starting containers..."
    docker compose up -d
    ;;
  down)
    echo "Stopping containers..."
    docker compose down -v
    ;;
  restart)
    echo "Restarting containers..."
    docker compose down
    docker compose up -d
    ;;
  logs)
    docker compose logs -f
    ;;
  clean)
    echo "cleaning..."
    docker compose down -v --rmi all --remove-orphans
    ;;
esac