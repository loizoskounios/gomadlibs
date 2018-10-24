# GoMadLibs

A very simple implementation of MadLibs written in Go using only the standard library.

## Running

You can run the app using Docker:

```shell
docker run --rm -it loizoskounios/gomadlibs
```

This will fetch the `loizoskounios/gomadlibs` image from Docker Hub and start up the app. You can skip ahead to the [Usage section](#usage) if you use this approach.

## Building

### Docker Image

To build the Docker image yourself, you can use the provided `Dockerfile`. From this directory:

```shell
docker build --tag gomadlibs .
```

This will start a multi-stage build that will leave you with a Docker image named `gomadlibs` and an intermediate image that you can safely discard. The intermediate images are assigned a label and can be easily removed like so:

```shell
docker rmi $(docker images --filter label="stage=intermediate" --quiet)
```

### Saving the Binary and / or Cross-Compiling

If you need to (cross-)compile the app and have the binary easily accessible outside of a Docker container, you should follow one of the two subsections.

#### Using a Docker Builder

Bind-mount this directory inside a Docker container running a `golang` image, and run `go build`. The resulting binary will be stored in this directory. From this directory:

```shell
docker run --rm \
  --user $(id -u "$USER"):$(id -g "$USER") \
  --volume "$(pwd)":/gomadlibs \
  --workdir /gomadlibs \
  golang:1.11.1-stretch \
  go build
```

To cross-compile, instantiate the Docker container with the appropriate environment variables via means of the `--env` flag. The following example will generate a binary for amd64 Windows.

```shell
docker run --rm \
  --user $(id -u "$USER"):$(id -g "$USER") \
  --volume "$(pwd)":/gomadlibs \
  --workdir /gomadlibs \
  --env GOOS=windows \
  --env GOARCH=amd64 \
  golang:1.11.1-stretch \
  go build
```

**Note:** If you encounter the warning `go: disabling cache (/.cache/go-build) due to initialization failure: mkdir /.cache: permission denied` during a build, it is a known issue. It occurs when the UID of the user that is running the container does not exist inside the container. In such a case, Docker sets `$HOME` to `/` and the Go compiler cannot write to `/.cache` since the process is not running as root. More details [here](https://go-review.googlesource.com/c/go/+/122487).

#### Using Your Local Golang Development Environment

From this directory, run

```shell
go build
```

This will generate a binary in this same directory.

Cross-compilation is also possible by setting the `GOOS` and `GOARCH` environment variables. The following example will generate a binary for amd64 Windows.

```shell
GOOS=windows GOARCH=amd64 go build
```

## Usage

The application reads a template and asks the user for input via `stdin`.
The user's answers are then plugged into the template and printed on `stdout`. 

```text
$ ./gomadlibs -help
Usage of gomadlibs:
  -help
        show this help message
  -stories-dir string
        where the app will look for story templates (default "./stories")
  -verify-integrity
        verifies the integrity of the story template

Examples:
    gomadlibs story1.mdlb
    gomadlibs /home/user/stories/astory.mdlb
    gomadlibs -verify-integrity story2.mdlb
    gomadlibs -stories-dir mystories astory.mdlb
```

## Story Templates

Two example story templates are provided in the `./stories` directory.
Unless a path to a template is explicitly provided, the app will look for files ending in `.mdlb` in the default stories directory (`./stories`) and choose a random one if any are available.

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

#### Example

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
$ ./gomadlibs example.mdlb
Template at 'stories/example.mdlb' has 4 blanks and 4 descriptions

----------
Name of Person: Mary
Noun: something
Name of Person: Jane

Let's Have a Talk
----------
Hi, Mary. My name is Jane. So, Mary, I need to talk to you about something.
```
