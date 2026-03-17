import toast from "react-hot-toast";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useEffect } from "react";

export default function Dashboard() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    if (searchParams.get("fresh") === "true") {
      toast.success("Logged in successfully!", { id: "login-success" });
      navigate("/dashboard", { replace: true });
    }
  }, [searchParams]);

  return <h1>Dashboard</h1>;
}
