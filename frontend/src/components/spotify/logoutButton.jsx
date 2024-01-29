export const LogoutButton = ({setIsAuthenticatedWithSpotify}) => {
  const handleLogout = () => {
    fetch("http://localhost:8080/auth/spotify/logout", { method: "POST" }).then(
      (response) => {
        if (response.ok) {
          setIsAuthenticatedWithSpotify(false); // Callback to update the parent component's state
        }
      }
    );
  };

  return (
    <>
      <button
        className="bg-green-600 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
        onClick={handleLogout}
      >
        Logout of Spotify
      </button>
    </>
  );
};
