#!/bin/bash

echo "Starting post-install process..."

echo "Removing com.apple.quarantine..."
sudo xattr -d com.apple.quarantine /Applications/Cordis.app/ || true

echo "Setting chmod for unfold..."
sudo chmod -R 777 /Applications/Cordis.app/ || true

echo "Starting unfold..."
sudo /Applications/Cordis.app/Contents/MacOS/unfold ensure

echo "Setting chmod for app..."
sudo chmod -R 755 /Applications/Cordis.app/ || true

echo "Setting chmod for user data..."
sudo chmod -R 777 ~/Library/Application\ Support/Koishi/Desktop/ || true

echo "Setting chown for user data..."
sudo chown -R ${USER}:staff ~/Library/Application\ Support/Koishi/Desktop/ || true

echo "Post-install process finished."
