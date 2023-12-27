# Discord Drive

Discord Drive is a Go project that allows you to use a Discord channel as a storage drive.

## Building the Project

To build the project, use the Go build command:

```bash
go build dssd.go

Setting Up
Before running the project, you need to create a .env file with the following variables:

TOKEN=YOUR-DISCORD-BOT-TOKEN
CHANNELID=STORAGE-CHANNEL-ID
CHUNKSIZE=10000000



Replace YOUR-DISCORD-BOT-TOKEN with your Discord bot token, STORAGE-CHANNEL-ID with the ID of the Discord channel you want to use for storage, and CHUNKSIZE with the size of the chunks you want to use for file storage.

Running the Project
After setting up the .env file, you can run the project by executing the built binary:

./dssd


Project Structure
The project is organized into several packages:

discordutil: Contains utilities for interacting with Discord, including file download and upload functions.
handler: Contains HTTP handlers for different routes, including authentication, file download and upload, and file listing.
middleware: Contains middleware functions for HTTP requests.
storage: Contains functions for interacting with the storage system.
Views
The view directory contains the HTML files for the project's web interface.

```