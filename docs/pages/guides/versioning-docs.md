---
layout: default
title: Versioning Documentation
description: This guide is to help add support for versioning to docs.
permalink: guides/versioning-docs
type: Guides
---

This guide helps to add support for versioning in docs for new major releases or for older releases. 

## For new major release

The structure which the docs follow right now is, The main `docs` folder has the most recent version of documentation, while there are sub-folders for previous versions, v0.x (x being the last major release).
On release of a new major version, the static html files for the most recent version is generated and is renamed as the release version (v0.x). 

#### Steps:
After cloning the Meshery repository
1. `cd docs` > `bundle install` > `make site`
1. On executing `make site` a `_site` folder is created which has static html files. 
1. The `_site` folder is renamed to `v0.x`. 
1. This `v0.x` folder is now the latest version of docs. 

##### __In the `v0.x` folder__
1. Search and replace all the instances where there is a direct path is defined to include the version name in the path, i.e, all paths to intra-page links and images should start with `/v0.x/`.
- Look for `href="/` and replace with `href="/0.x/`
- Look for `src="/`and replace with `src="/0.x/` <br/><br/>
<a href="{{ site.baseurl }}/assets/img/versioning-guide/search-and-replace.png">
  <img src="{{ site.baseurl }}/assets/img/versioning-guide/search-and-replace.png" />
</a>


## For old release

For older releases we have to travel back in time. Using the `Tags` in github we go to a previous release, `v0.X.x`, the `.x` here should be the latest version of the archived docs. 

#### Steps: 
1. Copy the commit ID for that release. <br/><br/>
<a href="{{ site.baseurl }}/assets/img/versioning-guide/commit-ID.png">
  <img src="{{ site.baseurl }}/assets/img/versioning-guide/commit-ID.png" />
</a>

1. `git checkout <commit ID>` > `cd docs` > `bundle install` > `make site`
1.  On executing `make site` a `_site` folder is created which has static html files. 
1.  The `_site` folder is renamed to `v0.X` and is copied into the `docs` folder of the present version. 

The above [steps]({{site.baseurl}}/guides/versioning-docs#in-the-v0x-folder) for replacing all the instances of direct path are to be followed. 


