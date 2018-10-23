# GoMadLibs

A very simple implementation of MadLibs written in Go using only the standard library.

## Building

Bind-mount this directory inside a Docker container running a `golang` image, and run `go build`. The resulting binary will be stored in this directory.

```shell
docker run --rm \
  --volume "$(pwd)":/gomadlibs \
  --user $(id -u "$USER"):$(id -g "$USER") \
  --workdir /gomadlibs \
  golang:1.11.1-stretch \
  go build
```

Additional environment variables can be provided to perform cross-platform builds. For example, this will generate a binary for amd64 Windows.

```shell
docker run --rm \
  --volume "$(pwd)":/gomadlibs \
  --user $(id -u "$USER"):$(id -g "$USER") \
  --workdir /gomadlibs \
  --env GOOS=windows \
  --env GOARCH=amd64 \
  golang:1.11.1-stretch \
  go build
```

**NOTE:** If you encounter the warning `go: disabling cache (/.cache/go-build) due to initialization failure: mkdir /.cache: permission denied` during a build, it is a known issue. It occurs when the UID of the user that is running the container does not exist inside the container. In such a case, Docker sets `$HOME` to `/` and the Go compiler cannot write to `/.cache` since the process is not running as root. More details [here](https://go-review.googlesource.com/c/go/+/122487).

## Running

```shell
go run main.go

# or

go build main.go
./main

# or

go install
gomadlibs # With $GOBIN in $PATH
```

## Story Templates

Two example story templates are provided in the `./stories` directory.
The application reads a template and asks the user for input via `stdin`.
The user's answers are then plugged into the template and printed on `stdout`. 
Unless a path to a template is explicitly provided, the app will look for files ending in `*.mdlb` in the default stories directory (`./stories`) and choose a random one if any are available.

### Template Format

Templates are split into 3 parts.

* Title
* Story
* Descriptions

Parts are separated by 5 hyphen symbols on a new line (`-----`).

#### Title

No special considerations.

#### Story

The blanks in the story are denoted by 5 underscore symbols (`_____`). There must be one description for every blank left in the story.

#### Descriptions

Line-separated, case-sensitive descriptions of what should go in each blank spot in the story. The second blank in the story is associated with the second line in the descriptions and so forth. Numbers should be used to differentiate between different adjectives, nouns, etc (e.g., `Noun 1`, `Noun 2`).

## Example

```text
Hi, John. My name is Jack. So, John, I need to talk to you about something.
```

Let's imagine that we want to convert the above snippet, titled `Let's Have a Talk`, into a template, replacing `John`, `Jack` and `something` with blanks. Note that `John` is mentioned twice.

```text
Let's Have a Talk
-----
Hi, _____. My name is _____. So, _____, I need to talk to you about _____.

-----
Name of Person 1
Name of Person 2
Name of Person 1
Noun 1
```

Note how `Name of Person 1` is used both as the first and third description, matching the first and third blanks of our story. When the user is asked to provided answers, they are only asked about unique descriptions.

We can save this in `./stories/example.mdlb` and load it by passing its name as the first argument to our application:

```text
$ go run main.go example.mdlb
Template at 'stories/example.mdlb' has 4 blanks and 4 descriptions

----------
Name of Person: Mary
Noun: something
Name of Person: Jane

Let's Have a Talk
----------
Hi, Mary. My name is Jane. So, Mary, I need to talk to you about something.
```
