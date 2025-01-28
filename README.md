# GoStreamer

GoStreamer is a tool for streaming and processing video frames from a webcam or files. It allows you to configure the server connection, select source and target files or webcam, and receive processed output.

It is only possible to send and recieve files manually as of now. 

Some future planned fueatures are 
  - adjusting server setting via client UI
  - auto upload
## Setup

1. **Modify `config.json`** to set up connection settings and to enable/disable the webcam. Folders can be set here, but it is also available to change in the UI.

    ```json
    {
      "server": {
        "ip": "192.168.1.2",
        "port": 8050
      },
      "local": {
        "webcam": {
          "enable": false,
          "target": "-1"
        },
        "source_folder": "sources",
        "target_folder": "targets",
        "output_folder": "output",
        "last_source": "",
        "last_target": "",
        "last_swapped": ""
      }
    }
    ```

2. **Install dependencies** by running:

    ```sh
    go mod tidy
    ```

## Running the Application

2. **Run the main application**:

    ```sh
    go run main.go
    ```

    This will start the server on the port set in `config.json`.
    This will start the GoStreamer UI.

## Usage

### Using Files

1. **Select Source Folder**: Choose the folder containing source files.
2. **Select Target Folder**: Choose the folder containing target files.
3. **Select Output Folder**: Choose the folder where the output will be saved.
4. **Submit**: Click the "Submit" button to start processing.
5. **Get Swapped**: Click the "Get swapped" button to receive the processed files. (Will be replaced by auto upload in the future.

### Using Webcam

1. **Select Source Face**: Choose a file for the source face.
2. **Enter Webcam Target**: Enter the webcam target (default is usually 0).
3. **Submit**: Click the "Submit" button to start streaming and processing frames.

## Future Features

- **Webcam**:
  1. Select webcam source
  2. Select target face
  3. Send target file to server
  4. Feed webcam frames to server
  5. Receive swapped webcam frames from server

- **Files**:
  1. Select source and target faces
  2. Select local output path
  3. Send target and source files to server
  4. Receive swapped output file from server
