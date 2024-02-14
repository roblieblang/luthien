import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyUserProfile from "../components/spotify/spotifyUserProfile";

export default function Music() {
  return (
    <div className="mb-2">
      <BasicHeading text="Music Page" />
      <SpotifyUserProfile />
      <LinkButton to="/" text="Back" />
    </div>
  );
}
