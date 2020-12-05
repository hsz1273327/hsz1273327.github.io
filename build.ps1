docker buildx build --load --platform=linux/arm64 -t hsz1273327/gh-oauth-server:arm64-0.0.2 -t hsz1273327/gh-oauth-server:arm64-latest .
docker buildx build --load --platform=linux/amd64 -t hsz1273327/gh-oauth-server:amd64-0.0.2 -t hsz1273327/gh-oauth-server:amd64-latest .
docker buildx build --load --platform=linux/arm/v7 -t hsz1273327/gh-oauth-server:armv7-0.0.2 -t hsz1273327/gh-oauth-server:armv7-latest .
docker push hsz1273327/gh-oauth-server

docker manifest create --amend hsz1273327/gh-oauth-server:0.0.2 hsz1273327/gh-oauth-server:arm64-0.0.2 hsz1273327/gh-oauth-server:amd64-0.0.2 hsz1273327/gh-oauth-server:armv7-0.0.2
docker manifest push --purge hsz1273327/gh-oauth-server:0.0.2

docker manifest create --amend hsz1273327/gh-oauth-server:latest hsz1273327/gh-oauth-server:arm64-latest hsz1273327/gh-oauth-server:amd64-latest hsz1273327/gh-oauth-server:armv7-latest
docker manifest push --purge hsz1273327/gh-oauth-server:latest