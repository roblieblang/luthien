# Luthien

## Table of Contents

- [Luthien](#luthien)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [🌟 Features](#-features)
  - [📸 Screenshots](#-screenshots)
  - [⚙️ Technologies Used](#️-technologies-used)
    - [Backend](#backend)
    - [Frontend](#frontend)
  - [💻 Getting Started](#-getting-started)
    - [Requirements](#requirements)
    - [Installation](#installation)
    - [Backend](#backend-1)
    - [Frontend](#frontend-1)
  - [Usage](#usage)
  - [Contribution](#contribution)
  - [License](#license)

## Overview

Luthien is a playlist conversion application designed to facilitate cross-platform listening between Spotify and YouTube. Users can easily transfer their playlists from one platform to another, ensuring seamless music enjoyment across different services.
<ul>
  <li>
    <a href="https://youtu.be/FOItY3HnoPs" target="_blank">
      Demo Video
    </a>
  </li>
</ul>

## 🌟 Features

- Convert playlists from Spotify to YouTube, or vice versa
- Integration with Spotify API and YouTube Data API
- Redis Caching strategies for optimizing 3rd party API calls
- Responsive design for various devices
- Authentication integrated with Auth0

## 📸 Screenshots

![Homepage](https://raw.githubusercontent.com/roblieblang/luthien/main/images/homepage.png)
![Spotify Account](https://raw.githubusercontent.com/roblieblang/luthien/main/images/spotify.png)
![Spotify Conversion](https://raw.githubusercontent.com/roblieblang/luthien/main/images/spotify-convert.png)
![Conversion Success](https://raw.githubusercontent.com/roblieblang/luthien/main/images/youtube-track-hits.png)

## ⚙️ Technologies Used

### Backend

- Go
- Gin
- Redis
- Docker
- Auth0
- Google Cloud Platform (Cloud Run)
- Upstash

### Frontend

- JavaScript
- React
- Vite
- Tailwind
- Auth0
- Vercel

## 💻 Getting Started

### Requirements

- [Docker](https://docs.docker.com/engine/install/)
- [Node.js](https://nodejs.org/en/download/package-manager)

### Installation

1. Clone the repository: `git clone https://github.com/roblieblang/luthien.git`
2. Navigate to the project directory: `cd luthien`

### Backend

1. Navigate to the backend directory: `cd backend`
2. Start Redis and MongoDB with Docker: `docker-compose up --build`
3. Start the server with live reloading: `air`

### Frontend

1. Navigate to the frontend directory: `cd frontend`
2. Install dependencies: `npm install`
3. Start the development server: `npm run dev`

## Usage

1. Open your browser and go to `http://localhost:5173/`
2. Sign in with your Spotify and YouTube accounts
3. Follow the on-screen instructions to convert your playlists

## Contribution

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a new branch: `git checkout -b feature-name`
3. Make your changes and commit them: `git commit -m 'Add new feature'`
4. Push to the branch: `git push origin feature-name`
5. Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
