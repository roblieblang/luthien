// TODO: Refactor some of auth-utils into this file

// Requests a Spotify API access token using an authorization code
export const getSpotifyAccessToken = async (code) => {
  let codeVerifier = localStorage.getItem("code_verifier");

  if (!codeVerifier) {
    console.log("Empty code verifier!");
  }

  const spotifyClientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID;
  const spotifyTokenUrl = new URL("https://accounts.spotify.com/api/token");
  const redirectUri = import.meta.env.VITE_REDIRECT_URI;

  const payload = {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: spotifyClientId,
      grant_type: "authorization_code",
      code,
      redirect_uri: redirectUri,
      code_verifier: codeVerifier,
    }),
  };

  const body = await fetch(spotifyTokenUrl, payload);
  const response = await body.json();
  console.log(`Response: ${JSON.stringify(response)}`);
  console.log(`getSpotifyAccessToken() access_token: ${response.access_token}`);

  // TODO: store tokens server-side(?)
  localStorage.setItem("access_token", response.access_token);
};

// Get a Spotify API refresh token
export const getSpotifyRefreshToken = async () => {
  // refresh token that has been previously stored
  const refreshToken = localStorage.getItem("refresh_token");
  const url = "https://accounts.spotify.com/api/token";

  const clientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID;

  const payload = {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      grant_type: "refresh_token",
      refresh_token: refreshToken,
      client_id: clientId,
    }),
  };
  const body = await fetch(url, payload);
  const response = await body.json();

  localStorage.setItem("access_token", response.accessToken);
  localStorage.setItem("refresh_token", response.refreshToken);
};

// Get the current user's Spotify profile
export const getSpotifyProfile = async () => {
  let accessToken = localStorage.getItem("access_token");

  const response = await fetch("https://api.spotify.com/v1/me", {
    headers: {
      Authorization: "Bearer " + accessToken,
    },
  });

  const data = await response.json();
};
