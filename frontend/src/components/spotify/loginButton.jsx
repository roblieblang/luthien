import { useAuth0 } from "@auth0/auth0-react";
// TODO: merge with Spotify logoutButton(?)
export const LoginButton = () => {
  const { isAuthenticated, user } = useAuth0();

  const handleLogin = () => {
    if (isAuthenticated && user) {
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
      className="bg-customButton hover:bg-green-700 text-customStroke font-bold py-2 px-4 rounded-full border-2 border-black"
      onClick={handleLogin}
    >
      Connect Spotify Account
    </button>
  );
};
