const getToken = async (code) => {
  let codeVerifier = localStorage.getItem("code_verifier");

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

  const body = await fetch(authUrl, payload);
  const response = await body.json();

  localStorage.setItem("access_token", response.access_token);
};

// getToken(code);

// TODO: add request access token functionality as shown towards bottom of this doc: https://developer.spotify.com/documentation/web-api/tutorials/code-pkce-flow

// TODO: integrate into react components

// TODO: Take a look at to the access token guide (https://developer.spotify.com/documentation/web-api/concepts/access-token)
// to learn how to make an API call using your new fresh access token.

// TODO: Learn about refresh tokens: https://developer.spotify.com/documentation/web-api/tutorials/refreshing-tokens
