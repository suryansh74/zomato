import toast from "react-hot-toast";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useEffect, useState } from "react";
import type { Role } from "../types/types";
import { Button } from "@/components/ui/button";
import { backendUrl } from "@/lib/config";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useAuth } from "@/context/useAuth";

export default function SelectRole() {
  const { setUser } = useAuth();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [role, setRole] = useState<Role>(null);

  useEffect(() => {
    if (searchParams.get("fresh") === "true") {
      toast.success("Logged in successfully!", { id: "login-success" });
      navigate("/select-role", { replace: true });
    }
  }, [searchParams]);

  const addRole = async () => {
    try {
      const res = await fetch(`${backendUrl}/auth/add_role`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ role }),
      });

      if (res.ok) {
        const data = await res.json();
        setUser(data.data.user);
        toast.success("Role selected!");
        console.log("ok working fine");
        navigate("/", { replace: true });
      } else {
        toast.error("Failed to save role", { id: "role-error" });
      }
    } catch (error) {
      toast.error("Something went wrong", { id: "role-error" });
      console.log(error);
    }
  };

  const roles: Role[] = ["customer", "restaurant_owner", "rider"];

  return (
    <div>
      <h1>Select your role</h1>
      <ul>
        {roles.map((r) => (
          <li key={r}>
            <button
              onClick={() => setRole(r)}
              style={{
                border: role === r ? "2px solid red" : "2px solid gray",
                color: role === r ? "red" : "inherit",
                background: role === r ? "#fff0f0" : "transparent",
              }}
            >
              Continue as {r}
            </button>
          </li>
        ))}
      </ul>

      {/* alert dialog on next click */}
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button disabled={!role}>Next</Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Confirm your role</AlertDialogTitle>
            <AlertDialogDescription>
              You are about to continue as <strong>{role}</strong>. This cannot
              be undone. Are you sure?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={addRole}>
              Yes, confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
