# GitHub Activity CLI

A simple Go command-line tool to fetch and summarize recent GitHub activity for a given username based on roadmap.sh project [github-user-activity](https://roadmap.sh/projects/github-user-activity).

## Features
- Verify if a GitHub user exists.
- Fetch the user's recent public events from GitHub.
- Classify events into types: PushEvent, WatchEvent, CreateEvent, IssuesEvent, PullRequestEvent, ForkEvent.
- Summarize push events by repository and number of commits.

## Installation
Ensure you have Go installed (version 1.20+ recommended).

Clone this repository

Build the executable:
```go
go build -o github-activity
```

## Usage

Run the program with a GitHub username as an argument:

```go
./github-activity <username>
```

## Example:
```go
./github-activity kaikelfalcao
```

Expected output:
```go
The last 30 activities from kaikelfalcao
- Starred repo_name
- Pushed 3 commits to repo_name
- Created repo new_repo
- Forked repo another_repo
- opened a PR in repo_name
...
```