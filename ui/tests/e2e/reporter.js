const github = require('@actions/github');
const core = require('@actions/core');

class MyReporter {
  async onBegin(config, suite) {
    if (!process.env.CI) {
      console.log(`Starting the run with ${suite.allTests().length} tests`);
      return;
    }

    async function run() {
      try {
        const github_token = core.getInput('GITHUB_TOKEN');

        const context = github.context;

        if (context.payload.pull_request) {
          const pr_number = context.payload.pull_request.number;

          const octokit = github.getOctokit(github_token);

          await octokit.rest.pulls.createReviewComment({
            owner: 'Meshery[bot]',
            repo: context.repo,
            pull_number: pr_number,
            body: `Starting the run with ${suite.allTests().length} tests`,
          });
        }
      } catch (error) {
        core.setFailed(error.message);
      }
    }

    run();
  }

  onTestBegin(test) {
    console.log(`Starting test ${test.title}`);
  }

  onTestEnd(test, result) {
    console.log(`Finished test ${test.title}: ${result.status}`);
  }

  onEnd(result) {
    if (!process.env.CI) {
      console.log(`Finished the run: ${result}`);
      return;
    }

    async function run() {
      try {
        const github_token = core.getInput('GITHUB_TOKEN');

        const context = github.context;

        if (context.payload.pull_request) {
          const pr_number = context.payload.pull_request.number;

          const octokit = github.getOctokit(github_token);

          await octokit.rest.pulls.createReviewComment({
            owner: 'Meshery[bot]',
            repo: context.repo,
            pull_number: pr_number,
            body: `Finished the run: ${result}`,
          });
        }
      } catch (error) {
        core.setFailed(error.message);
      }
    }

    run();
  }
}

module.exports = MyReporter;
