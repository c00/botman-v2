# botman

`botman` lets you talk to an LLM. It is optimized for use in the terminal (or as a library). You can use `stdin` or an argument for input, outputs content to `stdout` and errors to `stderr`.

Note that `botman` does not run any commands. It does not have the ability to _do_ anything, as having some automated LLM control your terminal couldn't possibly lead to anything good. So `botman` simply shows you the output.

You can pipe the output of the last received response to any other shell command, e.g. `botman -h 0 | bash`. Use this at your own risk.

## Supported LLMs

Currently `botman` can use OpenAi, Claude and FireworksAi. Popular choices are:

- GPT 4o
- Claude 3.5 sonnet
- Llama 3 8B / 70B
- Mixtral 7B / 22B

### Complete list of models

- For [OpenAi](https://platform.openai.com/docs/models)
- For [Anthropic](https://docs.anthropic.com/en/docs/about-claude/models)
- For [Fireworks Ai](https://fireworks.ai/models)

## Install from source

1. Clone the repo
2. Run `make install`
3. Run `botman --init` to setup the config. //todo this doesn't work yet.
4. (optional) Create an alias in your shell. e.g. `echo 'alias bot="botman"' >> ~/.bashrc`
5. Test that it works by running `botman "say hi"` or `bot "say hi"`

## Examples

```bash
# Open a new interactive chat.
botman

# Use an argument as input.
botman "tell me a joke about the golang gopher"
botman "git command for undoing last commit"
botman "untar .tar.gz file into new directory"

# Use stdin as input
echo Quote a Bob Kelso joke | botman

# Use both
ls -al | botman "Which files are hidden?"
cat deployment.yaml | botman "how many replicas will this run?"

# Print the last received response
botman -l

# Pipe the last response into bash (at your own risk)
botman -l | bash

# Pipe the last response into a file
botman -h 0 >> some-file.go

# Show the last conversation
botman --history 0

# Show the next-to-last conversation
botman --history 1

# Change model and/or set API keys
botman --init
```

![demo](https://github.com/c00/botman-v2/blob/main/assets/botman-demo.gif?raw=true)

## Interactive mode

In interactive mode, the program does not exit after a response, so you can continue the conversation.

By not supplying any arguments or stdin, `botman` will run in interactive mode.

Use interactive mode explicitly combined with stdin or arguments by giving the `-i` flag.

```bash
# Start interactive mode without an initial prompt
botman

# Start interactive mode with arguments
botman -i "How many bees in a bonnet?"
```

## Data privacy

`botman` talks directly to the API of your configured LLM. So assume that OpenAi / Anthropic / Fireworks knows about your plans to overthrow goverments and such. Other than that, botman does not reach out to any service. It does store your chat history locally in `~/.botman/history`. You can disable this in the settings file `~/.botman/config.yaml` by setting `saveHistory` to `false`.

## Motivation

I created it mainly for myself but thought it might be useful for others. My motivation stems from seeing some closed-source CLI-LLM integrations from companies I don't necessarily trust. So, I created something free and open source for those of us who value open source and transparency.

# Roadmap

I'm adding features as I feel they're useful. I'm open to suggestions and outside contributions. The aim is to be simple, non-intrusive and transparent about data.

- [ ] Add a terminal emulator (tcell, bubbletea, readline, ???)
- [ ] Auto cleanup old conversations
- [ ] Search in old conversation
- [ ] Consider input types like images and documents.
- [ ] Update [slackbot](https://github.com/c00/botman-slack) to use this one instead
- [x] reintroduce the config wizard
  - [ ] add config wizard for tools
- [ ] Separate tests between integration and normals.
- [ ] Tool Use
  - [ ] Implement tool calling in all models
    - [x] Claude (add it to request + deal with result)
    - [ ] OpenAi
    - [ ] Fireworks
  - [ ] Implement Dall-e Tool
