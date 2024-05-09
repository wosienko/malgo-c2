# Operator Panel Overview

## Introduction

The Operator Panel is a web application designed to work mainly on desktops (although it may work in a reduced mode on mobile devices). The main purpose is to manage implants running on client devices.

The app is designed to be used by multiple operators, each with their own account. Users may be assigned to roles, which determine what they can do within the application.

The app is also designed to enable execution of multiple projects. This way there is no need to completely scratch the server when a new project is started.

## Roles

The application has two roles:

- Admin: This role allows for the management of users and projects. Admins may create new users, assign roles, reset passwords. They may also create new projects, assign users to projects, and delete projects.
- Operator: This role allows for taking part in projects. They may interact with sessions, run commands, etc.

There may be situations where a user has both roles or neither.

## Projects

Projects have been incorporated to allow for the separation of data. This is a logical separation of various implants and enable conduction of multiple operations at the same time.

Every single implant needs to have a project assigned to it. This is done during the implant creation process.

## Startup

During the startup, the application will analyse whether it has started before or not. If not, it will perform basic setup and start the operator panel. The config consists of the following:

- Create 'Admin' role.
- Create 'Operator' role.
- Create a default admin account based on the environment variables.

It should be noted that the default admin account may have all the data changed. Environment variables are only used during the first startup.

!!! warning
    There is no recovering the server if all administrators are deleted. Make sure to have at least one admin account.

## User Session

User session is handled with the help of 'lucia' library. By default the session lasts for 30 minutes.

## Validation

Input validation is done using 'zod' library. All validations in use may be found in `src/lib/validationSchemas.ts`.

## Components

Components are stored in a seperate folder. It may be found in `src/lib/components`. Down there, the structure is divided into general components (inputs, modals, etc.) and components tailored for a particular page. The latter ones are stored in `src/lib/components/custom`.

## Database

Database-related files are stored in `src/lib/db`. The database is managed with the help of 'drizzle' library. Schemas are stored in `src/lib/db/schemas`.

!!! info
    To add another schema, remember to add it to `src/lib/db/db.server.ts`. This way TypeScript will know about the new schema.

    ```ts
    import * as usersSchema from './schema/users';
    import * as sessionsSchema from './schema/sessions';
    import * as configSchema from './schema/server-settings';
    import * as projectsSchema from './schema/projects';
    import * as c2SessionsSchema from './schema/c2_sessions';
    import * as c2CommandsSchema from './schema/c2_commands';
    import * as c2ResultsSchema from './schema/c2_results';
    import * as OutboxSchema from './schema/outbox';
    import { drizzle } from 'drizzle-orm/postgres-js';
    import postgres from 'postgres';
    import { DATABASE_URL } from '$env/static/private';

    const client = postgres(DATABASE_URL);
    export const db = drizzle(client, {
        schema: {
            ...usersSchema,
            ...sessionsSchema,
            ...configSchema,
            ...projectsSchema,
            ...c2SessionsSchema,
            ...c2CommandsSchema,
            ...c2ResultsSchema,
            ...OutboxSchema
        }
    });
    ```

    It's important that you don't use SvelteKit imports with '@' sign. Instead, use relative imports (important for drizzle).

## Services

To slightly abstract database operations, services have been created. They may be found in `src/lib/services`.

## Stores

Stores are used to manage global state. They may be found in `src/lib/stores`. Thus far, there is only one store for Websocket connection.

### Websocket

The store is used to manage the websocket connection. It is used to send and receive messages from the server. The store is created in `src/lib/stores/websocket.ts`.
