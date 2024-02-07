export const LogoutButton = ({ setIsAuthenticatedWithSpotify, userID }) => {
  const handleLogout = () => {
    fetch("http://localhost:8080/auth/spotify/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ userID }),
    }).then((response) => {
      if (response.ok) {
        setIsAuthenticatedWithSpotify(false); // Callback to update the parent component's state
      }
    });
  };

  return (
    <>
      <button
        className="bg-green-600 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
        onClick={handleLogout}
      >
        Log Out of Spotify
      </button>
    </>
  );
};
