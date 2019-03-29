#!/bin/sh
case $1 in
    start)
        if [ ! -f ./docker-compose.yaml ]; then
            echo "Getting the docker-compose.yml file"
            wget https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml
        fi
        if [ ! -f ./docker-compose.yaml ]; then
            echo "Unable to download the docker-compose.yaml file. Please try again later."
        else
            echo "Starting Meshery now. . ."
            docker-compose up -d && \
            echo "Please point your browser to http://localhost:9081 to access the Meshery UI"
        fi
    ;;

    stop)
        echo "Stopping Meshery. . ."
        docker-compose stop
        (docker-compose rm -f) || true
        (docker volume prune -f) || true
        echo "Meshery is now stopped"
    ;;
    logs)
        docker-compose logs -f
    ;;
    *)
        echo "Please enter 'start' or 'stop' to take an appropriate action"
    ;;
esac