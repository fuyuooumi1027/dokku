# Docker Installation Notes

Pull the dokku/dokku image:

```shell
docker pull dokku/dokku:v0.17.7
```

Next, run the image.

```shell
docker run \
  --env DOKKU_HOSTNAME=dokku.me \
  --name dokku \
  --publish 3022:22 \
  --publish 8080:80 \
  --publish 8443:443 \
  --volume /mount:/mount \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  dokku/dokku:v0.17.7
```

Dokku is run with the following options:

- The global hostname is set to `dokku.me` on boot.
- The container name is dokku.
- Container SSH port 22 is exposed on the host as 3022.
- Container HTTP port 80 is exposed on the host as 8000.
- Container HTTPS port 443 is exposed on the host as 8443.
- Data within the container is stored on the host within the `/mount` directory.
- The docker socket is mounted into container

Application repositories, plugin config, and plugin data is persisted to disk within the specified host directory for `/mount`.

To install custom plugins, create a `plugin-list` file in the host's `/mount` directory. The plugins listed herein will be automatically installed by Dokku on container boot. This file should be the following format:

```yaml
plugin_name: repository_url
```

An example for installing the postgres and redis plugins follows:

```yaml
postgres: https://github.com/dokku/dokku-postgres.git
redis: https://github.com/dokku/dokku-redis.git
```

To initialize ssh-keys within the container, use `docker exec` to enter the container and run the appropriate ssh-keys commands.

```shell
docker exec -it dokku bash
```

Please see the [user management documentation](/docs/deployment/user-management.md) for more information.
