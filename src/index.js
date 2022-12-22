const Twitter = require('twitter');
const core = require('@actions/core');
const github = require('@actions/github');

const client = new Twitter({
  consumer_key: core.getInput('consumer-key'),
  consumer_secret: core.getInput('consumer-secret'),
  access_token_key: core.getInput('access-token'),
  access_token_secret: core.getInput('access-token-secret'),
});

const repoName = github.context.payload.repository.full_name;
const { message } = github.context.payload.commits[0];

async function run() {
  try {
    client.post(
      'statuses/update',
      {
        status: `A new PR was merged into ${repoName}: ${message}`,
      },
      (error, tweet, response) => {
        if (error) {
          console.error(error);
        } else {
          console.log(tweet);
        }
      },
    );
  } catch (error) {
    core.setFailed(error.message);
  }
}

run();
