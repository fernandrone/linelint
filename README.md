# linelint

A linter that validates simple "newline" and "whitespace" rules in all sorts of files.

## Usage

This is a project in development. Use it at your own risk!

```console
go get github.com/fernandrone/linelint
```

Just run it and pass a list of file or directories as argument:

```console
linelint .
```

Or:

```console
linelint README.md LICENSE linter/config.go
```

## Configuration

Use a `.linelint.yml` in the same directory to fine-tune your settings. See the example file [.linelint.yml] for an up-to-date example:

```yaml
# 'true' will fix files
autofix: true

# list of paths to ignore, uses gitignore syntaxes
ignore:
  - .git/

rules:
  # checks if file ends in a newline character
  end-of-file:
    # set to true to enable this rule
    enable: true

    # set to true to disable autofix (if enabled globally)
    disable-autofix: false

    # will merge with global configuration
    ignore:
      - README.md

    # if true also checks if file ends in a single newline character
    single-new-line: true
```

## Rules

Right now it only supports a single rule, "End of File", which is enabled by default.

### EndOfFile

EndOfFileRule checks if the file ends in a newline character, or `\n`. You may find this rule useful if you dislike seeing these 🚫 symbols at the end of files on GitHub Pull Requests.

By default it also checks if it ends strictly in a single newline character. This behavior can be disabled by setting the `single-new-line` parameter to `false`.

```yaml
end-of-file:
  # set to true to enable this rule
  enable: true

  # set to true to disable autofix (if enabled globally)
  disable-autofix: false

  # will merge with global configuration
  ignore:
    - README.md

  # if true also checks if file ends in a single newline character
  single-new-line: true
```

## Docker Image

A public docker image exists at [docker.io/fernandrone/linelint](https://hub.docker.com/repository/docker/fernandrone/linelint). It is published manually and I at the moment there is no versioning and no guarantees of updates.

```console
docker run -it -v $(pwd):/data fernandrone/linelint /data
```

To add a configuration file, just share it with the root volume of the container:

```console
docker run -it -v $(pwd)/.linelint.yml:/.linelint.yml -v $(pwd):/data fernandrone/linelint /data
```
