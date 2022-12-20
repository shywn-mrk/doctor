# Doctor

A simple document extractor which look for a block of a code which starts with a certain patten and merge this groups as
a section of the documentation. This way your business/technical/... documents can be as close as possible to the code
for the developer and share the same life cycle. If the code changes, changing the document is easy and vice versa.

## Build and Run

- Install `ripgrep`
- Install `go` compiler

## How it works?

Using `rg` command to find block of comments matching what we define as *docblock*. Pass that to anothe layer of `rg` to
filter out unnecessary data. What remains are many lines separated by a `@doctor` at the beginning of the some of lines.
Now we can group these group based on simple groupings we have using our Go code.

## Notes

- [ ] How are we going to sort parapraghs of each section to keep the order the same when generating the document multiple times in a row?
- [ ] mermaid.js, pandoc, ... for markdown to html.
- [ ] some thing like tags maybe?

## Ideas

* Notify people when a section changes. They might need to know if the logic is changed.
