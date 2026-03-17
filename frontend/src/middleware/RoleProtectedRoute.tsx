import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/useAuth";

export default function RoleProtectedRoute() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (!user) return <Navigate to="/login" replace />;
  if (!user.role) return <Navigate to="/select-role" replace />;

  return <Outlet />;
}
