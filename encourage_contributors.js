#!/usr/bin/env node

/**
 * Script to post encouraging comments on the latest open pull requests
 * This promotes a welcoming and supportive community culture
 */

const fs = require('fs');
const path = require('path');

// Define the 6 open PRs and their encouraging messages
const prComments = [
  {
    number: 15713,
    title: "build(deps-dev): Bump sha.js from 2.4.11 to 2.4.12 in /ui",
    author: "dependabot[bot]",
    message: "🚀 Thanks for keeping our dependencies secure and up-to-date! These automatic security updates help maintain the health of our codebase. Every dependency bump contributes to a more secure Meshery experience! 🔒✨"
  },
  {
    number: 15711,
    title: "added error handling , wrote unit test cases",
    author: "ShivanshCharak",
    message: "🎉 Welcome to the Meshery community! Your first contribution focusing on error handling and unit tests is fantastic - these are crucial for maintaining code quality and reliability. Thank you for taking the initiative to improve our codebase robustness! Keep up the excellent work! 💪👏"
  },
  {
    number: 15709,
    title: "Modify connection state transition , with updated ui && description.",
    author: "FaheemOnHub",
    message: "🌟 Excellent work on improving the connection state transition UI! Providing users with clear information about state transitions and their effects makes for a much better user experience. The modal updates look really polished! 🎨✨"
  },
  {
    number: 15704,
    title: "fix(mesheryctl-e2e): update how port binding in e2e is managed",
    author: "lekaf974",
    message: "🔧 Great job tackling E2E test infrastructure improvements! Fixing CI issues and improving port binding management helps ensure our testing pipeline runs smoothly. These behind-the-scenes improvements are vital for project stability! 🚀"
  },
  {
    number: 15682,
    title: "Improve Snackbar: Implement Queue based system for snackbar notifications.",
    author: "FaheemOnHub", 
    message: "🎯 Brilliant solution to the snackbar notification problem! Implementing a queue-based system to prevent overwhelming users with simultaneous notifications shows great UX thinking. This will make the interface much more user-friendly and less chaotic! 🏆🎉"
  },
  {
    number: 15677,
    title: "add: deployment type && fix: option-alignment with chips",
    author: "FaheemOnHub",
    message: "✨ Nice work enhancing the connection details with deployment type information and fixing the chip alignment! These UI improvements make the interface more informative and visually appealing. Attention to detail like this really elevates the user experience! 👌"
  }
];

// Export for potential use by other scripts
module.exports = { prComments };

console.log('📝 Prepared encouraging comments for 6 open pull requests:');
prComments.forEach(pr => {
  console.log(`\n🔸 PR #${pr.number}: ${pr.title}`);
  console.log(`   Author: ${pr.author}`);
  console.log(`   Message: ${pr.message}`);
});

console.log('\n✅ Comments prepared! Use GitHub CLI or API to post these to the respective PRs.');