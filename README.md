# Atomic Mustache

A tiny, very opinionated package to render mustache templates based on the atomic design pattern. It makes it possible to use your mustache templates live in production, without duplicating your templates all over your apps.

It expects your template folder structure to be the following:

    /layouts      - Rails-style layouts for the application
    /styleguide   - The atomic design folder structure (atoms, molecules, organisms, etc)
    /views        - Rails-style views for the application.

In a normal scenario, you will have a main layout template that holds the HTML skeleton of the application. Then, you have views for each controller action (like `posts/index.mustache`), and those views uses the atomic design partials to render pages.

## Tests

`go test`
