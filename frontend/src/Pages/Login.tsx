import { useAuth } from "@/context/useAuth";
import toast from "react-hot-toast";
import { Navigate } from "react-router-dom";

export default function Login() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (user && user.role) return <Navigate to="/" replace />;
  if (user && !user.role) return <Navigate to="/select-role" replace />;
  const handleLogin = () => {
    toast.loading("Redirecting to Google...", { id: "google-login" });
    window.location.href = "http://localhost:8000/api/auth/login";
  };

  return <button onClick={handleLogin}>Login with Google</button>;
}
