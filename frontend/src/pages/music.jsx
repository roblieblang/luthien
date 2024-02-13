import BackButton from "../components/general/buttons/backButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyUserProfile from "../components/spotify/spotifyUserProfile";

export default function Music() {
  

  return (
    <div>
      <BasicHeading text="Music Page" />
      <SpotifyUserProfile />
      <BackButton linkTo="/" />
    </div>
  );
}
