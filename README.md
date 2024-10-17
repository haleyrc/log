# `log`

A set of custom `slog` handlers suitable for production use and testing.

---

## Cutting a release

Releases for this package are managed by the `goreleaser`[^1] utility which manages creating a GitHub release with a changelog, etc.
To create a new release, first make sure your local repo is in the desired state, and then tag the current commit with the next meaningful tag e.g.:

```sh
git tag -a v0.0.1 -m "v0.0.1"
```

Next, ensure all your commits and tags are pushed to GitHub:

```sh
git push --follow-tags
```

If you skip the previous step, the release notes will not function correctly and you will run into issues with duplicate tags the next time you try to push.

Finally, you can spin up `goreleaser` with:

```sh
make release
```

> **N.B.:** If you see an error about a missing GitHub token, you will need to create a `.env` file containing a token e.g.:
> ```
> set -a
> GITHUB_TOKEN=<YOUR-TOKEN-HERE>
> ```

[^1]: https://goreleaser.com/
