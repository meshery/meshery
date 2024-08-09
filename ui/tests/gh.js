const github = require('@actions/github');
const core = require('@actions/core');

async function runGithub() {
  try {
    const gh_token = core.getInput('GITHUB_TOKEN');
    const octokit = github.getOctokit(gh_token);

    const context = github.context;

    const pr_number = context.payload.pull_request.number;

    const sendComment = (message) => {
      octokit.issues.createComment({
        ...context.repo,
        issue_number: pr_number,
        body: message,
      });
    };

    return { sendComment };
  } catch (error) {
    core.setFailed(error.message);
  }

  return;
}

export default runGithub;
