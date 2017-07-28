# Babel

Babel is a library and a CLI to help manage tasks like update, upload,
download and use translations from [Babel](http://i18n.ml.com).

## Library Usage

Loading a bundled .zip file with translations:

```go
import (
  "github.com/mercadolibre/go-meli-toolkit/babel"
  "golang.org/x/text/language"
)

...
// "./conf/i18n/all.zip" is our mesage bundle download from babel API.
babel.Load("./conf/i18n/all.zip")

locale := language.MustParse("es-MX")
babel.Tr(locale, "Translation Key")
babel.Trn(locale, count, "Singular Key", "Plural Key")
```

Loading .po translations from an asset directory:

```go
import (
  "github.com/mercadolibre/go-meli-toolkit/babel"
  "golang.org/x/text/language"
)

...
// "./conf/i18n" is our asset directory with .po files downloaded from babel API.
babel.LoadDir("./conf/i18n")

locale := language.MustParse("es-MX")
babel.Tr(locale, "Translation Key")
babel.Trn(locale, count, "Singular Key", "Plural Key")

```


## Installing the CLI

```
$ go get github.com/mercadolibre/go-meli-toolkit/babel/babel
$ babel
babel is command line tool to help manage tasks like update, upload and download translations from Babel.

Usage:
  babel [command]

Available Commands:
  download    Download the message bundle.
  help        Help about any command
  scan        Scan project files.
  upload      Upload the message files.

Flags:
      --app string       github repository name
      --bundle string    message bundle with all translations (default "./conf/all.zip")
  -h, --help             help for babel
      --project string   Babel project name
      --source string    source filename (default "./conf/source.po")

Use "babel [command] --help" for more information about a command.
```


## Scanning your code for translations keys

Just run `babel scan` on your project root.
This will create a file called `./conf/source.po` which contains all the found translation keys.

> You can customize the output filename with the `--source` flag.


## Upload messages to Babel

Run `babel upload` and start translate your messages inside Babel.

> You must pass the `--project` and `--app` values.


## Download the translation bundle from Babel

Run `babel download`.

> You must pass the `--project` and `--app` values.
>
> You can customize the bundle location with the `--bundle` flag.

## Configuring your project

To make things easier, you can create a file `.babel.yaml` in your project root with the flag values,
or set env vars like `BABEL_PROJECT` and `BABEL_APIKEY`.

```
$ cat .babel.yaml
project: project-name
app: fury_my-app-name
source: ./assets/i18n/messages.po
```

## Still need help?
Contact me: jairo.dasilva@mercadolibre.cl
