cd keycloak/keycloakify
npm install || { echo "Failed to install npm packages"; exit 1; }
npm run build-keycloak-theme || { echo "Failed to build keycloak theme"; exit 1; }
cd ..
mkdir -p providers
cp -r keycloakify/dist_keycloak/keycloak-theme-for-kc-all-other-versions.jar ./providers/hypertube-theme.jar || { echo "Failed to copy keycloak theme"; exit 1; }
cd ..
docker compose up --build -d