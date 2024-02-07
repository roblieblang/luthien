import { useAuth0 } from "@auth0/auth0-react";

export const LoginButton = () => {
  const { isAuthenticated, user } = useAuth0();

  const handleLogin = () => {
    if (isAuthenticated && user) {
      sessionStorage.setItem("userID", user.sub);
      fetch("http://localhost:8080/auth/spotify/login")
        .then((response) => response.json())
        .then((data) => {
          sessionStorage.setItem("sessionID", data.sessionID);
          window.location.href = data.authURL;
        });
    }
  };

  return (
    <button
      className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
      onClick={handleLogin}
    >
      Connect Spotify Account
    </button>
  );
};
