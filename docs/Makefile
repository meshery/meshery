jekyll=bundle exec jekyll

site: install
	$(jekyll) serve --drafts --livereload --config _config.yml,_config_dev.yml

build: install
	$(jekyll) build --drafts

install:
	bundle install

docker:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.7 bash -c "bundle install; jekyll serve --drafts --livereload --config _config.yml,_config_dev.yml"