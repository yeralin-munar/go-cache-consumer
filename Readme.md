# Intro
Codding interview challenge that reproduce simple service which get data from URL and cache it

# To start cache service
Run
> go mod vendor


> make build


> docker-compose up -d --remove-orphans

# To make request to cache service
> go run consumer/main.go

# You can see log of cache service. Run it before you start consumer
> docker-compose logs -f --tail=10 cache
