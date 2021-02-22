# action-update-twirp

## This is not endorsed by or associated with GitHub, Dependabot, etc.

This action is plumbing to expose the https://github.com/thepwagner/action-update `Updater` interface through a [Twirp](https://github.com/twitchtv/twirp) service in a pluggable docker container.

This allows dependency parsing and updating logic to be implemented in the same language as the dependency manager, which tends to be productive and more approachable for the ecosystem's users.

This is refreshing an earlier concept, https://github.com/thepwagner/dependagot , with https://github.com/thepwagner/action-update .

## Implementations

- WIP: https://github.com/thepwagner/action-update-twirp-gradle
