set -ex

curl http://192.168.99.100:8000/routes -d '{
  "adapter": "statsd",
  "filter_sources": ["stdout" ,"stderr"],
  "address": "192.168.99.100:8125"
}'

docker build -t iron/emitter .
docker run --rm iron/emitter
