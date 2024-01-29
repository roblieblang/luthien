import { useAuth0 } from "@auth0/auth0-react";
import { Link } from "react-router-dom";

export default function Profile() {
  const { user, isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  return (
    <div className="flex flex-col items-center text-center justify-center my-5">
      <h1>Profile Page</h1>
      {isAuthenticated && (
        <div className="flex flex-col items-center text-center justify-center">
          <img src={user.picture} alt={user.name} />
          <h2>Name: {user.name}</h2> {/* defaults to email if no name given */}
          <p>Email: {user.email}</p>
        </div>
      )}
      <Link to="/" className="rounded-full bg-blue-600 px-3 py-1 hover:bg-blue-200 text-black my-5">
        Back
      </Link>
    </div>
  );
}
