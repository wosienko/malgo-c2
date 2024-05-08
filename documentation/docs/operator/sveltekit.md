# SvelteKit

As mentioned, Operator Panel has been created with SvelteKit. This document serves as a guide to help you understand how to work with SvelteKit. It should be noted that it is not a comprehensive guide to SvelteKit, but rather a guide to help you understand how to work with SvelteKit in the context of Operator Panel. For a more comprehensive guide to SvelteKit, please refer to the [official SvelteKit documentation](https://kit.svelte.dev/docs/introduction).

## Directory Structure

```plaintext
migrations/         # Drizzle migrations
src/                # SvelteKit source code
    lib/            # Shared code, e.g. components, stores, services, also server-side only code
    routes/         # Pages, endpoints, and layouts
    app.d.ts        # TypeScript types for the app
    app.html        # HTML template for the app
    hooks.server.ts # Server hooks - runs before any server route
    main.css
static/             # Static files
tests/              # Playwright tests
...                 # Lovely JS/TS configs...
```

## Routing

SvelteKit routing is filesystem based. It means that directories specify routes and files specify when and how code is executed.
The following files are used for routing:

- `+page.svelte` - the actual page content of the particular route.
- `+layout.svelte` - Common layout for this routes and all subroutes. Think of it as something consistent across subroutes, e.g. header and footer.
- `+page.ts` - data load source for the page. This runs once the page is accessed and serves data to `page.svelte`. May run both on the server and the client.
- `+layout.ts` - data load source for the layout. This runs once the layout is accessed and serves data to `layout.svelte`. May run both on the server and the client.
- `+page.server.ts` - data load source for the page. This runs once the page is accessed and serves data to `page.svelte`. Runs only on the server. Thus, may access server-only resources.
- `+layout.server.ts` - data load source for the layout. This runs once the layout is accessed and serves data to `layout.svelte`. Runs only on the server. Thus, may access server-only resources.

For API endpoints, there is only one file:

- `+server.ts` - endpoint logic. Has function for each HTTP method.

## Important Files & File Structures

### +page.svelte

It's as simple as it gets.

```js
<script lang="ts">
    // Your page logic here
    let num = $state(0);
</script>
<h1>{num}</h1>
```

### +layout.svelte

It's as simple as it gets, but with a twist - you need to add a slot for the page content.

```js
<script lang="ts">
    // Your layout logic here
    let num = $state(0);
</script>

<header>
    <h1>Operator Panel {num}</h1>
</header>

<main>
    <slot/>
</main>
```

### +page.ts

This is where you load data for the page. It's an async function that returns an object. The object is then passed to the page component as props.

```js
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    return {
        user: await fetch('/api/self').then((r) => r.json())
    };
};
```

!!! info
    You do not need to do await here. Moreover, it is not always recommended. SvelteKit does all load functions in parallel, so you can save time by not waiting for one request to finish before starting another.

### +layout.ts

Similar thing as with the page, but for the layout. Major difference is the TypeScript type.

```js
import type { LayoutLoad } from './$types';
import { uuidSchema } from '$lib/validationSchemas';
import { json } from '@sveltejs/kit';

export const load: LayoutLoad = async ({ params, data, fetch }) => {
    const { id } = params;

    const projectId = uuidSchema.safeParse(id);
    if (!projectId.success) {
        return json(projectId.error.issues, { status: 400 });
    }

    return {
        project: data.project,
        sessions: fetch(`/api/projects/${projectId.data}/sessions`).then((r) => r.json())
    };
};
```

### +page.server.ts

Same thing as with the page, but for the server only. This is where you can access server-only resources. Only major difference apart from that is the TypeScript type.

It can also include Actions - functions that handle other HTTP methods than GET.

```js
import type { Actions } from './$types';
import { loginSchema } from '$lib/validationSchemas';
import { fail, redirect } from '@sveltejs/kit';
import type { ZodIssue } from 'zod';
import { validatePasswordForEmail } from '$lib/services/user-service';
import { createSession } from '$lib/services/session-service';

export const actions = {
    default: async ({ request, cookies }) => {
        const data = await request.formData();
        const userLogin = Object.fromEntries(data.entries());

        const result = loginSchema.safeParse(userLogin);
        if (!result.success) {
            return fail(400, {
                issues: result.error.issues
            });
        }

        const userId = await validatePasswordForEmail(result.data.email, result.data.password);
        if (userId === '') {
            const issue: ZodIssue = {
                code: 'custom',
                message: 'Incorrect email or password',
                path: ['Login']
            };

            return fail(400, {
                issues: [issue]
            });
        }

        await createSession(userId, cookies);

        redirect(302, '/home');
    }
} satisfies Actions;

```

### +layout.server.ts

Same thing as with the layout, but for the server only. This is where you can access server-only resources. Only major difference apart from that is the TypeScript type.

```js
import type { LayoutServerLoad } from './$types';
import { isUserAdmin, isUserOperator } from '$lib/services/roles-service';

export const load: LayoutServerLoad = async ({ locals }) => {
    const userId = locals.user?.id ?? '00000000-0000-0000-0000-000000000000';
    return {
        loggedIn: !!locals.user && !!locals.session,
        isAdmin: await isUserAdmin(userId),
        isOperator: await isUserOperator(userId)
    };
};
```

### +server.ts

Just a collection of HTTP methods for the endpoint.

```js
import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { getBasicInformationForId, updateUserData } from '$lib/services/user-service';
import { updateUserSchema } from '$lib/validationSchemas';

export const GET: RequestHandler = async ({ locals }) => {
    const userId = locals.user!.id;

    return json(await getBasicInformationForId(userId));
};

export const PATCH: RequestHandler = async ({ locals, request }) => {
    const userId = locals.user!.id;

    const body = await request.json();

    const result = updateUserSchema.safeParse(body);
    if (!result.success) {
        return json(result.error.issues, { status: 400 });
    }

    const message = await updateUserData(
        userId,
        result.data.name,
        result.data.surname,
        result.data.email
    );
    if (message) {
        return json({ message }, { status: 400 });
    }

    return json({ id: userId });
};
```

### hooks.server.ts

This is where you can run code before any server route.
It's a good place to put code that you want to run before any server route (e.g. authorization checks).

```js
import { error, type Handle, type HandleServerError, redirect } from '@sveltejs/kit';
import { createDefaultAdminAndRoles } from '$lib/services/user-service';
import { deleteExpiredSessions, validateSession } from '$lib/services/session-service';

export const handle: Handle = async ({ event, resolve }) => {
    await deleteExpiredSessions();
    await validateSession(event);

    if (event.route.id?.startsWith('/(authorized)/')) {
        const nextRoute = event.route.id.replace('/(authorized)/', '');
        if (!event.locals.session && !event.locals.user) return redirect(303, '/login');
    }

    return resolve(event);
};

export const handleError: HandleServerError = async ({ error, status }) => {
    if (status === 401 || status === 403 || status === 404) {
        return;
    }
    console.error(error);
};

// on startup code - insert default roles and admin user
await createDefaultAdminAndRoles();
```
