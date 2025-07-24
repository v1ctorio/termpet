# TermPet

A terminal pet made to accompany you on your text adventures.

```
NAME:
   Termpet - Take care of your pet!

USAGE:
   Termpet [global options] [command [command options]]

COMMANDS:
   init
   greet     Greet your terminal pet!
   stat      View the statistics of your pet
   feed      Feed the pet
   coinflip  Make your pet flip a coin
   rps       Rock, paper, scissors
   dispatch  Dispatch a message to a service. Only Slack supported atm
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug     Print stack traces (default: false)
   --help, -h  show help
```

Termpet is someone who is always there to help you with small little things from the terminal.

From flipping a coin or playing a rock paper scissors to dispatching a message to Slack (more services to come).

Also, don't forget to take care of your terminal pet and feed them or they might get sick.



# Installation
## Termpet
You can install termpet downloading its release binaries from [the releases](https://github.com/v1ctorio/termpet/releases/latest), install it with the nix flake (`nix profile install github:v1ctorio/termpet`) or using go `go install github.com/v1ctorio/termpet@latest`
#### Cowsay
To use TermPet with it's default configuration, you need to have `cowsay` installed on your system. I recommend using [Neo-cowsay](https://github.com/Code-Hex/Neo-cowsay).

You can easily install it for windows using winget (`winget install neo-cowsay`) or using go `go install github.com/Code-Hex/Neo-cowsay/cmd/v2/cowsay@latest`

For Linux you also can install it using go or using nix (`nix profile install nixpkgs#neo-cowsay`). Also, the package `cowsay` is most likely in your distro repositories.


