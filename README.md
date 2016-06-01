
## What is this?

This is a [logspout](https://github.com/gliderlabs/logspout) adapter that will parse logs for metrics and 
forward them to statsd. 

## Trying it out

You can test this by running through the following walkthrough. 

NOTE: If you're on a mac, replace localhost below with your `docker-machine ip`.

First, start up stats/graphite with this prebuilt image: 

```sh
docker run -d\
 --name graphite\
 --restart=always\
 -p 80:80\
 -p 2003-2004:2003-2004\
 -p 2023-2024:2023-2024\
 -p 8125:8125/udp\
 -p 8126:8126\
 hopsoft/graphite-statsd
 ```

Test if it's working:

```sh
echo "foo:1|c" | nc -u -w0 localhost 8125
```

Go to `http://localhost/dashboard`, you should see `stats.foo`.

Copy the files in the `custom-build` in this repo, then run:

```sh
docker build -t mylogspouter .
```

Now start up logspout with new module, replace the papertrail URL to your own account. 

```sh
docker run --rm --name=logspout \
     -v=/var/run/docker.sock:/var/run/docker.sock \
     -p 8000:80 \
     mylogspouter \
     syslog://logs.papertrailapp.com:55555
```

Now we have to add the statsd route;

```sh
curl http://localhost:8000/routes -d '{
  "adapter": "statsd",
  "filter_sources": ["stdout" ,"stderr"],
  "address": "localhost:8125"
}'
```

Now run our [emitter image](https://github.com/iron-io/logspout-statsd/tree/master/test/stats-emitter) to test:

```sh
docker run --rm iron/emitter 
```

You should see some logging messages in papertrail and some metrics in graphite!

Any log line from any container that contains a `metric` key will be forwarded to statsd, so your applications should output log lines like this:

```
metric=someevent value=1 type=count
metric=somegauge value=50 type=gauge
```

## Development of this module

See iron-io's fork of logspout and read MODULES.md.

Copy run-logspout.sh to the logspout dir and run it with:

```sh
SYSLOG=syslog://logs.papertrailapp.com:55555 ./run-logspout.sh`
```
