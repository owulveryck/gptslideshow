# GPT-Powered Google Slides Automation Tool

This tool automates text updates in Google Slides presentations using AI services like OpenAI GPT and Ollama. It allows you to format content, generate AI-powered text responses, and seamlessly integrate them into your slides.

---

## Features

- **AI-Powered Text Processing**: Use GPT models (OpenAI and Ollama) to generate or format slide content dynamically.
- **Markdown Support**: Easily insert formatted text (via Markdown) into Google Slides.
- **Command Tags**: Annotate slide content with tags (`@format`, `@chatgpt`, `@ollama`) to trigger specific AI-powered actions.
- **Context-Aware AI Responses**: Include content from previous slides to provide context for GPT responses.

---

## Prerequisites

1. **Go**: Ensure you have Go installed ([Download Go](https://golang.org/dl/)).
2. **Google Cloud Project**:
   - Enable the Google Slides API.
   - Set up OAuth credentials to access Google services.
3. **AI Clients**:
   - OpenAI API Key for GPT.
   - Ollama setup for local LLMs.

---

## Installation

Clone the repository and build the project:

```bash
git clone <repository-url>
cd <project-folder>
go build -o gptslideshow
```

---

## Usage

Run the program with the `-id` flag to specify a Google Slides presentation ID:

```bash
./gptslideshow -id <presentation-id>
```

- If the `-id` flag is left empty, a new presentation will be created.

### Supported Commands in Slides

1. **`@format`**:
   - Triggers formatting of the text using Markdown.
2. **`@chatgpt`**:
   - Sends the slide content to OpenAI GPT for AI-generated responses.
   - Use `@withContext` to include all previous slide content for contextual AI responses.
3. **`@ollama`**:
   - Sends the slide content to the Ollama AI for local LLM processing.

**Example**:

To process a slide with GPT and include context:

```text
@chatgpt @withContext
Summarize the following content and format it nicely.
```

---

## Workflow

1. The tool retrieves all slides in the specified presentation.
2. For each slide:
   - Text content is extracted from shapes.
   - AI tags (`@format`, `@chatgpt`, `@ollama`) are detected.
   - Corresponding actions (e.g., AI queries or formatting) are performed.
3. The updated content is applied back to the slides via Google Slides API.

---

## Code Structure

- **`main.go`**: The core logic for reading slides, processing AI commands, and updating content.
- **`slidesutils`**: Utility functions for processing and inserting Markdown-formatted text.
- **`ai/ollama`** & **`ai/openai`**: Clients for interacting with AI services.

---

## Error Handling

- The tool logs any errors encountered (e.g., API failures, invalid slide IDs) and stops execution to ensure reliability.

---

## Contributing

Feel free to fork this repository and submit pull requests for improvements.

---

## License

This project is licensed under the MIT License.
