# git-dc-up
Git pull and execute docker compose up on a remote server when a push happens on a linked repo

### Usage

Start the webserver with following configuration:

```
docker run -p 5000:5000 \
-e WEBHOOK_SECRET=randomsecret \
-e DC_EXTRA_FLAGS="--build" \
-v /root/repo:/repo \
shahidh/git-dc-up
```

### Credits

Web server and build scripts from https://github.com/enricofoltran/simple-go-server/
