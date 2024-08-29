import { readFileSync, writeFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import path from 'path';

class MyReporter {
  introMessage = '';
  expectedTest = '';
  testTable = ` No | Project | Test Case | Status | Retry |
| :---: | :---: | :--- | :---: | :---: |`;
  passed = 0;
  failed = 0;
  skipped = 0;
  flaky = 0;
  count = 1;

  onBegin(_config, suite) {
    this.introMessage = `- Test run started at ${moment().format('MMMM Do YYYY, h:mm:ss a')}
- Number tests cases to run: ${suite.allTests().length}`;
  }

  // eslint-disable-next-line no-unused-vars
  onStdOut(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stdout.write(text);
  }

  // eslint-disable-next-line no-unused-vars
  onStdErr(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stderr.write(text);
  }

  onTestEnd(test, result) {
    const status = test.outcome();
    const project = test.parent.project().name;

    const message = `| ${this.count} | ${project} | ${test.title} | ${this.getStatusEmoji(
      status,
    )} | ${result.retry} |`;
    const logs = `${this.count}. Project: ${project}, Test: ${
      test.title
    }, Status: ${this.getStatusEmoji(status)}, Retry: ${result.retry}\n`;

    process.stdout.write(logs);
    this.addTestTable(message, status);
    this.countTest(status);
  }

  async onEnd(result) {
    const message = await this.buildMessage(result);

    try {
      writeFileSync('test-report.md', message);
    } catch (e) {
      console.log('Cannot write test reporter ', e);
    }
  }

  addTestTable(message, status) {
    if (status === 'expected') return;
    this.testTable += `\n${message}`;
  }

  countTest(status) {
    if (status === 'expected') {
      this.passed++;
    }
    if (status === 'unexpected') {
      this.failed++;
    }
    if (status === 'flaky') {
      this.flaky++;
    }
    if (status === 'skipped') {
      this.skipped++;
    }
    this.count++;
  }

  getStatusEmoji(status) {
    if (status === 'expected') {
      return '✅';
    }
    if (status === 'unexpected') {
      return '❌';
    }
    if (status === 'flaky') {
      return '⚠️';
    }
    if (status === 'skipped') {
      return '➖';
    }
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
      flaky: this.flaky,
      skipped: this.skipped,
      testTable: this.testTable,
    });
  }
}

export default MyReporter;
