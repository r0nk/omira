# Omira

Omira is a personal task manager based on the unix philosophy.
<p align="center">
  <img src="https://raw.githubusercontent.com/r0nk/omira/main/images/usage.png">
</p>

## Installation

```bash
go install github.com/r0nk/omira@latest
```

## Usage

```bash

# add a dance task with a time estimate of 15 minutes, due a week from now
omira add -t 15 dance

# schedule tasks to be done today, up to 4 hours worth
omira schedule -w 4

# pretty print the tasks to be done today
omira status

# print the tasks in a scriptable format
omira task
```
## Documentation
[Getting started guide](https://github.com/r0nk/omira/blob/main/doc/getting_started.md)
## Contributing
The project is currently still in alpha, but pull requests are welcome.

## License
[gpl3](https://www.gnu.org/licenses/gpl-3.0.en.html)
