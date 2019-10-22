jekyll=bundle exec jekyll

site:
	$(jekyll) serve

build:
	$(jekyll) build

docker:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.5 bash -c "bundle install; jekyll serve"
