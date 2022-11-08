#!/bin/bash

echo "Starting post-install process..."

echo "Removing com.apple.quarantine..."
sudo xattr -d com.apple.quarantine /Applications/Koishi.app/ || true

echo "Setting chmod for unfold..."
sudo chmod -R 777 /Applications/Koishi.app/ || true

echo "Starting unfold..."
sudo /Applications/Koishi.app/Contents/MacOS/unfold ensure

echo "Setting chmod for app..."
sudo chmod -R 755 /Applications/Koishi.app/ || true

echo "Setting chmod for user data..."
sudo chmod -R 777 ~/Library/Application\\ Support/Il\\ Harper/Koishi/ || true

echo "Setting chown for user data..."
sudo chown -R ${USER}:staff ~/Library/Application\\ Support/Il\\ Harper/Koishi/ || true

echo "Post-install process finished."
