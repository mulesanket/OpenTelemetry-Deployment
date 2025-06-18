#!/bin/bash

# Step 1: Create a new branch and orphan the history
git checkout --orphan new-main

# Step 2: Add everything from scratch
git add .
git commit -m "Initial commit - imported project"

# Step 3: Delete old branch and rename
git branch -D main
git branch -m main

# Step 4: Force push to your repo
git push -f origin main

