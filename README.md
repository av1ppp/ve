# Ve

**Ve** is a minimal versioning tool for applications.

It stores the version in the format: `YYMMDD.micro`.
Where `micro` - incrementing number for changes made on the same day (starting from 1).

## Installation
You can install the **ve** using the command:

```sh
go install github.com/av1ppp/ve/cmd/ve@latest
```

Ensure you have Go installed and properly configured in your environment.

## Commands
- `init` - initialize **ve** in the current directory by creating a `VERSION` file
- `incr` - increment the version number in the `VERSION` file
- `get` - get the current version number from the `VERSION` file
- `help` - print the help message

## License
This project is licensed under the MIT License.
