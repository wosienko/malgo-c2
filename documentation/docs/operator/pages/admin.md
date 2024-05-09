# Admin Panel

Admin panel is divided into two main sections: users and projects. The admin panel is only accessible to users with the admin role.

Relevant components are stored in `src/lib/components/custom/admin`.

Relevant API routes are stored in `src/routes/api/(admin)/`.

If you were to open `+layout.svelte`, you may notice, that anchor tags have `data-sveltekit-preload-data` attribute. This is used to load data on mouse hover. This way, the data is loaded before the user clicks on the link.

```html hl_lines="7"
{#each ALL_TABS as tab}
      <a
      role="tab"
      class="tab pb-9 hover:opacity-50"
      class:tab-active={currentUrl === tab}
      href={tab}
      data-sveltekit-preload-data>
            {capitalize(tab.replace('/admin/', ''))}
      </a>
{/each}
```

## Users

The users section allows for the management of users. The following actions are available:

- Create new users
- Edit existing users
      - Change name
      - Change surname
      - Change email
      - Change role
- Change user password
- Delete users

## Projects

The projects section allows for the management of projects. The following actions are available:

- Create new projects
- Edit existing projects
      - Change project name
      - Change project description
- Assign users to projects
- Delete projects
