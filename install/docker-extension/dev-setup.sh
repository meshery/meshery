
#!/usr/env/bin bash

set -e

EXTENSION_IMAGE="meshery/docker-extension-meshery:edge-latest-dev"
EXTENSION_NAME="meshery/docker-extension-meshery"

echo "🔧 Removing existing Meshery Docker extension if present..."
docker extension rm $EXTENSION_NAME || echo "No extension to remove."

make extension

echo "📦 Installing Meshery Docker extension from $EXTENSION_IMAGE..."
docker extension install $EXTENSION_IMAGE --force


echo "🐞 Enabling Docker Desktop debug mode..."
docker extension dev debug meshery/docker-extension-meshery

echo "✅ Meshery Docker extension installed successfully."
