import { createContext } from "react";
import type { AuthContextType } from "../types/types";

export const AuthContext = createContext<AuthContextType>({
  user: null,
  setUser: () => {}, // empty function as default
  loading: true,
  location: null,
  loadingLocation: false,
  city: "Fetching Location...",
});
