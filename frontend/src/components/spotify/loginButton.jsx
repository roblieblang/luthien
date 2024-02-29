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
      className="bg-customHeadline hover:bg-customButton text-customStroke hover:text-slate-800 font-bold py-1 px-2 my-6 rounded-md border border-black"
      onClick={handleLogin}
    >
      Connect Spotify Account
    </button>
  );
};
