import { useAuth0 } from "@auth0/auth0-react";

const AuthenticationButton = () => {
  const { loginWithRedirect, logout, isAuthenticated } = useAuth0();

  return isAuthenticated ? (
    <button
      className="bg-customHeadline hover:bg-customButton text-md  text-customStroke hover:text-slate-800 font-bold py-1 px-2 rounded-md border border-black"
      onClick={() =>
        logout({ logoutParams: { returnTo: window.location.origin } })
      }
    >
      Log Out
    </button>
  ) : (
    <button
      className="bg-customHeadline hover:bg-customButton text-md text-customStroke hover:text-slate-800 font-bold py-1 px-2 rounded-md border border-black"
      onClick={() =>
        loginWithRedirect({
          access_type: "offline",
          connection_scope: "https://www.googleapis.com/auth/youtube",
          offline: true,
          prompt: "consent",
        })
      }
    >
      Get Started
    </button>
  );
};

export default AuthenticationButton;
