jekyll=bundle exec jekyll

site: install
	$(jekyll) serve --drafts --livereload --config _config_dev.yml

build: install
	$(jekyll) build --drafts

build-graphql-docs:
	bundle exec rake graphql:compile_docs

install:
	bundle install

docker:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:4.0.0 bash -c "bundle install; jekyll serve --drafts --livereload --config _config.yml,_config_dev.yml"
