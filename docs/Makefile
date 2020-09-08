jekyll=bundle exec jekyll

site:
	$(jekyll) serve --drafts --livereload

build:
	$(jekyll) build --drafts

docker:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.7 bash -c "bundle install; jekyll serve --drafts --livereload"
