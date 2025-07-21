#!/bin/sh
# Ensure required environment variables are set
: "${INTRA_CLIENT_ID:?Environment variable INTRA_CLIENT_ID is required}"
: "${INTRA_CLIENT_SECRET:?Environment variable INTRA_CLIENT_SECRET is required}"
rm -f /run/rsyslogd.pid
rsyslogd
find /opt/keycloak/data/import -type f -name "*.json" -exec sed -i \
  -e "s|INTRA_CLIENT_ID|${INTRA_CLIENT_ID}|g" \
  -e "s|INTRA_CLIENT_SECRET|${INTRA_CLIENT_SECRET}|g" {} +
/opt/keycloak/bin/kc.sh build
exec /opt/keycloak/bin/kc.sh start \
    -Dkeycloak.migration.action=import \
    -Dkeycloak.migration.provider=dir \
    -Dkeycloak.migration.dir=/opt/keycloak/data/import \
    -Dkeycloak.migration.strategy=OVERWRITE_EXISTING \
    -Dkeycloak.migration.replace-placeholders=true \
    --proxy-headers=xforwarded