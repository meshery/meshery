# Docsy Jekyll Theme

[![CircleCI](https://circleci.com/gh/vsoch/docsy-jekyll/tree/master.svg?style=svg)](https://circleci.com/gh/vsoch/docsy-jekyll/tree/master)
<a href="https://jekyll-themes.com/docsy-jekyll/">
    <img src="https://img.shields.io/badge/featured%20on-JT-red.svg" height="20" alt="Jekyll Themes Shield" >
</a>

![https://raw.githubusercontent.com/vsoch/docsy-jekyll/master/assets/img/docsy-jekyll.png](https://raw.githubusercontent.com/vsoch/docsy-jekyll/master/assets/img/docsy-jekyll.png)

This is a [starter template](https://vsoch.github.com/docsy-jekyll/) for a Docsy jekyll theme, based
on the Beautiful [Docsy](https://github.com/google/docsy) that renders with Hugo. This version is intended for
native deployment on GitHub pages. The original [Apache License](https://github.com/vsoch/docsy-jekyll/blob/master/LICENSE) is included.

## Changes

The site is intended for purely documentation, so while the front page banner
is useful for business or similar, this author (@vsoch) preferred to have
the main site page go directly to the Documentation view. Posts
are still provided via a feed.

## Usage

### 1. Get the code

You can clone the repository right to where you want to host the docs:

```bash
git clone https://github.com/vsoch/docsy-jekyll.git docs
cd docs
```

### 2. Customize

To edit configuration values, customize the [_config.yml](https://github.com/vsoch/docsy-jekyll/blob/master/_config.yml).
To add pages, write them into the [pages](https://github.com/vsoch/docsy-jekyll/blob/master/pages) folder. 
You define urls based on the `permalink` attribute in your pages,
and then add them to the navigation by adding to the content of [_data/toc.myl](https://github.com/vsoch/docsy-jekyll/blob/master/_data/toc.yml).
The top navigation is controlled by [_data/navigation.yml](https://github.com/vsoch/docsy-jekyll/blob/master/_data/navigation.yml)

### 3. Options

Most of the configuration values in the [_config.yml](https://github.com/vsoch/docsy-jekyll/blob/master/_config.yml) are self explanatory,
and for more details, see the [getting started page](https://vsoch.github.io/docsy-jekyll/docs/getting-started)
rendered on the site.

### 4. Serve

Depending on how you installed jekyll:

```bash
jekyll serve
# or
bundle exec jekyll serve
```
