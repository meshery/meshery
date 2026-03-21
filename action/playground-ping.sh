# Pings the Playground active server
# returns true if the playground responds wit 200
# else returns false

PLAYGROUND_PING_URL=https://playground.meshery.io/api/system/version
export V=($(curl --silent  $PLAYGROUND_PING_URL  -w "%{http_code}" | grep 200))

if [ ! -z "$V" ] && [ $V == 200 ]
then
  echo true
else 
  echo false
fi
