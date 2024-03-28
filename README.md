# Luthien

## Overview
I built a playlist conversion application for cross-platform listening between Spotify and YouTube.
<ul>
  <li>
    <a href="" target="_blank">
      Try here (awaiting verification by API providers)
    </a>
  </li>
  <li>
    <a href="" target="_blank">
      Video (coming soon)
    </a>
  </li>
</ul>

## Features
- Complete integration with Spotify API and YouTube Data API
- Robust Redis Caching strategies for expensive 3rd party API calls (such as search)
- Responsive design
- Full authentication integrated with Auth0

## Technologies Used
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

## Local Setup
### Requirements
- Docker
### Backend
1. Run `docker-compose up --build` from the root of `./backend` to start Redis and MongoDB
2. Start the server from `./backend` with command `air` for live reloading

### Frontend
1. Start the dev server from `./frontend` directory with command `npm run dev`


