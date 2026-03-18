import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/useAuth";
export default function PublicRoute() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;

  // already logged in — send them to right place
  if (user && !user.role) return <Navigate to="/select-role" replace />;
  if (user && user.role) return <Navigate to="/" replace />;

  // not logged in → show the public page
  return <Outlet />;
}
