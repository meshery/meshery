@echo off
REM Git Line Ending Fix Script for Windows
REM This script configures Git to handle line endings properly on Windows systems

echo 🔧 Fixing Git Line Ending Configuration...

REM Check if we're in a Git repository
git rev-parse --git-dir >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Not in a Git repository. Please run this script from within a Git repository.
    pause
    exit /b 1
)

echo 📁 Current directory: %CD%

REM Configure Git settings for the current repository
echo ⚙️  Configuring local Git settings...
git config core.autocrlf false
git config core.safecrlf warn
git config core.eol lf

echo ✅ Local Git configuration updated:
echo    - core.autocrlf: false
echo    - core.safecrlf: warn
echo    - core.eol: lf

REM Ask user if they want to apply global settings
set /p "global=🌍 Do you want to apply these settings globally for all Git repositories? (y/N): "
if /i "%global%"=="y" (
    echo ⚙️  Configuring global Git settings...
    git config --global core.autocrlf false
    git config --global core.safecrlf warn
    echo ✅ Global Git configuration updated
) else (
    echo ℹ️  Skipped global configuration
)

REM Check if .gitattributes exists
if exist ".gitattributes" (
    echo 📄 .gitattributes file already exists
    set /p "backup=🔄 Do you want to backup and update it? (y/N): "
    if /i "!backup!"=="y" (
        copy .gitattributes .gitattributes.backup >nul
        echo 💾 Backup created: .gitattributes.backup
        echo ⚠️  Please manually merge the new .gitattributes content with your existing file
    )
) else (
    echo 📄 No .gitattributes file found. You should create one with proper line ending rules.
)

REM Show current Git configuration
echo.
echo 📋 Current Git configuration:
echo    Local settings:
git config --local --get-regexp "core\.(autocrlf|safecrlf|eol)" 2>nul || echo    No local core settings found
echo    Global settings:
git config --global --get-regexp "core\.(autocrlf|safecrlf|eol)" 2>nul || echo    No global core settings found

echo.
echo 🎉 Git line ending configuration completed!
echo.
echo 📚 Next steps:
echo    1. Ensure you have a proper .gitattributes file
echo    2. Test by adding files: git add .
echo    3. Check for reduced line ending warnings
echo.
echo 💡 For more information, see GIT_LINE_ENDING_FIX.md

pause
