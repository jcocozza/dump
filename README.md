# Dump

Dump is a CLI tool that makes it easy to keep track of, create and share files across machines that you trust.
The idea is peer-to-peer file sync with version control.

It assumes you know about git. You will need to handle merge conflicts yourself.
By default, it runs a rebase when syncing.
The intended use case of dump should rarely produce merge conflicts.
If you are getting lots of conflicts, this tool is probably not for you.

## What this thing does

I have a set of files that I like all of computers to have.
These files have at least 1 of 2 properties:
1. They contain sensitive or personal information (like passwords).
As such, I don't really really want them up in some cloud repo.
I want them on machines I control.

2. They're "scratch" files for projects or ideas I'm working on
It is very, very annoying when I create a well thought-out file on one machine only to later be on another and realize I don't have it. I want a way to have these kinds of files synced across all my machines so I just have them.

3. (2b) A small but useful set of files that I can't find a good place for.
For example, I keep my bookmarks in a text file and the list of books I want to read in another.

(1) Eliminates any kind of cloud based servers. I don't control them and don't want my passwords on them.
While I could have a "remote" git repo running on my home lab set up, (2) is not well suited for it.
Suppose I add a book to my book list.
I don't want to have to think about remembering to push to some central hub repo.
I am far too forgetful for that.

Putting this together, it's pretty clear that I want a peer to peer system.
So I've settled on the following architecture:
Every machine has a local git repo.
That git repo keeps track of all the other machines as "remotes".
"Syncing" is just looping across all remotes and pulling in the various changes from each.

## Some weaknesses:

No machine ever pushes changes to others.
Machines only accept them.
That means if you write some stuff then turn off the computer and go over to another one, it won't be able to access those.
For my particular case this isn't an issue.
If it becomes one then I might figure out how to implement some "push" sync functionality.
