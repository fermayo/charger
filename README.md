# fermayo/charger

**NOTE: WORK IN PROGRESS**

Service aggregation framework based on Docker where applications implement their own HTTP API, CLI and web UI/UX.

# Requirements

* Docker 1.12+ in Swarm mode


# Usage

Building the images

	make all

or

	make cli     # Builds fermayo/charger-cli
	make server  # Builds fermayo/charger-server


Run the charger server

	docker run --name chargerd -p 9000:9000 -d fermayo/charger-server


## Deploy an app

TBD


# Interfaces

## HTTP interface

TBD

## CLI interface

	alias charger="docker run -it --link chargerd fermayo/charger-cli"
	charger [<command>]


## Web interface

TBD


# Charger service reference

## Creating a new charger service

TBD

