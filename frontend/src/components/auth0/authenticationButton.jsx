import { useAuth0 } from "@auth0/auth0-react";

const AuthenticationButton = () => {
  const { loginWithRedirect, logout, isAuthenticated } = useAuth0();

  return isAuthenticated ? (
    <button
      className="bg-gray-500 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
      onClick={() =>
        logout({ logoutParams: { returnTo: window.location.origin } })
      }
    >
      Log Out (Auth0)
    </button>
  ) : (
    <button
      className="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
      onClick={() =>
        loginWithRedirect({
          access_type: "offline",
          connection_scope: "https://www.googleapis.com/auth/youtube",
          offline: true,
          prompt: "consent",
        })
      }
    >
      Log In (Auth0)
    </button>
  );
};

export default AuthenticationButton;
