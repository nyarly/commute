# Commute

Commute maintains a list of git repositories, and then allows you to track them across machines easily.

`commute add` in a git repo will add the remote to `~/.config/commute/config.yaml`
and a symlink to the checked out workspace in `~/.config/commute`.
If the remote already appears in `config.yaml`, the symlink will be created
(if it doesn't exists already)
but no new entry will be added to the file.

`commute list` enumerates those repositories, with the annotation "-> MISSING"
if commute doesn't know about that repo being checked out.

That's the whole story.
It's easy,
with this limited tool,
to do more useful stuff though.

The author has an `end-of-day` script
that `git status`es everything in the `commute list`,
to make sure local changes are all pushed.
