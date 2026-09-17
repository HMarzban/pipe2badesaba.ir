# pipe2badesaba.ir

A Go experiment that fetches Badesaba calendar events and writes JSON and iCalendar (`.ics`) files.

## Status and scope

This is a historical utility. `main.go` currently selects Jalali years **1400–1403**; it does not automatically select the current year. The upstream API and generated data have not been revalidated by this documentation update.

## Run

The module declares Go 1.17. From the repository root, with a compatible Go toolchain and network access:

```sh
go mod download
mkdir -p dist
go run .
```

To build without contacting the calendar service:

```sh
go build -o /tmp/pipe2badesaba .
```

The generator makes requests to `https://badesaba.ir/api/site/getDataCalendar` and writes `dist/data-<year>.json` and `dist/events-<year>.ics`. Change the `getYears` list in `main.go` deliberately if generating another range.

## Limitations

There is no CLI or automated fixture suite yet. The original implementation does not report output-file write errors; verify that both files exist and inspect month/event counts after a run. Upstream availability and date-conversion correctness need verification before using these files as an authoritative calendar.

Related JavaScript project: [pipe2time.ir](https://github.com/HMarzban/pipe2time.ir).
