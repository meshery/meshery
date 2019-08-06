jekyll=bundle exec jekyll

site:
	bundle exec just-the-docs rake search:init
	$(jekyll) serve --drafts --livereload

build-search:
	bundle exec just-the-docs rake search:init

build: build-search
	$(jekyll) build --drafts --livereload

docker:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.5 bash -c "bundle install; bundle exec just-the-docs rake search:init; jekyll serve --drafts --livereload"
