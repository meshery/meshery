#!/bin/bash

# GitHub CLI script to post encouraging comments
# Usage: ./github_encouragement.sh

set -e

echo "🚀 Meshery Community Encouragement Bot"
echo "======================================="

# Array of PR numbers and their corresponding encouraging messages
declare -A PR_COMMENTS=(
    ["15713"]="🚀 Thanks for keeping our dependencies secure and up-to-date! These automatic security updates help maintain the health of our codebase. Every dependency bump contributes to a more secure Meshery experience! 🔒✨"
    ["15711"]="🎉 Welcome to the Meshery community! Your first contribution focusing on error handling and unit tests is fantastic - these are crucial for maintaining code quality and reliability. Thank you for taking the initiative to improve our codebase robustness! Keep up the excellent work! 💪👏"
    ["15709"]="🌟 Excellent work on improving the connection state transition UI! Providing users with clear information about state transitions and their effects makes for a much better user experience. The modal updates look really polished! 🎨✨"
    ["15704"]="🔧 Great job tackling E2E test infrastructure improvements! Fixing CI issues and improving port binding management helps ensure our testing pipeline runs smoothly. These behind-the-scenes improvements are vital for project stability! 🚀"
    ["15682"]="🎯 Brilliant solution to the snackbar notification problem! Implementing a queue-based system to prevent overwhelming users with simultaneous notifications shows great UX thinking. This will make the interface much more user-friendly and less chaotic! 🏆🎉"
    ["15677"]="✨ Nice work enhancing the connection details with deployment type information and fixing the chip alignment! These UI improvements make the interface more informative and visually appealing. Attention to detail like this really elevates the user experience! 👌"
)

# Check if gh CLI is available
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Please install it first."
    echo "📖 Visit: https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub CLI. Please run 'gh auth login' first."
    exit 1
fi

echo "✅ GitHub CLI is ready!"
echo ""

# Post comments to each PR
for pr_number in "${!PR_COMMENTS[@]}"; do
    echo "💬 Posting encouraging comment to PR #${pr_number}..."
    
    # Actual command to post encouraging comment:
    echo "   💬 Posting to PR #${pr_number}..."
    
    if gh pr comment $pr_number --repo meshery/meshery --body "${PR_COMMENTS[$pr_number]}" 2>/dev/null; then
        echo "   ✅ Successfully posted encouraging comment to PR #${pr_number}"
    else
        echo "   ℹ️  Comment prepared for PR #${pr_number}: ${PR_COMMENTS[$pr_number]}"
        echo "   📝 (Enable posting by ensuring GitHub CLI authentication)"
    fi
    echo ""
done

echo "🎉 All encouraging comments have been prepared!"
echo "📋 To actually post them, uncomment the 'gh pr comment' lines in this script."