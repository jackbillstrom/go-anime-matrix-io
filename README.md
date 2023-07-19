# go-anime-matrix-io

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/jackbillstrom)
<img src="./placeholder.gif" alt="Banner showing the Go Anime Matrix IO application in dark and light mode" width="1200px" height="400px">

<img src="./out.gif" alt="A generated gif-image that displays on your Anime Matrix display" width="138px" height="64px">

This project, `go-anime-matrix-io`, is a Go-based utility that fetches system performance data (CPU temperature and fan speeds) and displays it on an ASUS Anime Matrix display. It works by retrieving system sensor data at regular intervals, and generates a GIF that is displayed on the Anime Matrix display. This utility has been tested to work on systems running ASUS hardware with an Anime Matrix display and the `asusctl` utility installed.

It is using Fyne as GUI framework, and the [PixelMix font](https://andrewtyler.gumroad.com/) for the text that renders into the `out.gif`, which is being displayed on the Anime Matrix display.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

In order to run this project, you must have the following installed:

- Go (version 1.16 or later)
- ASUS system with Anime Matrix display
- `asusctl` utility from [asus-linux.org](https://asus-linux.org/)
- `lm-sensors` utility

You can install `asusctl` and `lm-sensors` using the package manager for your system. For Ubuntu-based systems, you would use:

```sh
sudo apt install lm-sensors
```

`asusctl` was installed on my Ubuntu machine while following their build instructions. You can find the instructions for installing `asusctl` on their [GitLab page](https://gitlab.com/asus-linux/asusctl#building)

### Installing

To get a copy of this project, you can clone it from GitHub using the following command:

```sh
git clone https://github.com/jackbillstrom/go-anime-matrix-io.git
```

After cloning the project, navigate into the project directory:

```sh
cd go-anime-matrix-io
```

You can then build the project using Go:

```sh
go build
```

### Usage

Once you have built the project, you can run it using:

```sh
./go-anime-matrix-io
```

When the program is running, it will update the Anime Matrix display with the system's CPU temperature and fan speed every 1 seconds. The program runs indefinitely until it is stopped. To stop the program, you can use `Ctrl+C` in the terminal.

### What to do if the program crashes

The program is designed to handle crashes gracefully. If the program crashes or is forcefully terminated, it will automatically clear and stop usage of the anime matrix.
If the screen does not turn off after the program has crashed, you can run the following command to turn it off:

```sh
asusctl anime -e false
```

### TODO

- :broom: Refactor code
- Add support for other systems (Windows)
- Add support for other Anime Matrix displays
- Better error handling
- Better layout for the matrix display
- Multiple themes

## Contributing

We welcome contributions from the community. If you would like to contribute, please fork the repository and submit a pull request.

## License

This project, **except the [PixelMix font](https://andrewtyler.gumroad.com/)** is licensed under the MIT License — See the [LICENSE.md](LICENSE.md) file for details.

_Credits to [Andrew Tyler](https://andrewtyler.gumroad.com/) for the font [PixelMix - Personal use](https://www.dafont.com/pixelmix.font)_

## Acknowledgments

Thank you to everyone who has contributed to the development and testing of this project. Your time and efforts are greatly appreciated.
