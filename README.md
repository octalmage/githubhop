# githubhop [![Build Status](https://travis-ci.org/octalmage/githubhop.svg?branch=master)](https://travis-ci.org/octalmage/githubhop)

Uses [GH Archive](https://gharchive.org) to create Timehop for GitHub.

Based on [githop](https://github.com/neonichu/githop) but written in go and uses the https://gharchive.org gzipped archives instead of Big Query.

This project streams the events for every user from GH Archive and extracts the events relevant to the specified user. This is all done in memory and it usually takes around a minute to make it through a full day of events.

## Installation

### Homebrew

```
brew tap octalmage/githubhop
brew install githubhop
```

### Manual Install

```
go get -u github.com/octalmage/githubhop
```
## Usage

```bash
$ githubhop --date 2018-05-27
Fetching hours (24/24)   20s [====================================================================] 100%
At 2018-05-26 12:57am, you created the repository octalmage/githubhop
At 2018-05-26 12:57am, you created the branch master on octalmage/githubhop
At 2018-05-26 1:44pm, you pushed commits to octalmage/githubhop
At 2018-05-26 5:57pm, you watched the repo remeh/sizedwaitgroup
```

## License

MIT
