# GoSlideShow

GoSlideShow is a proof-of-concept (POC) command-line tool written in Go. It automates the creation and modification of Google Slides presentations by leveraging Google APIs and OpenAI's language model. The goal is to demonstrate how to interact with an AI based on large language models (LLMs) to generate structured content.

## Features

- **Generate Slides from Markdown**: Convert Markdown files into structured Google Slides presentations.
- **Template Support**: Create presentations based on a specified Google Slides template.
- **OAuth2 Authentication**: Uses Google's OAuth2 for secure API access.

## Demo

https://github.com/owulveryck/gptslideshow/raw/refs/heads/main/demo.mov

## Prerequisites

- Go 1.18 or later
- Google API credentials (`credentials.json`)
- Access to the Google Slides and Drive APIs
- **OpenAI API Key**: Set the `OPENAI_API_KEY` environment variable with your OpenAI API key.

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/goslideshow.git
   cd goslideshow
   ```

2. **Install dependencies**:

   Ensure you have the required Go packages installed. You can use Go modules to handle dependencies:

   ```bash
   go mod tidy
   ```

3. **Setup Google API Credentials**:

   Place your `credentials.json` file in the root directory of the project. This file should contain your Google API credentials.

4. **Set OpenAI API Key**:

   Export your OpenAI API key as an environment variable:

   ```bash
   export OPENAI_API_KEY=your_openai_api_key
   ```

## Usage

Run the program with the following command:

```bash
go run main.go -content <path-to-markdown-file> [-t <template-id>] [-id <presentation-id>]
```

- `-content`: Path to the Markdown file to convert into slides.
- `-t`: (Optional) ID of the Google Slides template to use.
- `-id`: (Optional) ID of an existing presentation to update.

## File Structure

- **main.go**: The entry point of the application. It handles command-line arguments, initializes services, and orchestrates the creation of slides.

- **presentation.go**: Contains data structures and functions for generating presentations from Markdown content. It uses OpenAI's API to convert content into slide format.

- **google_auth.go**: Manages OAuth2 authentication, including token storage and retrieval. It ensures secure access to Google APIs.

- **slides_operations.go**: Contains functions for interacting with Google Slides, such as copying templates and creating new slides. It handles the insertion of content into slide placeholders.

## Authentication

The tool uses OAuth2 to authenticate with Google APIs. The first time you run the tool, it will prompt you to visit a URL and authorize the application. The token will be cached for future use.

## Disclaimer

This project is a proof-of-concept intended for educational purposes. It demonstrates how to interact with AI models and Google APIs to automate tasks. It is not intended for production use. I will not provide any support for this code.

## Contribution

As a POC, I will not make this project evolve. Feel free to use, fork, or modify it as you see fit.

## License

This project is licensed under the MIT License.
