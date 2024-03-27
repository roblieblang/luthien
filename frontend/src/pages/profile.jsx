import { useAuth0 } from "@auth0/auth0-react";
import { Bars } from "react-loader-spinner";
import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";

export default function Profile() {
  const { user, isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Bars
          height="80"
          width="80"
          color="#e2714a"
          ariaLabel="bars-loading"
          visible={true}
        />
      </div>
    );
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
