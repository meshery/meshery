	if [ ! "$(docker ps -q -f name=pg)" ]; then
		if [ "$(docker ps -aq -f status=exited -f name=pg)" ]; then
			# cleanup
			docker rm pg
		fi
    # run your container
    docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=meshery --name pg postgres
	fi