# worktimer

`worktimer` records active hours in X11 to make writing time sheets less of a pain in the ass.

## Building

You need to have the X11 libs installed.

```shell
make
```

## Installing

On Linux with X11 only:

```shell
make
sudo make install
```

## Running

Starting manually:
```
systemctl --user daemon-reload
systemctl --user enable --now worktimer
```

Alternatively you can use the makefile:
```
make run-as-service
```

## Commands

Use `worktimer help` to view the help. Use `worktimer help [COMMAND]` to view the help for any command.

These are the most commonly used commands:
```shell
worktimer start # Starts recording times
worktimer stop # Stops recording times
worktimer note "Some text" # Records a note in the current time slice so you know what you did at that time
worktimer status # Prints the status of the timer
```

## How it works

When the timer is started, it records times to `${HOME}/.config/worktimer` in JSON format.
The times are recorded in the form of time slices:

```json5
[
  {
    "Started": "2022-05-09T20:09:58.753476314+02:00",
    "Ended": "2022-05-09T20:10:02.008421633+02:00",
    "Duration": 3254945309,
    "Notes": null,
    "StartedBy": "Manual start",
    "EndedBy": "X11 idle"
  },
  {
    "Started": "2022-05-09T20:10:36.579078677+02:00",
    "Ended": "2022-05-09T20:11:51.624152165+02:00",
    "Duration": 75045073468,
    "Notes": null,
    "StartedBy": "X11 activity",
    "EndedBy": "X11 idle"
  },
  {
    "Started": "2022-05-09T20:11:52.997820441+02:00",
    "Ended": "2022-05-09T20:11:55.20356786+02:00",
    "Duration": 2205747419,
    "Notes": [
      "Did a thing"
    ],
    "StartedBy": "Note added",
    "EndedBy": "Manual stop"
  }
]
```

Each time the timer is stopped by the user or the user is inactive for 5 minutes, a time slice is recorded.
Inactivity is interrupted by using X11 (i.e. moving your mouse or typing) or by adding notes.
The daemon writes the current times to disk every hour, or when it shuts down.