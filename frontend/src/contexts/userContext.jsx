import { useAuth0 } from "@auth0/auth0-react";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { config } from "../utils/config";

const UserContext = createContext();

export const useUser = () => useContext(UserContext);

export const UserProvider = ({ children }) => {
  const { isAuthenticated, user, logout } = useAuth0();
  const [userID, setUserID] = useState(null);
  const [spotifyUserID, setSpotifyUserID] = useState(null);
  const [googleUserID, setGoogleUserID] = useState(null);
  const [spotifyAuthStatus, setSpotifyAuthStatus] = useState(false);
  const [youTubeAuthStatus, setYouTubeAuthStatus] = useState(false);

  useEffect(() => {
    if (isAuthenticated && user) {
      sessionStorage.setItem("userID", user.sub);
      setUserID(user.sub);
    } else {
      setUserID(null);
    }
  }, [isAuthenticated, user]);

  const checkAuthStatus = useCallback(
    (serviceURL, setIsAuthenticated) => {
      if (userID) {
        fetch(`${serviceURL}?userID=${userID}`)
          .then((res) => {
            if (!res.ok) {
              throw new Error(`Failed to fetch, status code: ${res.status}`);
            }
            return res.json();
          })
          .then((data) => {
            setIsAuthenticated(data.isAuthenticated);
          })
          .catch((error) => {
            console.error("Fetching error:", error);
            logout({ returnTo: window.location.origin });
          });
      }
    },
    [userID, logout]
  );

  // Spotify
  useEffect(() => {
    checkAuthStatus(
      `${config.backendUrl}/auth/spotify/check-auth`,
      setSpotifyAuthStatus
    );
  }, [userID, checkAuthStatus]);

  // Google/YouTube
  useEffect(() => {
    checkAuthStatus(
      `${config.backendUrl}/auth/google/check-auth`,
      setSpotifyAuthStatus
    );
  }, [userID, checkAuthStatus]);

  return (
    <UserContext.Provider
      value={{
        userID,
        spotifyUserID,
        updateSpotifyUserID: setSpotifyUserID,
        spotifyAuthStatus,
        updateSpotifyAuthStatus: setSpotifyAuthStatus,
        googleUserID,
        updateGoogleUserID: setGoogleUserID,
        youTubeAuthStatus,
        updateYouTubeAuthStatus: setYouTubeAuthStatus,
      }}
    >
      {children}
    </UserContext.Provider>
  );
};
