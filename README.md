# fermayo/charger

**NOTE: WORK IN PROGRESS**

Service aggregation framework based on Docker where applications implement their own HTTP API, CLI and web UI/UX.

# Requirements

* Docker 1.12+ in Swarm mode


# Usage

Building the images

	make all

or

	make mac-cli  # Builds the `build/charger` binary for macOS
	make server   # Builds the `fermayo/charger-server` docker image


Run the charger server

	docker run --name chargerd -p 9000:9000 -d fermayo/charger-server


## Deploy an app

TBD


# Interfaces

## HTTP interface

TBD

## CLI interface

	charger [<command>]


## Web interface

TBD


# Charger service reference

## Creating a new charger service

TBD

