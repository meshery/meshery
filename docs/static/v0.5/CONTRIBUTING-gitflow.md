<p align="center">
<img src="https://img.shields.io/badge/open%20source-%E2%9D%A4%EF%B8%8F-violetgreen">
<img src="https://img.shields.io/badge/PRs-welcome-green">
</p>

## Setting Project on Local System :-

<ul>
<li><h3> Fork the Repository of Project</h3></li><br>

<p><img src="/v0.5/.github/assets/images/fork.png" align="center"></p>

( You will see this on Top Right of Github Repository !)<br><br>
<li><h3>Clone your fork to your local machine</h3></li><br>
<img width="50%" align="center"  src="/v0.5/.github/assets/images/clone.png"><br>
( Click on the Green Code button and Copy the link `https://github.com/........` )<br><br>
<li><h3>Open Git bash where you want to clone the project ( Avoid On Desktop )</h3></li>

<li><h3>Run Command</h3></li><br>

`git clone <insert-link>`
<br>
```(NOTE: In Place of insert-link paste the link you copied)```
<br>
<h3>Project Cloned in System🌟</h3>
<br><br>

<li><h3>Add 'upstream' repo to list of remotes</h3></li><br>
Keeping Your Fork Updated
In order to get the latest updates from the development trunk do a one-time setup to establish the main GitHub repo as a remote by entering:<br><br>

`git remote add upstream https://github.com/meshery/meshery.git`
<br><br>
("meshery" is used as the example repo. Be sure to reference the _actual_ repo you're contributing to e.g. "meshery-linkerd").
<br><br>

<li><h3>Verify the new remote named 'upstream'</h3></li><br>

`git remote -v`
<br><br>
<li><h3>Fetch from upstream remote</h3></li><br>
You'll need to fetch the upstream repo's branches and newest commits to bring them into your repository whenever you wish to update your fork with the latest upstream changes:<br><br>

`git fetch upstream`
<br><br>
<li><h3>Checkout your master branch and merge upstream</h3></li>
<br>
Now, checkout your master branch and merge it with the master branch of the upstream repo:<br><br>

`git checkout master`<br>
`git merge upstream/master`
<br><br>
If the local master branch has no unique commits, git will simply execute a fast-forward. However, if you've been making modifications to master (which, in the vast majority of circumstances, you shouldn't be - see the next section), you may run into issues. Always keep in mind the changes made upstream when doing so.
<br>
Your local master branch is now up to date with everything that has been changed upstream.<br><br>
<li><h3>Create a Branch (to avoid conflicts)</h3></li>
<br>
It's essential to create a new branch whenever you start working on a new feature or bugfix. Not only is it a standard git workflow, but it also organises and separates your modifications from the main branch, allowing you to simply submit and manage several pull requests for each task you finish.
<br>
Follow the steps below to establish a new branch and begin working on it.
<br><br>
<li><h3>Check out the master branch; from which your new branch will be derived.</h3></li><br>

`git checkout master`
<br><br>
<li><h3>Create a new branch</h3></li> (Give your branch a simple, informative name.)
<br>
For continuous integration changes use<br><br>

`ci/your_username/issue#`  <br>OR<br>
 `feature/your_username/name_of_feature`
<br><br>
For bugs use

`bug/your_username/issue#` 
 <br>OR<br>
 `bug/your_username/name_of_bug`
<br>

`git branch feature/jdoe/567`
<br><br><li><h3>Switch to your new branch</h3></li><br>

`git checkout feature/jdoe/567`<br>
```(NOTE: Use the name of the branch you created instead of 'feature/jdoe/567'.)```
<br><br>
Now you may start hacking and make any changes you desire.🚀
<br><br>

<li><h3>Stage the Changes</h3></li><br>

`git add [files-changed]`
<br>
```(NOTE: This will stage all the changes you have made.)```
<br><br>

<li><h3>Commit Changes</h3></li><br>

`git commit -m "MESSAGE"`
<br>
```(NOTE: Instead of 'MESSAGE,' include a commit message so the maintainer can see what you've done.Also make sure to get the DCO signed.)```
<br><br>

<li><h3>Creating Pull Request on Github</h3></li>
<br>
Before submitting your pull request, you should clean up your branch and make it as easy as possible for the maintainer of the original repository to test, accept, and integrate your work.
<br><br>
If any commits to the upstream master branch have been made during the period you've been working on your changes, you'll need to rebase your development branch so that merging it will be a simple fast-forward with no conflict resolution work.
<br><br>

<li><h3>Fetch upstream master and merge with your repo's master branch</h3></li><br>

`git fetch upstream`
<br>
`git checkout master`
<br>
`git merge upstream/master`
<br><br>

<li><h3>If there were any new commits, rebase your development branch</h3></li><br>

`git checkout feature/jdoe/567`<br>
`git rebase master`
<br><br>
Now, it may be desirable to squash some of your smaller commits down into a small number of larger more cohesive commits. You can do this with an interactive rebase:
<br><br>
<li><h3>Rebase all commits on your development branch</li></h3>
<br>

`git checkout`<br>
`git rebase -i master`
<br>

This will open up a text editor where you can specify which commits to squash.
<br><br>
<li><h3>References</h3></li>
<a href="https://git-scm.com/docs">Git Reference Docs</a>
<br>
<a href="https://git-scm.com/docs/git-rebase#_interactive_mode">git-rebase / Interactive Mode</a>
<br><br>
<li><h3>Submit the Changes</h3></li>
Go to the page for your fork on GitHub, select your development branch, then click the pull request button once you've committed and submitted all of your changes. Simply upload the changes to GitHub if you need to make any changes to your pull request.

Your pull request will track and update changes in your development branch automatically.🌸
