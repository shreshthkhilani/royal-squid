# royal-squid

```sh
cp .env.sample .env
docker build -t royal-squid api
docker run -it --rm --name royal-squid -p 8080:8080 --env-file .env royal-squid
```