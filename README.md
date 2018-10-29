# PACE Workshop Builder

## Hugo + Builder = PA Workshop

1. Download the correct `pace` CLI binary from the releases tab.

1. Setup a `config.json`. The format should follow the `sampleConfig.json`.

1. Run `pace build`. Notice the new `workshopGen` folder. This contains your new workshop.

1. [Install the Hugo CLI](https://gohugo.io/getting-started/installing/)

1. Change directory `cd` into workshopGen.

1. Run `hugo serve` to run the site locally on port `1313`

1. CF Push! `cf push workshop -m 64M -p public -b https://github.com/cloudfoundry/staticfile-buildpack.git`