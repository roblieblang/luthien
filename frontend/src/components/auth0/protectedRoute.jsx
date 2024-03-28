import { useAuth0 } from "@auth0/auth0-react";
import { Navigate } from "react-router-dom";
import Loading from "../general/modals/loading";

export default function ProtectedRoute({ children }) {
  const { isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <Loading />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return children;
}
