.PHONY: docs build install docker graphql-docs-build
jekyll=bundle exec jekyll

docs: install
	$(jekyll) serve --drafts --incremental --config _config_dev.yml

build: install
	$(jekyll) build --drafts

graphql-docs-build:
	bundle exec rake graphql:compile_docs

install:
	bundle install

docker:
	docker build -t jekyll-with-ruby3.2.2 .
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll-with-ruby3.2.2 bash -c "bundle install; jekyll serve --drafts --livereload --config _config.yml,_config_dev.yml"
