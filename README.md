# PACE Workshop Builder

[![asciicast](https://asciinema.org/a/213808.svg)](https://asciinema.org/a/213808?autoplay=1)

## Quick Start

1. Download the correct `pace` CLI binary from the releases tab. 
    - *MAC OS Users Optional:* If you have the `brew tap pivotal/tap` installed you can install the `pace` CLI with `brew install pace-cli`

1. Run `pace init`.

1. Edit the `config.json`. The format should follow the `sampleConfig.json`.

1. Run `pace build`. Notice the new `workshopGen` folder. This contains your new workshop.

1. ***Optional** Run `pace serve` to view your workshop. View local running site at http://localhost:1313

1. Run `pace push`. This will push your workshop to Pivotal Web Services.

## Notes

1. Content is pulled from the [pace-workshop-content](https://github.com/Pivotal-Field-Engineering/pace-workshop-content) github repo. Feel free to add any content there that you can then use to build a workshop with `pace build`

1. While `pace` will build a generic homepage for your workshop you can setup a custom one by supplying a markdown file via the `workshopHomepage` field in the `config.json` file. This is not required.

1. `pace push` will automatically generate a random hostname for your workshop but you can specifiy a custom one with the use of the `workshopHostname` attribute inside your `config.json`.
