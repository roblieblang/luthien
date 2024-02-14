import { useAuth0 } from "@auth0/auth0-react";
import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";

export default function Profile() {
  const { user, isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  return (
    <div className="flex flex-col items-center text-center justify-center">
      <BasicHeading text="Profile Page" />
      {isAuthenticated && (
        <div className="flex flex-col items-center text-center justify-center">
          <img src={user.picture} alt={user.name} />
          <h2>Name: {user.name}</h2> {/* defaults to email if no name given */}
          <p>Email: {user.email}</p>
        </div>
      )}
      <LinkButton to="/" text="Back" />
    </div>
  );
}
