import { useEffect, useState } from "react";
import { useUser } from "../../contexts/userContext";
import SpotifyUserPlaylists from "./spotifyUserPlaylists";

export default function SpotifyUserProfile() {
  const { userID, updateSpotifyUserID } = useUser();
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    if (userID) {
      fetch(`http://localhost:8080/spotify/current-profile?userID=${userID}`)
        .then((res) => {
          if (!res.ok) {
            throw new Error("Response from server was not ok");
          }
          return res.json();
        })
        .then((data) => {
          setProfile(data);
          updateSpotifyUserID(data.id);
        })
        .catch((error) => {
          console.error("Error fetching user profile:", error);
        });
    }
  }, [userID, updateSpotifyUserID]);

  if (!profile) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1 className="font-bold text-lg">
        {profile.display_name}&apos;s Spotify Account
      </h1>
      <SpotifyUserPlaylists />
    </div>
  );
}
