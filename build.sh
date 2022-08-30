#! /bin/sh
#comment this statement for swagger generate on build due to problem on building. Now Swagger.yaml is manually generated and push to branch develop
#swagger generate spec -o swagger.yaml --scan-models
docker build -f docker/Dockerfile.local -t aprilsea/as:shipping_svc .
docker push aprilsea/as:shipping_svc