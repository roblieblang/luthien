import { logoutClick } from "../../utils/auth-utils";
export const SpotifyLogoutButton = () => {
  return (
    <>
      <button
        className="bg-green-600 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full"
        onClick={logoutClick}
      >
        Logout of Spotify
      </button>
    </>
  );
};
