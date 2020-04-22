# Meshery Docs

![assets/img/docs-screenshot.png](assets/img/docs-screenshot.png)
*Check out the [Meshery Docs](https://meshery.layer5.io/docs/)!*

## Contributing to the Meshery Documentation

Before contributing, please review the [Documentation Contribution Flow](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). In the following steps you will set up your development environment, fork and clone the repository, run the site locally, and finally commit, sign-off, and push any changes made for review. 

### 1. Set up your development environment

* *The Meshery Docs site is built using Jekyll - a simple static site generator! You can learn more about Jekyll and setting up your development environment in the [Jekyll Docs](https://jekyllrb.com/docs/).*

* First [install Ruby](https://jekyllrb.com/docs/installation/), then install Jekyll and Bundler.

### 2. Get the code

* Fork and then clone the [Meshery repository](https://github.com/layer5io/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```
* Change to the docs directory
  ```bash
  $ cd docs
  ```
* Install any Ruby dependencies
  ```bash
  $ bundle install
  ```

### 3. Serve the site

* Serve the code locally
  ```bash
  $ make site
  ```
  *Note: From the Makefile, this command is actually running `$ bundle exec jekyll serve --drafts --livereload`*

### 4. Create a Pull Request

* After making changes, don't forget to commit with the sign-off flag (-s)!
  ```bash
  $ commit -s -m “my commit message w/signoff”
  ```
* Once all changes have been committed, push the changes.
  ```bash
  $ git push origin <branch-name>
  ```
* Then on Github, navigate to the [Meshery repository](https://github.com/layer5io/meshery) and create a pull request from your recently pushed changes!

---
*See the [Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) for additional reference.*