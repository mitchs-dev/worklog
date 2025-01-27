# Worklog

Worklog is a CLI tool that helps you track your personal work log. The idea is that you can track your tasks and the time you spend on them. This allows you to see the amount of work that you have accomplished over the course of a day, week, month, quarter, or year.

> **Note**: Worklog is intended to be used for personal work tracking. The goal is to provide the user with a simple means to track their work and provide insights into their productivity. It is not intended to be used as a time tracking tool for performance evaluations. Any reports of misuse will land you on the wall of shame.

## Principles

A work log is intended to provide a simple means to provide you with insights on your productivity. It is meant to help you evaluate your work habits, have a reference point of tasks, recognize patterns, and most importantly, show how much you have accomplished. 

A work log should not be used as a tool for comparison or "bragging rights." It should be used as a tool for reference , self-improvement and self-reflection and should be treated as such.

### What do you add to your work log?

You can add whatever you like! Add as much or as little as you want. The general rule of thumb is add what is useful and meaningful to you. 

For example, if you are the type that easily "forgets what you had for breakfast," you might add a lot of details to your work log. But if you remember everything you do, you might add less. There is no wrong way to use a work log as long as you are using it in a way that is meaningful to you.

### Always forward, never back

When adding to a worklog, you should only be logging what you have done for the day. What was done yesterday is yesterday's news. What is done tomorrow is tomorrows problem. You should only be logging what you have done today.

### Log at your own pace

People are different. Some people may find logging as you go easier. But some people may want to log at the end of the day. The important thing is that you are logging in a way that you are not forgetting details of what you have done for the day.

### Start the day off right

In the context of review, unlike adding logs, it's always a good idea to start the day off by reviewing what you did yesterday. This will help you remember where you left off and help you get back into the groove of things.

> **Pro Tip**: At the end of the week, it's not a bad idea to review the entire week to see what you have accomplished. This can help you see patterns and help you plan for the next week.

### Log what you do, not what you plan to do

A work log is a record of what you have done, not what you plan to do. There are plenty of tools out there to help you manage your upcoming work, but a work log is meant to be a record of what you have done.

### Be honest

A work log is like a diary or a journal, only you will see it, so be honest with yourself.

### Log it and forget it

Don't _really_ forget it! But the idea is that you shouldn't be dwelling on what to log or wanting to go back and edit it. You should just log it and move on. If you happen to make a grammatical error or something, unless it affects the meaning of the entry, it's not worth going back and editing it. Stuff happens!



## Getting Started

### Installation

To install worklog, you simply need to download the binary for your platform from the [releases page](https://github.com/mitchs-dev/worklog/releases) and run:

```bash
chmod +x worklog
```
and it is ready to use. You can move it to a directory in your PATH to make it easier to access. For example:

```bash
mv worklog /usr/local/bin
```

### Usage

#### Add an entry

To add an entry you simply need to run:

```bash
worklog add <entry>
```

This will add an entry to your work log for the current day.

For example:
```bash
$ worklog add View the worklog readme
[2025-01-23T15:14:23|INFO|add.go:44(command.go:989)]: Entry ID: 0123-12
```

> **Note**: There is no plan to be able to add entries for previous days. This is intentional. (See [Principles - Always forward, never back](#always-forward-never-back))

#### List entries

To list entries you simply need to run:

```bash
worklog list
```

This will list all the entries for the current day. There are many other options for the list command so it is recommended to play around with it!

Here is an example of what the output might look like:

```bash
list
Period: today
Worklog:
- [0123-7] Optimized database queries. Reduced response time from "grab a coffee" to "blink and you'll miss it".
- [0123-8] Conducted code review for PR #1337. Suggested renaming variables from "x" to something more descriptive, like "y".
- [0123-9] Researched microservices architecture. Concluded that "micro" is a relative term.
```


### Enable sync with Git

Worklog has the ability to sync your work log with a Git repository. This is useful if you want to keep a backup of your work log, or use it across multiple devices.

While `worklog` _does_ have the ability to assist with the creation of a git repository, it is highly recommended to follow the steps below instead. This is because the interactive setup for a git repository is not very robust and may not work as expected.

#### Pre-requisites

This assumes that:

1. You have installed `git` on your machine.
2. You have configure authentication with your Git hosting service (e.g. GitHub, GitLab, Bitbucket).
3. Your remote is configured and the name is set to `origin`. 

#### Steps

1. Create an empty Git repository in your favorite Git hosting service (e.g. GitHub, GitLab, Bitbucket).
2. Clone the repository to your local machine.
3. Set the `.settings.logs.path` to the path of your local repository.
4. Set the `.settings.git.sync` to `true`.
5. Set the `.settings.git.uri` to the URI of your remote repository. (HTTPS or SSH is supported)
6. Set the `.settings.git.branch` to the branch you want to push to.

Once you have completed the above, you can run:

```bash
worklog sync
```

This will bi-directionally sync your work log with the remote repository.

> **Pro Tip**: It is recommended to run `worklog sync` at the end of the day to ensure that your work log is backed up. And also make sure that you run it before swapping devices (if you are using multiple devices).

It's not recommended to manually mess with your worklogs repository. If you need to make changes, it is recommended to do so through the `worklog` CLI. If required, you can also use the `--force` flag to overwrite the remote repository with your local repository.


## Planned Features

As of now, the existing implementation was designed with the [Principles](#principles) outlined above in mind; which is what I wanted for a `v1` release. 

However, there are some other features that I would like to add in the future that I feel would complement the existing implementation. These are:

- [ ] Add a `remove` command for accidental entries. (Keeping [Log it and forget it](#log-it-and-forget-it) in mind)
- [ ] Add time tracking capabilities. Such as `start`,`pause`,`resume`,`stop`.
  * This wouldn't affect those that don't want to use this and it also wouldn't affect backwards compatibility.
- [ ] Add (optional) workday restrictions to make sure you are not logging entries in the off hours. 🙂
- [ ] Configuration editing from the CLI.