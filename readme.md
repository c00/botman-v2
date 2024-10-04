# Refactor Todo list

As learn more about Go, I see I did things poorly. To refactor, I should:

- [x] Fix the file structure with an internal folder and a cmd folder for better consumption
- [x] Create a proper interface for a chat bot
- [x] Test said interface, also create a fake chatbot, so it becomes very easy to test new ones.
- [x] Implement openAi
- [x] Implement fireworksLlms
- [x] Test that history is saved
- [x] Consider a proper implementation for tool use. (And consider shit like attachments like generated images.)
- [ ] Consider input types like images and documents.

- [ ] Replace allll of the old repo with this one.
- [ ] Update slackbot to use this one instead

## Tool Use

- [x] ChatMessage should have a []toolCall property
- [x] Mainloop executes the tool calls
- [x] Configuration should just have a bunch of tools enabled by default
- [ ] Implement tool calling in
  - [x] Claude (add it to request + deal with result)
  - [ ] OpenAi
  - [ ] Fireworks
- [x] Implement SDXL Tool
- [x] Add s3 as storage option
- [ ] Implement Dall-e Tool
