# githubhop

Uses [GH Archive](https://gharchive.org) to create Timehop for GitHub.

Based on [githop](https://github.com/neonichu/githop) but written in go and uses https://gharchive.org gzipped archives instead of Big Query.

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
