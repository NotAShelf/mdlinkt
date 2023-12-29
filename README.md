# 🚨 mdlinkt

A CLI tool for detecting dead or inaccessible links in markdown files.

## Usage

To check a markdown file for dead or inaccessible links

```console
mdlinkt -file your-file.md
```

This will return each tested link and a summary message in an example markdown file containing 1 valid and 1 invalid file.

```console
mdlinkt -file test.md
2023/12/12 12:00:00 ERROR Invalid link: https://probablynotvalid.com
2023/12/12 12:00:00 INFO Summary: 1 valid links, 1 invalid links
```

Should you want more details on the links that are scanned, you may pass the `-verbose` flag.

```console
2023/12/12 12:00:00 INFO https://github.com: 200
2023/12/12 12:00:00 ERROR Invalid link: https://probablynotvalid.com
2023/12/12 12:00:00 INFO Summary: 1 valid links, 1 invalid links
```

In case of invalid links, the program will exit with exit code 1, making it perfect for usage in
GitHub actions or other pipelines.

### Performance

For a file containing **10,000 invalid links**, the **Hyperfine** benchmark is as follows.

| Command                           |    Mean [ms] | Min [ms] | Max [ms] | Relative |
| :-------------------------------- | -----------: | -------: | -------: | -------: |
| `mdlinkt -verbose -file links.md` | 822.2 ± 22.6 |    787.4 |    959.8 |     1.00 |

The test has been conducted on a **Ryzen 5 3600X**, at a **95 ± 5 MB/s** bandwitdh speed.

## Hacking

A nix flake is provided. Use `direnv allow` or `nix develop` to enter the development shell.

### Contributing

PRs are always welcome.

## License

**mdlinkt** is licensed under the GPLv3. See [LICENSE](LICENSE) for more details.
