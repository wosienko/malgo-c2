# Projects

This page is available only to operators. By default it shows all future, current and projects ended within last 14 days. One may toggle the view to show all projects assigned to the operator.

Opening this page also establishes a Websocket connection.

Upon clicking, one will see on the left a list of all sessions and at the top a reminder of the time frame and current project.
It will also send information to the WS server to subscribe to the project's events.

## Session list

The session list shows all sessions for the current project. It is updated in real time as new sessions are created or updated.

Websocket events for the session list:

- 'new-session'
- 'session-heartbeat-updated'
- 'renamed-session'
