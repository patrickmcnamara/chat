# chat

![latest version](https://img.shields.io/github/v/tag/patrickmcnamara/chat?label=latest%20version)
![last commit](https://img.shields.io/github/last-commit/patrickmcnamara/chat)
![top language](https://img.shields.io/github/languages/top/patrickmcnamara/chat)
![licence](https://img.shields.io/github/license/patrickmcnamara/chat?label=licence)

An end-to-end encrypted chat application.
I wrote this to learn about end-to-end encryption and networking.
Don't use this.
It does run however.

## Installation

Run `go get -u github.com/patrickmcnamara/chat/...`.

## Usage

Run `chat-client`, assuming your $GOPATH and $PATH are set up correctly, to run the client.
Or run `chat-server` to run a server instead.

## Client configuration

The config files are located in:
- Unix systems, `$XDG_CONFIG_HOME` if non-empty, else `$HOME/.config`.
- Darwin, `$HOME/Library/Application Support`.
- Windows, `%AppData%`.
- Plan 9, `$home/lib`.

There are three config files with this directory:
- the `config` file just contains the server and port of the chat server.
- the `contacts` file is a JSON mapping from a contact name to a 32-length byte array of the contact's public key.
- the `profile` file contains a public and private key pair and is automatically be generated.

There are examples of each of these in the `misc/example-config-files` directory.

## Server configuration

There is no configuration for the server.
It uses port `6969`.
You might need to allow through your firewall or NAT setup.

## Miscellaneous

To talk to someone with the client you must first select the contact using `/msg CONTACT_NAME`.
Use quotes if the name contains a space.

## Licence

This project is licenced under the European Union Public Licence v1.2.
