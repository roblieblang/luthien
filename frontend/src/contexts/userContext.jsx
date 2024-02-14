import { useAuth0 } from "@auth0/auth0-react";
import { createContext, useContext, useEffect, useState } from "react";

const UserContext = createContext();

export const useUser = () => useContext(UserContext);

export const UserProvider = ({ children }) => {
  const { isAuthenticated, user } = useAuth0();
  const [userID, setUserID] = useState(null);
  const [spotifyUserID, setSpotifyUserID] = useState(null);

  const updateSpotifyUserID = (spotifyUserID) => {
    setSpotifyUserID(spotifyUserID);
  };

  useEffect(() => {
    if (isAuthenticated && user) {
      setUserID(user.sub);
    } else {
      setUserID(null);
    }
  }, [isAuthenticated, user]);

  return (
    <UserContext.Provider
      value={{ userID, spotifyUserID, updateSpotifyUserID,  }}
    >
      {children}
    </UserContext.Provider>
  );
};
