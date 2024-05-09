# Commands

Page for sending new commands to the implant.

At the bottom, operator sees last command sent to the implant.

To the left of the command window, key-values are stored. It is a simple store for useful data.

!!! warning
    Right now, newestCommand may have one of three values:
    - `null` - loading data from the server
    - `undefined` - no command at all
    - `CommandType` - command could have been retrieved from the server

## Key-Value Store

Available events to listen to:

- session-key-value-modified
- session-key-value-deleted

Creating new key:

```javascript
const createNewKey = async (): Promise<boolean> => {
    websocketStore.ws?.send(
        JSON.stringify({
            type: 'session-new-key',
            sessionId: get(page).params.sessionid,
            key: newKey.key,
            value: newKey.value
        })
    );
    newKey = { key: '', value: '' };
    return true;
};
```

## Command

Available events to listen to:

- new-commands
- command-status-updated
- command-result-retrieved
- result-chunk-inserted

To send a command, websocket store `sendCommand` is executed.
