import github from '@actions/github';
import core from '@actions/core';

async function runGithub() {
  if (!process.env.CI) {
    return;
  }
  try {
    core.debug('Initialize GH action command');

    const gh_token = core.getInput('GITHUB_TOKEN');
    const octokit = github.getOctokit(gh_token);

    const context = github.context;

    const pr_number = context.payload.pull_request.number;

    const sendComment = (message) => {
      core.debug('Send comment to an issue');
      octokit.rest.issues.createComment({
        ...context.repo,
        issue_number: pr_number,
        body: message,
      });
    };

    return { sendComment };
  } catch (err) {
    core.setFailed(`GH Action api failed with error, ${err}`);
  }
}

export default runGithub;
