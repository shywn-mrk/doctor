# Doctor

A simple document extractor which look for a block of a code which starts with a certain patten and merge this groups as
a section of the documentation. This way your business/technical/... documents can be as close as possible to the code
for the developer and share the same life cycle. If the code changes, changing the document is easy and vice versa.

## Build and Run

- Install `ag` (You can find it searching *The Silver Searcher*)


## TODO

- [ ] What is commenting in the target text?
- [ ] What is the pattern that flags a part of code as a documentation?
- [ ] We need to include h1, h2, ..., h6 into consideration.
- [ ] Can we use `ag -i --php -0 --nofilename --nobreak --nocolor '/\*+?\s*?\@doctor(.*?\n)*?.*?\*/' ~/Documents/Code/Digikala/supernova-env-dev/vendor/digikala`
    - Then we feed the output to a program to convert it to markdown.
    - Pass that markdown to another program to create a web page from it.
    - [ ] Which one is faster? more customizable? more presentable?

## Ideas

* Notify people when a section changes. They might need to know if the logic is changed.