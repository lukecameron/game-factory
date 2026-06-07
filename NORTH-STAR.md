Note: This file (NORTH-STAR.md) is never to be edited by coding agents. Human-only. Agents can certain read from it.

# What is the game factory?

It's an experiment in developing fun games with very little human iteration. It is inspired by ralph loops.


# How does the game factory work?

That is defined by this repo (game-factory). It's like a ralph loop. Each iteration will run a coding agent with a prompt defined in a file (could call it ITERATION-PROMPT.md). This prompt tells the agent to
* load the current state of the repo (any logs and files left by previous iterations, DEV-LOG.md, TEST-LOG.md, etc)
* pick up a single task to do (could be from something like a ROADMAP.md, that is not defined here)
* complete that task and only that task, resulting in a new wasm artifact for the game
* agenticaly test (could be in a sub-agent, could be using the test harness directly) and come up with gameplay feedback on how well the game plays, feels, etc.
* update DEV-LOG.md, TEST-LOG.md, and ROADMAP.md so that next steps are well represented.

After a certain number of iterations (let's call it $reflectInterval, probably a decent default would be 10), we run a system reflection. The system reflection process also invokes a coding agent, has its own REFLECTION-LOG.md which it reads from and logs its changes to. Its purpose is to go over the dev and test logs, and decide whether the current framework is working well. It is capable of editing everything in the repo except NORTH-STAR.md, which is human-edited ONLY. It looks for ways to optimise the factory architecture. The REFLECTION-LOG should be readable both by humans and by the reflection process itself. Interventions applied by the refelction agent would be things like:
* adjusting $reflectInterval to a new value
* improving the test harness
* editing ITERATION-PROMPT.md to make the iteration process more effective
* editing the supervisor itself. It needs a way to safely reboot the supervisor and resume it if that happens

# What kind of games are we looking for?
* Simple in design, but fun for humans to play
* Beautiful in TUI design (think claude code, btop, k9s, charm's bubble tea TUI framework etc)
* Keyboard-driven
* Familiar - could be a clone of a common game, or a twist on something people know/have played before
* Playable over a ssh link - shellcade is very constrained, it's largely ANSI + a few extensions, low screen size.
* Make use of 24-bit RGB since shellcade/kit supports it
* Low input complexity (ideal is 1-5 separate keys. Don't go crazy)
* Low code complexity. Don't go for a sprawling game design. Probably 5 code files of no more than 1000 lines each is enough for a fun game
* Single-player

# Tools to use
* Use antigravity CLI to run each iteration (`agy --dangerously-skip-permissions -p "$(cat some-prompt.md)"` or something like this
* shellcade/kit wasm CLI game framework. The code and tooling is at ~/dev/kit on this machine, always will be. It is a repo, you can go there and run git pull to get the latest.
* commit automatically to github on each iteration

# Tools built/maintained by us, here in this repo
* A test harness, that coding agents can use to drive the game and test agentic-ly. Likely CLI-based for easy agent-driving. It is able to output a recording that can be played back by a human to see how the testing went.
* A supervisor program. This is in charge of running the iteration agents, and running the reflection agent every $reflectInterval iterations. It also monitors google antigravity rate limits and injects delays between iterations (sleep time) so that the antigravity (google) account never exhausts its rate limit.

# Rules that can never be broken
* The only game artifacts we ever build or test are wasm binaries, as defined in the shellcade/kit repo. 
* Each iteration produces exactly one new wasm artifact, along with test artifacts

