# Operator Panel

## Fullstack Framework

While the original idea for the Malgo Command & Control assumed that there will be a clear distinction between frontend and backend technologies (Vue + Python, HTMX + Go + SSR), it was later deemed impractical. After all, there is no intention of creating multiple frontend views. After numerous attempts at creating Malgo, it was decided that fullstack framework is the best choice for application this complex. This way, one is able to easily reuse logic (e.g. validation schemas).

### SvelteKit

There aren't as many fullstack frameworks as there are frontend libraries. The main contenders were rather obvious:

- [Next.js](https://nextjs.org/)
- [Nuxt](https://nuxt.com/)
- [SvelteKit](https://kit.svelte.dev/)

Why weren't Next.js or Nuxt used:

- Next.js uses React - unreadable, constantly changing framework that is unnecessarily complicated.
- Nuxt is a bit cluncy and has some quirks that make it slightly tedious when developing, although it's okay.

Moreover, SvelteKit has several advantages:

- It is by far the simplest framework of all.
- Routing is entirely directory-based.
- Order of execution and quirks of this framework are easy to understand.
- Compiles directly to JavaScript, HTML and CSS, with no additional runtime.
- [Documentation](https://kit.svelte.dev/docs/introduction) is on par with Vue's.

## Runtime

Sadly, the state of the JavaScript ecosystem is... peculiar. There are so many runtimes to choose from.

### Node.js

The default Node.js is a good and reliable choice. While development is fully compatible with this runtime, the final choice was not to use Node.js in a production version.

!!! info
    To develop this project, Node.js needs to be installed either way.

### Bun

That said, Bun is a faster runtime than Node.js. It outperforms Node.js in nearly any action, especially the download of the infamous `node_modules/`. It is highly encouraged to use Bun instead of Node.js during development, especially considering the fact that it is a drop-in replacement.

That said, Bun is a hardcoded runtime for production environment. Switching runtimes in production will be described in chapters directly related to the development of the Operators' Web Application.

## ORM

Object-Relational Mapping (ORM) is commonly used to abstract away database queries. Naturally, it was also used in this project.
This way, developers may easily manage schema, relations and migrations with no additional tools and gain syntax highlighting when adding new functionalities.

### Drizzle ORM

Until recently, there weren't many popular ORMs in JS/TS world. Thankfully, some time ago Drizzle ORM was created - a simple, zero dependency ORM. It was a pretty blatant choice:

- It is much faster than Prisma.
- Schema is created in JS/TS.
- Schema may be easily extended through raw SQL.
- Provides all possible means of querying the database:
      - Raw SQL.
      - SQL query builder.
      - Abstracted queries (FindOne, FindMany).
- Zero dependencies.
- Additional Drizzle Studio tool so that developers don't need to use Adminer.
