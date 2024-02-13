import { useEffect, useState } from "react";
import SpotifyUserPlaylists from "./spotifyUserPlaylists";

export default function SpotifyUserProfile() {
  const [profile, setProfile] = useState(null);

  const userID = sessionStorage.getItem("userID");

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
        })
        .catch((error) => {
          console.error("Error fetching user profile:", error);
        });
    }
  }, []);

  if (!profile) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>{profile.display_name}&apos;s Spotify Account</h1>
      <p>{profile.followers.total} followers</p>
      <SpotifyUserPlaylists userID={userID} />
    </div>
  );
}
