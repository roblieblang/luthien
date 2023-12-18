import {
  base64encode,
  generateRandomString,
  sha256,
} from "../../utils/auth-utils";

export const SpotifyLoginButton = () => {
  // TODO: Refactor and move this function to ../../utils/auth-utils.js
  const handleLogin = async () => {
    const codeVerifier = generateRandomString(64);
    const hashedVerifier = await sha256(codeVerifier);
    const codeChallenge = base64encode(hashedVerifier);

    // console.log(`Code challenge: ${codeChallenge}`);

    const spotifyClientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID;
    const redirectUri = import.meta.env.VITE_REDIRECT_URI;

    const scope = "user-read-private user-read-email";
    const authUrl = new URL("https://accounts.spotify.com/authorize");

    window.localStorage.setItem("code_verifier", codeVerifier);

    const params = {
      response_type: "code",
      client_id: spotifyClientId,
      scope,
      code_challenge_method: "S256",
      code_challenge: codeChallenge,
      redirect_uri: redirectUri,
    };

    authUrl.search = new URLSearchParams(params).toString();
    window.location.href = authUrl.toString();
  };

  return (
    <button
      className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 my-5 rounded-full"
      onClick={handleLogin}
    >
      Login to Spotify
    </button>
  );
};
