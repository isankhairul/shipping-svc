#! /bin/sh

swagger generate spec -o swagger.yaml --scan-models
docker build -f docker/Dockerfile.local -t aprilsea/as:shipping_svc .
docker push aprilsea/as:shipping_svc