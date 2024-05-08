# Settings

Settings for each session.

## Table Entry

To make a generic and reusable table entry component, sending data is slightly more complex.

Firstly, `createSendType` function is used to make a new JS object with all the required fields (including one dynamic one).

!!! warning
    It wouldn't be necessary if not for the stupid TypeScript.

```ts
function createSendType<T extends InputProps>(
    props: T,
    type: string,
    sessionId: string,
    value: string
): { type: string; sessionId: string; [key: string]: string } {
    return {
        type,
        sessionId,
        [props.value_key]: value // Using value_key to set the dynamic property
    };
}
```

Next, we create a disposable type from which we can extract the value key.

```ts
let disposable: InputProps = { key, value, ws_type, value_key };
```

Finally, we send the data to the server.

```ts
websocketStore.ws?.send(
    JSON.stringify(createSendType(disposable, ws_type, get(page).params.sessionid, editableValue))
);
```
