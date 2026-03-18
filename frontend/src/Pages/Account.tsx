import { useAuth } from "@/context/useAuth";
import { backendUrl } from "@/lib/config";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
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
import toast from "react-hot-toast";
import { Link, useNavigate } from "react-router-dom";
import { Package, LogOut } from "lucide-react";

const Account = () => {
  const { user, setUser } = useAuth();
  const navigate = useNavigate();

  const logout = async () => {
    await fetch(`${backendUrl}/auth/logout`, {
      method: "POST",
      credentials: "include",
    });
    setUser(null);
    toast.success("Logged out successfully", { id: "logout-success" });
    navigate("/login", { replace: true });
  };

  return (
    <div className="max-w-md mx-auto mt-10 rounded-2xl border bg-white shadow-sm overflow-hidden">
      {/* user info */}
      <div className="flex items-center gap-4 p-6">
        <Avatar className="h-14 w-14">
          <AvatarImage src={user?.image} alt={user?.name} />
          <AvatarFallback className="bg-red-500 text-white text-xl">
            {user?.name?.charAt(0).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        <div>
          <p className="font-semibold text-lg">{user?.name}</p>
          <p className="text-sm text-muted-foreground">{user?.email}</p>
        </div>
      </div>

      <Separator />

      {/* menu items */}
      <div className="flex flex-col py-2">
        <button className="flex items-center gap-3 px-6 py-3 hover:bg-gray-50 transition-colors text-left">
          <Package className="h-5 w-5 text-red-500" />
          <Link to="/orders">
            <span className="text-sm font-medium">Your Orders</span>
          </Link>
        </button>

        <Separator className="mx-6" />

        {/* logout with confirmation */}
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <button className="flex items-center gap-3 px-6 py-3 hover:bg-gray-50 transition-colors text-left text-red-500">
              <LogOut className="h-5 w-5" />
              <span className="text-sm font-medium">Logout</span>
            </button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Confirm Logout</AlertDialogTitle>
              <AlertDialogDescription>
                Are you sure you want to logout?
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={logout}
                className="bg-red-500 hover:bg-red-600"
              >
                Logout
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
};

export default Account;
