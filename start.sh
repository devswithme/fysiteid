docker compose down -v --rmi all
docker system prune -af --volumes

docker network create nginx

docker compose down -v
docker compose up -d

docker ps