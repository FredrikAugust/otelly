# Otelly

A TUI for viewing and debugging with OTEL signals locally.

![a screenshot showing otelly in action](./assets/screenshot.png)

## Development

This project uses [Taskfile.dev](https://taskfile.dev) to simplify running commands.

If you want to install all deps on MacOS:

```bash
brew install krzko/tap/otelgen
brew install go-task
```

And then the different tasks:

```bash
task run               # Builds and runs the app
task debug             # Builds and runs the app for attaching debugger from neovim
task test              # Run tests
task ingest-dummy-data # Uses otelgen to send some data for testing
task logs              # Tail (follow) logs
```

Since the TUI takes up the main window, we write logs to `debug.log`.

During testing, to avoid having to re-seed my development environment all the time,
the application state is persisted to `local.db`. You can remove this to clear the
state or use the application "clear state" feature (`ctrl+l` from table view).
This will be changed before the first release so you can choose if you want
persistence or not.

> [!NOTE]  
> This projects is partially intended as a learning project, and as such use of
> AI is kept to a minimum.
> I consult it from time to time to brainstorm or fix my SQL queries, but an
> attempt is made to keep the
> project "hand made in Denmark" (other countries are encounraged to contribute).
