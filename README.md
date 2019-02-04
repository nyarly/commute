# Commute

Commute maintains a list of git repositories,
and then allows you to track them across machines easily.

`commute add`
in a git repo
will add the remote to `~/.config/commute/config.yaml`
and a symlink to the checked out workspace
in `~/.config/commute`.
If the remote already appears in `config.yaml`,
the symlink will be created
(if it doesn't exists already)
but no new entry will be added to the file.

`commute list` enumerates those repositories,
with the annotation "-> MISSING"
if commute doesn't know about that repo being checked out.

That's the whole story.
It's easy,
with this limited tool,
to do more useful stuff though.

The author has an `end-of-day` script
that `git status`es everything in the `commute list`,
to make sure local changes are all pushed.

## Gist synchronization

There's a bit of
a chicken-and-egg problem
with the `commute/config.yaml` file.
You're using commute
to keep the files
on various systems
synchronized.
But you need to synchronize
the `commute/config.yaml` file itself.

Since git repos often are hosted on github,
it seemed reasonable to have commute
synchronize itself via a gist.
You'll notice that commute
warns you that your oauth token
hasn't been set up yet
with every invocation.

First,
create personal access token on github,
and plug that into the field in the config.yaml.
Run `commute list` once
(or any subcommand)
and commute will create a gist.
Copy the token and the gist id
through some other channel
to other machines you want to sync to
`et voila!`
