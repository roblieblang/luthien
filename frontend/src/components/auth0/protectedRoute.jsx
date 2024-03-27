import { useAuth0 } from "@auth0/auth0-react";
import { Bars } from "react-loader-spinner";
import { Navigate } from "react-router-dom";

export default function ProtectedRoute({ children }) {
  const { isAuthenticated, isLoading } = useAuth0();

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

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return children;
}
