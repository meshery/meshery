import { readFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import path from 'path';

class MyReporter {
  introMessage = '';
  failsMessage = '';
  passed = 0;
  failed = 0;
  skipped = 0;

  onBegin(config, suite) {
    this.introMessage = `- Test run started at ${moment().format('MMMM Do YYYY, h:mm:ss a')}
- Number tests cases to run: ${suite.allTests().length}`;
  }

  onTestEnd(test, result) {
    switch (result.status) {
      case 'failed':
      case 'timedOut':
        this.addFailMessage(`❌ Test ${test.title} failed\n>${result.error?.message}`);
        this.failed++;
        break;
      case 'skipped':
        this.addFailMessage(`⚠️ Test ${test.title} skipped`);
        this.skipped++;
        break;
      case 'passed':
        this.passed++;
        break;
    }
  }

  async onEnd(result) {
    const message = await this.buildMessage(result);
    console.log(message);
    process.exit(this.failed > 0 ? 130 : 0); // Return non-zero status code if there are failed tests
  }

  addFailMessage(message) {
    this.failsMessage += `\n${message}`;
  }

  async buildMessage(result) {
    const duration = moment.duration(result.duration, 'milliseconds');
    const minutes = Math.floor(duration.asMinutes());
    const seconds = duration.seconds();
    const templateStr = readFileSync(path.join(__dirname, 'reporterSummary.md'), 'utf8');
    return template(templateStr)({
      introMessage: this.introMessage,
      minutes,
      seconds,
      passed: this.passed,
      failed: this.failed,
      skipped: this.skipped,
      failsMessage: this.failsMessage,
    });
    // const details = `
    //     <details>
    //             <summary>Click Here for more details</details>
    //             ${this.failsMessage}
    //     </details>`;

    // const resultMarkdownMessage = `
    // Test run results
    //   ---
    //   ${this.introMessage}
    //   ---
    //   Summary:
    //     - ⌛ Duration  of test run: ${minutes} minutes and ${seconds} seconds
    //     - 📦 Tests results:
    //       - ✅ ${this.passed}
    //       - ❌ ${this.failed}
    //       - ⏩ ${this.skipped}

    //   ${this.failsMessage ? "👍 All tests passed successfully!" : "👎 Some tests failed!"}

    //   ${this.failsMessage && details}

    //   To see the full report, please visit our CI/CD pipeline with reporter.`;
    // return resultMarkdownMessage;
  }
}

export default MyReporter;
