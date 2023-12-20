
# Discord Drive

Discord Drive is a Go project that allows you to use a Discord channel as a storage drive.

## Building the Project

To build the project, use the Go build command:

```bash
go build dssd.go
```

#Setting Up

Before running the project, you need to create a ```.env``` file with the following variables:

```bash
TOKEN=YOUR-DISCORD-BOT-TOKEN
CHANNELID=STORAGE-CHANNEL-ID
CHUNKSIZE=10000000
```

Replace ```YOUR-DISCORD-BOT-TOKEN``` with your Discord bot token, ```STORAGE-CHANNEL-ID``` with the ID of the Discord channel you want to use for storage, and ```CHUNKSIZE``` with the size of the chunks you want to use for file storage. ```10000000 = 9 mb```

## Running the Project

After setting up the ```.env``` file, you can run the project by executing the built binary:

```bash
./dssd
```

## API Reference

#### Get Index Page

```http
  GET /
```

No parameters required.

#### Get File List

```http
  GET /files
```

No parameters required.

#### Download File

```http
  GET /download
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`    | `string` | **Required**. Name of file        |


#### Upload File

```http
  POST /upload
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `file`    | `binary` | **Required**. File to be uploaded |

#### Delete File

```http
  DELETE /files/:name
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`    | `string` | **Required**. Name of file        |

#### Rename File

```http
  PUT /files/:oldFileName/:newFileName
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `oldFileName` | `string` | **Required**. Current name of the file |
| `newFileName` | `string` | **Required**. New name for the file |

#### Share File

```http
  GET /share
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`    | `string` | **Required**. Name of file        |
