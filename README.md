# SummVideo

SummVideo is a service that creates concise summaries from video content using AI technology. It extracts key points from any video, saving your time and providing only the most important information.

## 🌟 Features

- Upload video files up to 500MB
- Extract audio from video files automatically
- Transcribe audio to text using OpenAI's Whisper model
- Generate customized summaries using GPT-4o-mini
- Choose from predefined templates or create your own custom prompt
- View complete transcriptions alongside summaries
- Containerized deployment with Docker
- CI/CD pipeline using GitHub Actions

## 🛠️ Tech Stack

- **Backend**: Golang, Gin Web Framework
- **Frontend**: Vue.js 3, HTML, CSS
- **Media Processing**: FFmpeg (via Go wrappers)
- **AI Services**: OpenAI API (Whisper for transcription, GPT-4o-mini for summarization)
- **Containerization**: Docker, Docker Compose
- **Deployment**: DigitalOcean, GitHub Actions

## 📋 Prerequisites

- Docker and Docker Compose
- OpenAI API key
- Git

## 🚀 Installation

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/trofimovm/summvideo.git
   cd summvideo
   ```

2. Create a `.env` file in the project root directory:
   ```bash
   echo "OPENAI_API_KEY=your_openai_api_key" > .env
   ```

3. For frontend development (optional):
   ```bash
   cd frontend
   npm install
   npm run serve
   ```

4. Build and start the Docker container:
   ```bash
   docker compose up -d
   ```

5. Access the application at `http://localhost:8000`

### Production Deployment

The project includes a GitHub Actions workflow that automatically deploys to a DigitalOcean Droplet when changes are pushed to the main branch. To set this up:

1. Create a DigitalOcean Droplet

2. Set up GitHub repository secrets:
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_PASSWORD`: Your Docker Hub password
   - `DEPLOY_KEY`: SSH private key for accessing the DigitalOcean Droplet
   - `OPENAI_API_KEY`: Your OpenAI API key

3. Push changes to the main branch to trigger automatic deployment

## 💻 Usage

1. Open the application in your web browser
2. Select a video file (max 500MB)
3. Choose a template prompt or create your own custom prompt
4. Click "Submit" to start processing the video
5. View the generated summary and, optionally, the full transcription

### Available Templates

- **Meeting Summary**: Ideal for meeting recordings, focuses on context, key points, decisions, and next steps
- **Sales Meeting Summary**: Specialized for sales meetings, includes information about products/services, client questions, and agreements
- **Educational Video**: Optimized for educational content, highlights learning objectives, main topics, and key takeaways

## 📁 Project Structure

```
summvideo/
├── .github/workflows/   # GitHub Actions workflow configurations
├── backend-go/          # Golang backend application
│   ├── handlers/        # HTTP request handlers
│   ├── services/        # Business logic layer
│   ├── models/          # Data structures
│   ├── utils/           # Helper functions
│   ├── Dockerfile       # Go service containerization
│   ├── go.mod           # Go dependencies
│   └── .env.example     # Example environment variables
├── frontend/            # Vue.js frontend application
│   ├── src/             # Vue components and application logic
│   ├── public/          # Static assets for Vue.js
│   └── package.json     # Frontend dependencies
├── static/              # Static assets (CSS, JS)
│   └── vue/             # Built Vue.js application (after build)
├── templates/           # HTML templates
├── .gitignore           # Git ignore file
├── docker-compose.yml   # Docker Compose configuration
├── README.md            # Project documentation
└── CLAUDE.md            # Helper documentation for LLM assistants
```

## ⚙️ How It Works

1. **Video Upload**: Users upload video files through the Vue.js interface
2. **Concurrent Processing**: Go backend handles the request in a separate goroutine
3. **Audio Extraction**: FFmpeg extracts audio from the video
4. **Transcription**: OpenAI's Whisper model transcribes the audio to text
5. **Summarization**: OpenAI's GPT-4o-mini generates a summary based on the chosen prompt
6. **Result Display**: The summary is rendered in markdown format in the Vue.js interface, with an option to view the full transcription

## 📝 Logging

The application logs all processed videos, including:
- Timestamp
- Video filename
- Prompt used
- Transcription
- Generated summary

Logs are stored in `/var/log/summvideo/log.txt` inside the container.

## 🔒 Security Notes

- API keys are stored as environment variables
- Temporary files are created for processing and deleted afterward
- File size limits are enforced to prevent abuse

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Open a pull request

## 📄 License

This project is licensed under the MIT License.