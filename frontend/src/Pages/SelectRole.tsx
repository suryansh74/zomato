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
import { ShoppingBag, UtensilsCrossed, Bike } from "lucide-react";

const roleConfig = {
  customer: {
    icon: ShoppingBag,
    label: "Customer",
    description: "Order food from your favourite restaurants",
  },
  restaurant_owner: {
    icon: UtensilsCrossed,
    label: "Restaurant Owner",
    description: "List and manage your restaurant on Zomato",
  },
  rider: {
    icon: Bike,
    label: "Delivery Rider",
    description: "Deliver orders and earn money",
  },
};

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
        navigate("/", { replace: true });
      } else {
        toast.error("Failed to save role", { id: "role-error" });
      }
    } catch (error) {
      toast.error("Something went wrong", { id: "role-error" });
      console.error(error);
    }
  };

  const roles: Role[] = ["customer", "restaurant_owner", "rider"];

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
      <div className="w-full max-w-lg">
        {/* header */}
        <div className="text-center mb-8">
          <span className="text-3xl font-bold text-red-500">zomato</span>
          <h1 className="text-2xl font-semibold text-gray-800 mt-4">
            How will you use Zomato?
          </h1>
          <p className="text-gray-500 text-sm mt-2">
            Choose your role to get started. This cannot be changed later.
          </p>
        </div>

        {/* role cards */}
        <div className="flex flex-col gap-3 mb-6">
          {roles.map((r) => {
            const config = roleConfig[r!];
            const Icon = config.icon;
            const isSelected = role === r;

            return (
              <button
                key={r}
                onClick={() => setRole(r)}
                className={`flex items-center gap-4 p-4 rounded-xl border-2 text-left transition-all duration-150 bg-white
                  ${
                    isSelected
                      ? "border-red-500 bg-red-50"
                      : "border-gray-200 hover:border-gray-300 hover:bg-gray-50"
                  }`}
              >
                {/* icon */}
                <div
                  className={`p-3 rounded-xl ${isSelected ? "bg-red-100" : "bg-gray-100"}`}
                >
                  <Icon
                    className={`h-6 w-6 ${isSelected ? "text-red-500" : "text-gray-500"}`}
                  />
                </div>

                {/* text */}
                <div className="flex-1">
                  <p
                    className={`font-medium ${isSelected ? "text-red-500" : "text-gray-800"}`}
                  >
                    {config.label}
                  </p>
                  <p className="text-sm text-gray-500">{config.description}</p>
                </div>

                {/* selected indicator */}
                <div
                  className={`h-5 w-5 rounded-full border-2 flex items-center justify-center
                  ${isSelected ? "border-red-500" : "border-gray-300"}`}
                >
                  {isSelected && (
                    <div className="h-2.5 w-2.5 rounded-full bg-red-500" />
                  )}
                </div>
              </button>
            );
          })}
        </div>

        {/* next button */}
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button
              disabled={!role}
              className="w-full h-11 bg-red-500 hover:bg-red-600 text-white"
            >
              Continue
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Confirm your role</AlertDialogTitle>
              <AlertDialogDescription>
                You are about to continue as{" "}
                <strong>{role ? roleConfig[role].label : ""}</strong>. This
                cannot be undone. Are you sure?
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={addRole}
                className="bg-red-500 hover:bg-red-600"
              >
                Yes, confirm
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
}
